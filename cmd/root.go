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
	"os"

	"github.com/Gozap/certmonitor/pkg/alarm"

	"github.com/Sirupsen/logrus"

	"github.com/Gozap/certmonitor/pkg/monitor"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var debug bool

var (
	Version   string
	BuildTime string
	CommitID  string
)

var rootCmd = &cobra.Command{
	Use:   "certmonitor",
	Short: "A simple website certificate monitor tool",
	Long: `
A simple website certificate monitor tool.`,
	Run: func(cmd *cobra.Command, args []string) {

		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})

		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}

		if len(args) > 0 {
			cmd.Help()
			return
		}
		monitor.Start()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "certmonitor.yaml", "config file (default is certmonitor.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug")
}

func initConfig() {

	viper.SetConfigFile(cfgFile)

	if _, err := os.Stat(cfgFile); err != nil {
		os.Create(cfgFile)
		viper.Set("monitor", monitor.ExampleConfig())
		viper.Set("alarm", alarm.ExampleConfig())
		viper.Set("smtp", alarm.SMTPExampleConfig())
		viper.Set("webhook", alarm.WebHookExampleConfig())
		viper.WriteConfig()
	}

	viper.AutomaticEnv()
	viper.ReadInConfig()
}
