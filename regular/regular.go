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

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/99nil/gopkg/logger"
)

type TaskFunc func(context.Context) error

func (tf TaskFunc) Run(ctx context.Context) error {
	return tf(ctx)
}

type TaskInterface interface {
	Run(ctx context.Context) error
}

func New(cfg *Config) (*Engine, error) {
	return NewWithLogger(cfg, nil)
}

func NewWithLogger(cfg *Config, log logger.UniversalInterface) (*Engine, error) {
	for k, v := range cfg.Periods {
		start, err := time.Parse("15:04", v.Start)
		if err != nil {
			return nil, fmt.Errorf("analysis time period %d start_time(%s) failed: %v", k, v.Start, err)
		}
		end, err := time.Parse("15:04", v.End)
		if err != nil {
			return nil, fmt.Errorf("analysis time period %d end_time(%s) failed: %v", k, v.Start, err)
		}
		cfg.Periods[k].startHour = start.Hour()
		cfg.Periods[k].startMinute = start.Minute()
		cfg.Periods[k].endHour = end.Hour()
		cfg.Periods[k].endMinute = end.Minute()
	}
	if log == nil {
		log = logger.NewEmpty()
	}
	return &Engine{cfg: cfg, log: log}, nil
}

type Engine struct {
	m      sync.Mutex
	cfg    *Config
	log    logger.UniversalInterface
	cancel context.CancelFunc
}

func (e *Engine) SetConfig(cfg *Config) {
	if cfg == nil {
		return
	}

	e.m.Lock()
	defer e.m.Unlock()
	e.cfg = cfg
}

func (e *Engine) GetConfig() *Config {
	e.m.Lock()
	defer e.m.Unlock()
	return e.cfg
}

func (e *Engine) Start(ctx context.Context, task TaskInterface) error {
	for {
		second := time.Now().Second()
		if second == 0 {
			break
		}
		sleepInterval := 60 - second
		e.log.Warnf("The current seconds is not 0, need to wait for %ds to start the automatic assistant", sleepInterval)
		time.Sleep(time.Duration(sleepInterval) * time.Second)
	}

	if len(e.GetConfig().Periods) == 0 {
		return e.run(ctx, task)
	}

	currentStartHour, currentStartMinute := -1, -1
	ticker := time.NewTicker(time.Minute)
	for {
		if e.cancel == nil {
			e.log.Info("Start mission reconnaissance")
		}
		now := time.Now()
		hour := now.Hour()
		minute := now.Minute()

		for _, v := range e.GetConfig().Periods {
			if currentStartHour > -1 && (currentStartHour != v.startHour || currentStartMinute != v.startMinute) {
				continue
			}

			start, end := false, false
			if v.startHour < hour {
				start = true
			}
			if v.startHour == hour && v.startMinute <= minute {
				start = true
			}
			if v.endHour < hour {
				end = true
			}
			if v.endHour == hour && v.endMinute <= minute {
				end = true
			}

			if start && !end && currentStartHour != v.startHour {
				currentStartHour = v.startHour
				currentStartMinute = v.startMinute

				ctx, e.cancel = context.WithCancel(ctx)
				go func() {
					if err := e.run(ctx, task); err != nil {
						e.log.Errorf("Execution ends with error: %v", err)
					}
					e.log.Info("The execution of the current time period is over, please wait for the next time period")
				}()
				break
			}
			if start && end && e.cancel != nil {
				e.cancel()
				e.cancel = nil
			}
		}
		<-ticker.C
	}
}

func (e *Engine) run(ctx context.Context, task TaskInterface) error {
	cfg := e.GetConfig()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := task.Run(ctx); err != nil {
			if cfg.FailInterval < 0 {
				return err
			}
			e.log.Errorf("Execution ends with error: %v", err)
			e.log.Warnf("Will continue after %dms", cfg.FailInterval)
			fmt.Println()
			time.Sleep(time.Duration(cfg.FailInterval) * time.Millisecond)
			continue
		}
		cfg = e.GetConfig()

		if cfg.SuccessInterval < 0 {
			return nil
		}
		e.log.Infof("Executed successfully, will continue after %dms", cfg.SuccessInterval)
		time.Sleep(time.Duration(cfg.SuccessInterval) * time.Millisecond)
	}
}
