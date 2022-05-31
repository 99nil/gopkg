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

import (
	"github.com/99nil/go/logger"
)

// logInstance defines the log processing method
var logInstance = logger.NewEmpty()

// Logger returns the default logger instance
func Logger() logger.Interface {
	return logInstance
}

// InitLogger initialization the log processing method
func InitLogger(ins logger.Interface) {
	logInstance = ins
}

// InitStdLogger initialization the Std logger processing method
func InitStdLogger(log logger.StdInterface) {
	logInstance = logger.NewEmptyWithStd(log)
}
