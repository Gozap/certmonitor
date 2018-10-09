/*
 * Copyright 2018 Gozap, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var versionTpl = `
Name: certmonitor
Version: %s
Arch: %s
BuildTime: %s
CommitID: %s
`

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Long: `
Print version.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(versionTpl, Version, runtime.GOOS+"/"+runtime.GOARCH, BuildTime, CommitID)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
