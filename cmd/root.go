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
	"os"

	"github.com/gozap/certmonitor/conf"

	"github.com/sirupsen/logrus"

	"github.com/gozap/certmonitor/monitor"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

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

		if conf.Monitor.Debug {
			logrus.SetLevel(logrus.DebugLevel)
		}

		if len(args) > 0 {
			_ = cmd.Help()
			return
		}
		monitor.Start()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Info(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "certmonitor.yaml", "config file (default is certmonitor.yaml)")
}

func initConfig() {

	viper.SetConfigFile(cfgFile)

	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		_, err = os.Create(cfgFile)
		if err != nil {
			logrus.Panic(err)
		}
		viper.Set("monitor", conf.MonitorExampleConfig())
		viper.Set("alarm", conf.AlarmExampleConfig())
		viper.Set("acme", conf.ACMEExampleConfig())
		err = viper.WriteConfig()
		if err != nil {
			logrus.Panic(err)
		}
	}

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Panic(err)
	}

	err = viper.UnmarshalKey("monitor", &conf.Monitor)
	if err != nil {
		logrus.Panic(err)
	}
	err = viper.UnmarshalKey("alarm", &conf.Alarm)
	if err != nil {
		logrus.Panic(err)
	}
	err = viper.UnmarshalKey("acme", &conf.ACME)
	if err != nil {
		logrus.Panic(err)
	}
}
