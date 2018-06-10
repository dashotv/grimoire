// Copyright © 2018 Shawn Catanzarite <me@shawncatz.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/kr/pretty"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dashotv/grimoire/config"
	"github.com/dashotv/grimoire/parser"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the server and read feeds",
	Long:  "Run the server and read feeds",
	Run:   runServer,
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runServer(cmd *cobra.Command, args []string) {
	cfg := &config.Config{}
	err := viper.Unmarshal(cfg)
	if err != nil {
		fmt.Errorf("failed to parse config: %s", err)
	}

	cfg.Compile()

	gp := parser.NewParser(cfg)
	// use cron to run this regularly
	// use goroutines for multiple feed processors
	list := gp.Parse()
	for _, r := range list {
		//fmt.Printf("release: %#v\n", r)
		r.CalculateChecksum()
		pretty.Print(r, "\n")
	}
}
