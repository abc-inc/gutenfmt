/*
Copyright © 2021 The gutenfmt authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/abc-inc/gutenfmt/gfmt"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "whatever",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		f, err := cmd.Flags().GetString("format")
		if err != nil {
			cmd.Help()
			os.Exit(1)
		}

		m := map[string]interface{}{}
		in := bufio.NewScanner(os.Stdin)
		for in.Scan() {
			b := []byte(in.Text())
			if json.Valid(b) {
				json.Unmarshal(b, &m)
			} else if strings.Contains(in.Text(), "=") {
				kv := strings.SplitN(in.Text(), "=", 2)
				m[kv[0]] = kv[1]
			}
		}

		var o gfmt.InterfaceWriter
		switch strings.ToLower(f) {
		case "coljson":
			o = gfmt.NewColJSON(os.Stdout)
		case "json":
			o = gfmt.NewJSON(os.Stdout)
		case "text":
			o = gfmt.NewText(os.Stdout)
		case "tab":
			o = gfmt.NewTab(os.Stdout)
		default:
			cmd.Help()
			os.Exit(1)
		}

		o.Write([]map[string]interface{}{m})
	},
}

func main() {
	Execute()
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().String("format", "", "todo")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
