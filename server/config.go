// Copyright © 2021 zc2638 <zc2638@qq.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// ParseConfig Parsing configuration files
func ParseConfig(configPath string, data interface{}) error {
	return ParseConfigWithEnv(configPath, data, "")
}

// ParseConfigWithEnv Parsing configuration files with env
// You can customize your judgment according by `err.(*os.PathError)`
// You can get the configuration file path by `viper.ConfigFileUsed()`
func ParseConfigWithEnv(configPath string, data interface{}, envPrefix string) error {
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}
		viper.AddConfigPath(home)
		viper.SetConfigName("config.yaml")
	}
	if envPrefix != "" {
		viper.SetEnvPrefix(envPrefix)
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()
	}
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(*os.PathError); !ok {
			return err
		}
	}
	return viper.Unmarshal(data, func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "json"
	})
}
