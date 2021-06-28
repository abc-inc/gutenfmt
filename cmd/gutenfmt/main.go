/**
 * Copyright 2021 The gutenfmt authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/abc-inc/gutenfmt/gfmt"
	"github.com/mattn/go-colorable"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gutenfmt",
	Short: "Formats the input as JSON, YAML, ASCII table or name and value pairs.",
	Long: `The gutenfmt utility formats its input to various output formats.

Supported input formats:
- JSON
- Name and value pairs, separated by equal sign or colon.
- Tab-separated name and value pairs

The following output formats are supported:
- json: JSON string. This setting is the default for non-terminals.
- jsonc: Colorized JSON. This setting is the default for interactive terminals.
- table: ASCII table.
- text: Name and value pairs, separated by equal sign.
- tsv: Tab-separated name and value pairs (useful for grep, sed, or awk).
- yaml: YAML, a machine-readable alternative to JSON.
- yamlc: Colorized YAML.
`,
	Run: func(cmd *cobra.Command, args []string) {
		ff, err := cmd.Flags().GetString("output")
		if err != nil {
			_ = cmd.Help()
			os.Exit(1)
		}

		m := parse()
		var w gfmt.Writer
		switch strings.ToLower(ff) {
		case "":
			w = gfmt.NewAutoJSON(os.Stdout)
		case "json":
			w = gfmt.NewJSON(os.Stdout)
		case "jsonc":
			w = gfmt.NewPrettyJSON(os.Stdout)
		case "table":
			w = gfmt.NewTab(os.Stdout)
		case "text":
			w = gfmt.NewText(os.Stdout)
			w.(*gfmt.Text).Sep = "="
		case "tsv":
			w = gfmt.NewText(os.Stdout)
			w.(*gfmt.Text).Sep = "\t"
		case "yaml":
			w = gfmt.NewYAML(os.Stdout)
		case "yamlc":
			w = gfmt.NewPrettyYAML(os.Stdout)
		default:
			_ = cmd.Help()
			os.Exit(1)
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

	rootCmd.Flags().StringP("output", "o", "",
		"The formatting style for command output (json, jsonc, table, text, tsv, yaml).")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// parse attempts to detect the input format e.g., JSON and returns the value,
// which could be a key-value pairs (map) or a slice thereof.
func parse() interface{} {
	m := map[string]interface{}{}
	in := bufio.NewScanner(os.Stdin)
	for in.Scan() {
		s := in.Text()
		b := []byte(s)
		if json.Valid(b) {
			if b[0] == '[' {
				m2 := []interface{}{}
				if err := json.Unmarshal(b, &m2); err != nil {
					log.Fatalln("Cannot output JSON:", err)
				}
				m[""] = m2
			} else if err := json.Unmarshal(b, &m); err != nil {
				log.Fatalln("Cannot output JSON:", err)
			}
		} else if idx := strings.IndexAny(s, "=:\t"); idx > 0 {
			m[s[:idx]] = s[idx+1:]
		}
	}
	if _, ok := m[""]; ok {
		return m[""]
	}
	return m
}
