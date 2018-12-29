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

var ACME ACMEConfig

type ACMEProviderConfig struct {
	Name      string `yaml:"name" mapstructure:"name"`
	APIKey    string `yaml:"api_key" mapstructure:"api_key"`
	APISecret string `yaml:"api_secret" mapstructure:"api_secret"`
}

type ACMEConfig struct {
	Email     string               `yaml:"email" mapstructure:"email"`
	Providers []ACMEProviderConfig `yaml:"providers" mapstructure:"providers"`
}

func ACMEExampleConfig() ACMEConfig {
	return ACMEConfig{
		Email: "mritd@mritd.me",
		Providers: []ACMEProviderConfig{
			{
				Name:      "godaddy",
				APIKey:    "e7UFaqWtRfMcPvvosuDqFYkQNVmkEWNY",
				APISecret: "XjfmmniwA3tyHYRLjTCXFcBTGRxFTpYU",
			},
			{
				Name:      "alidns",
				APIKey:    "CVYGfcEJBjgzGuNFbyUXYkj7QpLqMGfM",
				APISecret: "hTeQbsdAkXGTsPVJpsdzUE7RjaYtXAym",
			},
		},
	}
}
