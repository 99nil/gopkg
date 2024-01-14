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

package ctr

import "fmt"

type Log interface {
	Error(vs ...any)
}

// logInstance defines the log processing method
var logger Log = new(emptyLogger)

// Logger returns the default logger instance
func Logger() Log {
	return logger
}

func SetLog(l Log) {
	if l != nil {
		logger = l
	}
}

// InitLogger initialization the log processing method
// Deprecated: SetLog instead.
var InitLogger = SetLog

type emptyLogger struct{}

func (l *emptyLogger) Error(vs ...any) {}

func CoverKVLog(logger KVLog) Log {
	return &coverLog{logger: logger}
}

type KVLog interface {
	Error(msg string, vs ...any)
}

type coverLog struct {
	logger KVLog
}

func (l *coverLog) Error(vs ...any) {
	for _, v := range vs {
		str, ok := v.(string)
		if !ok {
			str = fmt.Sprintf("%v", v)
		}
		l.logger.Error(str)
	}
}
