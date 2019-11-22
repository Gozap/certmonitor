package cmd

import (
	"os"

	"github.com/gozap/certmonitor/alarm"

	"github.com/sirupsen/logrus"

	"github.com/gozap/certmonitor/monitor"

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
			_ = cmd.Help()
			return
		}
		monitor.Start()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
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
		_, err = os.Create(cfgFile)
		if err != nil {
			logrus.Fatal(err)
		}
		viper.Set("monitor", monitor.ExampleConfig())
		viper.Set("alarm", alarm.ExampleConfig())
		viper.Set("smtp", alarm.SMTPExampleConfig())
		viper.Set("webhook", alarm.WebHookExampleConfig())
		err = viper.WriteConfig()
		if err != nil {
			logrus.Fatal(err)
		}
	}

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatal(err)
	}
}
