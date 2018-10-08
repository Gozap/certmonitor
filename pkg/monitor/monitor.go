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

package monitor

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Gozap/certmonitor/pkg/alarm"

	"github.com/robfig/cron"

	"github.com/Sirupsen/logrus"

	"github.com/spf13/viper"

	"github.com/Gozap/certmonitor/pkg/utils"
)

type Config struct {
	WebSites   []string
	Cron       string
	BeforeTime time.Duration
}

func ExampleConfig() Config {
	return Config{
		WebSites: []string{
			"https://google.com",
			"https://mritd.me",
		},
		Cron:       "@every 1h",
		BeforeTime: 7 * 24 * time.Hour,
	}
}

func check(address string, beforeTime time.Duration) error {

	logrus.Infof("Check website [%s]...", address)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(address)
	if !utils.CheckErr(err) {
		return err
	}
	defer resp.Body.Close()

	for _, cert := range resp.TLS.PeerCertificates {
		if !cert.NotAfter.After(time.Now()) {
			return errors.New(fmt.Sprintf("Website [%s] certificate has expired: %s", address, cert.NotAfter.Local().Format("2006-01-02 15:04:05")))
		}

		if cert.NotAfter.Sub(time.Now()) < beforeTime {
			return errors.New(fmt.Sprintf("Website [%s] certificate will expire, remaining time: %fh", address, cert.NotAfter.Sub(time.Now()).Hours()))
		}
	}

	return nil
}

func Start() {
	var config Config
	err := viper.UnmarshalKey("monitor", &config)
	if err != nil {
		logrus.Fatalf("Can't parse server config: %s", err)
	}

	c := cron.New()

	for _, website := range config.WebSites {
		addr := website
		c.AddFunc(config.Cron, func() {
			err := check(addr, config.BeforeTime)
			if err != nil {
				alarm.Alarm(err.Error())
			}
		})
	}
	c.Start()
	select {}
}
