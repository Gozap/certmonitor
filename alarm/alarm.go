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
	"bytes"
	"strings"
	"text/template"
	"time"

	"github.com/gozap/certmonitor/conf"

	"github.com/sirupsen/logrus"
)

func Alarm(website conf.WebsiteConfig, err error) {

	logrus.Debugf("Website alarm: %s", err.Error())
	switch strings.ToLower(conf.Monitor.AlarmType) {
	case "all":
		smtpAlarm(website, err)
		webhooksAlarm(website, err)
	case "smtp":
		smtpAlarm(website, err)
	case "webhook":
		webhooksAlarm(website, err)
	default:
		logrus.Error("Alarm type not support!")
	}
}

func smtpAlarm(_ conf.WebsiteConfig, err error) {
	config := SMTPConfig{
		User:     conf.Alarm.SMTP.User,
		Password: conf.Alarm.SMTP.Password,
		Server:   conf.Alarm.SMTP.Server,
		From:     conf.Alarm.SMTP.From,
	}
	config.Send(conf.Alarm.SMTP.Targets, err.Error())
}

func webhooksAlarm(website conf.WebsiteConfig, err error) {

	for _, cfg := range conf.Alarm.HTTP {

		var targets []string
		config := WebHookConfig{
			Method:  strings.ToLower(cfg.Method),
			TimeOut: 5 * time.Second,
		}
		if cfg.NeedParse {
			for _, addr := range cfg.Targets {
				t, err := template.New("").Parse(addr)
				if err != nil {
					logrus.Error(err)
					continue
				}
				var buf bytes.Buffer
				err = t.Execute(&buf, website)
				if err != nil {
					logrus.Error(err)
					continue
				}
				targets = append(targets, buf.String())
			}
		} else {
			targets = cfg.Targets
		}
		config.Send(targets, err.Error())
	}

}
