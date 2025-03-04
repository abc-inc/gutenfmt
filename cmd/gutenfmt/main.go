// Copyright 2021 The gutenfmt authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/abc-inc/gutenfmt/gfmt"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var rootCmd = &cobra.Command{
	Use:   "gutenfmt",
	Short: "Formats the input as CSV, JSON, YAML, ASCII table, or name and value pairs.",
	Long: `The gutenfmt utility formats its input to various output formats.

Supported input formats:
- JSON
- Name and value pairs, separated by equal sign or colon.
- Tab-separated name and value pairs

The following output formats are supported:
- csv: Comma-separated values.
- json: JSON string. This setting is the default. Optionally, use --pretty.
- table: ASCII table.
- text: Name and value pairs, separated by equal sign.
- tsv: Tab-separated name and value pairs (useful for grep, sed, or awk).
- yaml: YAML, a machine-readable alternative to JSON. Optionally, use --pretty.
`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ff, err := cmd.Flags().GetString("output")
		if err != nil {
			_ = cmd.Help()
			os.Exit(1)
		}

		if list, _ := cmd.Flags().GetBool("list-themes"); list {
			if !isatty.IsTerminal(os.Stdout.Fd()) {
				fmt.Println(strings.Join(styles.Names(), "\n"))
			} else {
				ex := `{"types": [true, 1, "y"]} // example`
				l := chroma.Coalesce(lexers.Get("json"))
				for _, s := range styles.Names() {
					fmt.Print("Theme: " + s + "\n    ")
					it, _ := l.Tokenise(nil, ex)
					if err := formatters.TTY16m.Format(os.Stdout, styles.Get(s), it); err != nil {
						log.Fatal(err)
					}
					fmt.Print("\n\n")
				}
			}
			return
		}

		if (isatty.IsTerminal(os.Stdin.Fd()) && len(args) == 0) ||
			(!cmd.Flags().Changed("jq") && (cmd.Flags().Changed("arg") || cmd.Flags().Changed("argjson"))) {
			_ = cmd.Help()
			os.Exit(1)
		}

		m := parse(append(args, "-")[0])
		th, _ := cmd.Flags().GetString("theme")
		p, _ := cmd.Flags().GetString("pretty")
		p = strings.ToLower(p)

		var w gfmt.Writer
		switch strings.ToLower(ff) {
		case "csv":
			w = gfmt.NewText(os.Stdout)
			w.(*gfmt.Text).Sep = ","
		case "":
			fallthrough
		case "json":
			if p == "true" || p == "always" || (p == "auto" && isatty.IsTerminal(os.Stdout.Fd())) {
				w = gfmt.NewJSON(os.Stdout, gfmt.WithStyle[gfmt.JSON](styles.Get(th)), gfmt.WithPretty[gfmt.JSON]())
				w.(*gfmt.JSON).Indent = "  "
			} else {
				w = gfmt.NewJSON(os.Stdout, gfmt.WithStyle[gfmt.JSON](styles.Get(th)))
			}
		case "table":
			w = gfmt.NewTab(os.Stdout)
		case "text":
			w = gfmt.NewText(os.Stdout)
			w.(*gfmt.Text).Sep = "="
		case "tsv":
			w = gfmt.NewText(os.Stdout)
			w.(*gfmt.Text).Sep = "\t"
		case "yaml":
			if p == "true" || p == "always" || (p == "auto" && isatty.IsTerminal(os.Stdout.Fd())) {
				w = gfmt.NewYAML(os.Stdout, gfmt.WithStyle[gfmt.YAML](styles.Get(th)), gfmt.WithPretty[gfmt.YAML]())
			} else {
				w = gfmt.NewYAML(os.Stdout, gfmt.WithStyle[gfmt.YAML](styles.Get(th)))
			}
		default:
			_ = cmd.Help()
			os.Exit(1)
		}

		if jq, _ := cmd.Flags().GetString("jq"); jq != "" {
			var allArgs []gfmt.Arg
			args, _ := cmd.Flags().GetStringSlice("arg")
			for _, a := range args {
				arg, err := gfmt.NewArg(a, true)
				if err != nil {
					log.Fatal(err)
				}
				allArgs = append(allArgs, *arg)
			}

			args, _ = cmd.Flags().GetStringSlice("argjson")
			for _, a := range args {
				arg, err := gfmt.NewArg(a, false)
				if err != nil {
					log.Fatal(err)
				}
				allArgs = append(allArgs, *arg)
			}

			raw, _ := cmd.Flags().GetBool("raw-output")
			if raw {
				w = gfmt.NewJQWithArgs(w, jq, allArgs, gfmt.WithRaw())
			} else {
				w = gfmt.NewJQWithArgs(w, jq, allArgs)
			}
		} else if q, _ := cmd.Flags().GetString("query"); q != "" {
			w = gfmt.NewJMESPath(w, q)
		}

		if _, err := w.Write(m); err != nil {
			log.Fatalln("Cannot write output:", err)
		}
	},
}

func main() {
	if runtime.GOOS == "windows" {
		colorable.EnableColorsStdout(nil)
	}

	theme := ""
	if isatty.IsTerminal(os.Stdout.Fd()) {
		theme = "native"
	}

	rootCmd.Flags().StringSlice("arg", nil, "Pass a string value to the jq filter as a predefined variable.")
	rootCmd.Flags().StringSlice("argjson", nil, "Pass a JSON-encoded value to the jq filter as a predefined variable.")
	rootCmd.Flags().String("jq", "", "Specify a jq filter for modifying the output.")
	rootCmd.Flags().Bool("list-themes", false, "Display a list of supported themes for syntax highlighting.")
	rootCmd.Flags().StringP("output", "o", "", "The formatting style for command output (csv, json, table, text, tsv, yaml).")
	rootCmd.Flags().String("pretty", "auto", `Pretty-print the output (JSON or YAML). Possible values are "true"/"always", "false"/"never", "auto".`)
	rootCmd.Flags().StringP("query", "q", "", "Specify a JMESPath query to use in filtering the output")
	rootCmd.Flags().BoolP("raw-output", "r", false, "If the filter's result is a string, then it will be written directly to standard output rather than being formatted as a JSON string with quotes.")
	rootCmd.Flags().String("theme", theme, "Set the theme for syntax highlighting. Use '--list-themes' to see all available themes.")

	rootCmd.MarkFlagsMutuallyExclusive("jq", "query")
	rootCmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Name != "list-themes" {
			rootCmd.MarkFlagsMutuallyExclusive("list-themes", f.Name)
		}
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println()
}

// parse attempts to detect the input format e.g., JSON and returns the value,
// which could be a key-value pairs (map) or a slice thereof.
func parse(name string) any {
	var err error
	r := os.Stdin
	if name != "-" {
		r, err = os.Open(name)
		if err != nil {
			log.Fatalln(err)
		}
		defer func() { _ = r.Close() }()
	}

	bs, err := io.ReadAll(r)
	if err != nil {
		log.Fatalln(err) //nolint:gocritic
	}

	var m any
	kv := map[string]any{}
	d := json.NewDecoder(bytes.NewReader(bs))
	if err = d.Decode(&m); err != nil {
		s := bufio.NewScanner(bytes.NewReader(bs))
		for s.Scan() {
			if idx := bytes.IndexAny(s.Bytes(), "=:\t"); idx > 0 {
				kv[string(s.Bytes()[:idx])] = string(s.Bytes()[idx+1:])
			}
		}
	}
	if len(kv) > 0 {
		return kv
	}
	return m
}
