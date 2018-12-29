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

type SMTPConfig struct {
	Username string `yaml:"username" mapstructure:"username"`
	Password string `yaml:"password" mapstructure:"password"`
	From     string `yaml:"from" mapstructure:"from"`
	Server   string `yaml:"server" mapstructure:"server"`
}

func SMTPExampleConfig() SMTPConfig {
	return SMTPConfig{
		Username: "mritd",
		Password: "password",
		From:     "mritd@mritd.me",
		Server:   "smtp.qq.com:465",
	}
}
