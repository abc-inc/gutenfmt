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
	"strings"

	"github.com/abc-inc/gutenfmt/gfmt"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gutenfmt",
	Short: "Formats the input as JSON, key-value pairs or table.",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ff, err := cmd.Flags().GetString("output")
		if err != nil || ff == "" {
			_ = cmd.Help()
			os.Exit(1)
		}

		m := map[string]interface{}{}
		in := bufio.NewScanner(os.Stdin)
		for in.Scan() {
			s := in.Text()
			b := []byte(s)
			if json.Valid(b) {
				if err := json.Unmarshal(b, &m); err != nil {
					log.Fatalln("Cannot output JSON:", err)
				}
			} else if idx := strings.IndexAny(s, "=:\t"); idx > 0 {
				m[s[:idx]] = s[idx+1:]
			}
		}

		var f gfmt.GenericWriter
		switch strings.ToLower(ff) {
		case "json":
			f = gfmt.NewJSON(os.Stdout)
		case "jsonc":
			f = gfmt.NewPrettyJSON(os.Stdout)
		case "table":
			f = gfmt.NewTab(os.Stdout)
		case "text":
			f = gfmt.NewText(os.Stdout)
		case "tsv":
			f = gfmt.NewText(os.Stdout)
			f.(*gfmt.Text).Sep = "\t"
		default:
			_ = cmd.Help()
			os.Exit(1)
		}

		if _, err := f.Write(m); err != nil {
			log.Fatalln("Cannot write output:", err)
		}
	},
}

func main() {
	rootCmd.Flags().StringP("output", "o", "",
		"The formatting style for command output (json, jsonc, table, text, tsv).")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
