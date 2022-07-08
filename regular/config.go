// Copyright © 2022 zc2638 <zc2638@qq.com>.
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

package regular

import "time"

var (
	AllDay = Period{Start: "00:00", End: "23:59"}
	NoDay  = Period{Start: "00:00", End: "00:00"}
)

func NewDefaultConfig() *Config {
	return &Config{
		SuccessInterval: -1,
		FailInterval:    -1,
	}
}

type Config struct {
	// SuccessInterval defines successful execution interval time (ms) for re-execution,
	// -1 is to stop execution.
	SuccessInterval int `json:"success_interval"`

	// FailInterval defines failure execution interval time (ms) for re-execution,
	// -1 is to stop execution.
	FailInterval int `json:"fail_interval"`

	// Periods defines the execution cycle
	Periods []*Period `json:"periods"`
}

type Period struct {
	// Start defines the start time (e.g. 00:00)
	Start string `json:"start"`

	// End defines the end time (e.g. 23:59)
	End string `json:"end"`

	startHour   int
	startMinute int
	endHour     int
	endMinute   int
}

func (p *Period) Parse() error {
	start, err := time.Parse("15:04", p.Start)
	if err != nil {
		return err
	}
	end, err := time.Parse("15:04", p.End)
	if err != nil {
		return err
	}
	p.startHour = start.Hour()
	p.startMinute = start.Minute()
	p.endHour = end.Hour()
	p.endMinute = end.Minute()
	return nil
}
