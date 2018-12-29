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
	"fmt"
	"net/http"
	"time"

	"github.com/gozap/certmonitor/alarm"

	"github.com/robfig/cron"

	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"

	"github.com/gozap/certmonitor/utils"
)

type Config struct {
	WebSites   []string
	Cron       string
	BeforeTime time.Duration
	TimeOut    time.Duration
}

type WebSiteError struct {
	Message string
}

func (e *WebSiteError) Error() string {
	return e.Message
}

func NewWebSiteError(msg string) *WebSiteError {
	return &WebSiteError{
		Message: msg,
	}
}

func ExampleConfig() Config {
	return Config{
		WebSites: []string{
			"https://google.com",
			"https://mritd.me",
		},
		Cron:       "@every 1h",
		BeforeTime: 7 * 24 * time.Hour,
		TimeOut:    10 * time.Second,
	}
}

func check(address string, beforeTime, timeout time.Duration) *WebSiteError {

	logrus.Infof("Check website [%s]...", address)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}
	resp, err := client.Get(address)
	if !utils.CheckErr(err) {
		return nil
	}
	defer resp.Body.Close()

	for _, cert := range resp.TLS.PeerCertificates {
		if !cert.NotAfter.After(time.Now()) {
			return NewWebSiteError(fmt.Sprintf("Website [%s] certificate has expired: %s", address, cert.NotAfter.Local().Format("2006-01-02 15:04:05")))
		}

		if cert.NotAfter.Sub(time.Now()) < beforeTime {
			return NewWebSiteError(fmt.Sprintf("Website [%s] certificate will expire, remaining time: %fh", address, cert.NotAfter.Sub(time.Now()).Hours()))
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
			err := check(addr, config.BeforeTime, config.TimeOut)
			if err != nil {
				alarm.Alarm(err.Error())
			}
		})
	}
	c.Start()
	select {}
}
