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

package database

import (
	"time"

	"github.com/go-sql-driver/mysql"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Config struct {
	Addr   string `json:"addr" yaml:"addr"`
	User   string `json:"user" yaml:"user"`
	Pwd    string `json:"pwd" yaml:"pwd"`
	DBName string `json:"dbname" yaml:"dbname"`
	Debug  bool   `json:"debug" yaml:"debug"`
}

func (c *Config) Clone() *Config {
	return &Config{
		Addr:   c.Addr,
		User:   c.User,
		Pwd:    c.Pwd,
		DBName: c.DBName,
	}
}

func New(cfg *Config) (*gorm.DB, error) {
	config := mysql.Config{
		Addr:                 cfg.Addr,
		User:                 cfg.User,
		Passwd:               cfg.Pwd,
		DBName:               cfg.DBName,
		Net:                  "tcp",
		Collation:            "utf8mb4_general_ci",
		ParseTime:            true,
		Loc:                  time.UTC,
		AllowNativePasswords: true,
	}
	dbConnect, err := gorm.Open(
		gormmysql.Open(config.FormatDSN()),
		&gorm.Config{
			SkipDefaultTransaction: true,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	if cfg.Debug {
		dbConnect = dbConnect.Debug()
	}
	return dbConnect, nil
}
