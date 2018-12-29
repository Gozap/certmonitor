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

var Alarm AlarmConfig

type SMTPAlarmConfig struct {
	Username string   `yaml:"username" mapstructure:"username"`
	Password string   `yaml:"password" mapstructure:"password"`
	From     string   `yaml:"from" mapstructure:"from"`
	Server   string   `yaml:"server" mapstructure:"server"`
	Targets  []string `yaml:"targets" mapstructure:"targets"`
}

type HTTPAlarmConfig struct {
	Method    string   `yaml:"method" mapstructure:"method"`
	Targets   []string `yaml:"targets" mapstructure:"targets"`
	NeedParse bool     `yaml:"need_parse" mapstructure:"need_parse"`
}

type AlarmConfig struct {
	SMTP SMTPAlarmConfig   `yaml:"smtp" mapstructure:"smtp"`
	HTTP []HTTPAlarmConfig `yaml:"http" mapstructure:"http"`
}

func AlarmExampleConfig() AlarmConfig {
	return AlarmConfig{
		SMTP: SMTPAlarmConfig{
			Username: "mritd",
			Password: "password",
			From:     "mritd@mritd.me",
			Server:   "smtp.qq.com:465",
			Targets: []string{
				"mritd@mritd.me",
				"mritd@mritd.com",
			},
		},
		HTTP: []HTTPAlarmConfig{
			{
				Targets: []string{
					"https://mritd.me",
					"https://mritd.com",
				},
				Method:    "GET",
				NeedParse: false,
			},

			{
				Targets: []string{
					"https://baidu.com?q={{ .Domain }}",
					"https://google.com?q={{ .Domain }}",
				},
				Method:    "GET",
				NeedParse: true,
			},
		},
	}
}
