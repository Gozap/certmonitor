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
		logrus.Printf("Can't parse alarm config: %s", err)
		return
	}
	for _, a := range alarm {
		switch strings.ToLower(a.Type) {
		case "smtp":
			var s SMTPConfig
			err := viper.UnmarshalKey("smtp", &s)
			if err != nil {
				logrus.Printf("Can't parse smtp config: %s", err)
				return
			}
			s.Send(a.Targets, message)
		case "webhook":
			var w WebHookConfig
			err := viper.UnmarshalKey("webhook", &w)
			if err != nil {
				logrus.Printf("Can't parse webhook config: %s", err)
				return
			}
			w.Send(a.Targets, message)
		default:
			logrus.Print("Alarm type not support!")

		}
	}
}