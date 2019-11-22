package alarm

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Alarm(message string) {

	logrus.Debugf("Website alarm: %s", message)

	var alarm []Config
	err := viper.UnmarshalKey("alarm", &alarm)
	if err != nil {
		logrus.Infof("Can't parse alarm config: %s", err)
		return
	}
	for _, a := range alarm {
		switch strings.ToLower(a.Type) {
		case "smtp":
			var s SMTPConfig
			err := viper.UnmarshalKey("smtp", &s)
			if err != nil {
				logrus.Errorf("Can't parse smtp config: %s", err)
				return
			}
			s.Send(a.Targets, message)
		case "webhook":
			var w WebHookConfig
			err := viper.UnmarshalKey("webhook", &w)
			if err != nil {
				logrus.Errorf("Can't parse webhook config: %s", err)
				return
			}
			w.Send(a.Targets, message)
		default:
			logrus.Error("Alarm type not support!")
		}
	}
}
