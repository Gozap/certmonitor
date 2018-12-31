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

package conf

import "time"

var Monitor MonitorConfig

type WebsiteConfig struct {
	Domain      string `yaml:"domain" mapstructure:"domain"`
	Command     string `yaml:"command" mapstructure:"command"`
	AutoRenew   bool   `yaml:"auto_renew" mapstructure:"auto_renew"`
	DNSProvider string `yaml:"dns_provider" mapstructure:"dns_provider"`
}

type MonitorConfig struct {
	Debug       bool            `yaml:"debug" mapstructure:"debug"`
	Websites    []WebsiteConfig `yaml:"websites" mapstructure:"websites"`
	Cron        string          `yaml:"cron" mapstructure:"cron"`
	AlarmType   string          `yaml:"alarm_type" mapstructure:"alarm_type"`
	HttpTimeout time.Duration   `yaml:"http_timeout" mapstructure:"http_timeout"`
	BeforeTime  time.Duration   `yaml:"before_time" mapstructure:"before_time"`
}

func MonitorExampleConfig() MonitorConfig {
	return MonitorConfig{
		Debug:       true,
		Cron:        "@every 4h",
		AlarmType:   "all",
		BeforeTime:  7 * 24 * time.Hour,
		HttpTimeout: 5 * time.Second,
		Websites: []WebsiteConfig{
			{
				Domain:    "mritd.me",
				Command:   "bash ~/copy_cert.sh",
				AutoRenew: true,
			},
		},
	}
}
