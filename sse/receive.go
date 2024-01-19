// Copyright © 2024 zc2638 <zc2638@qq.com>.
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

package sse

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
)

type ReceiveDataEvent string

func NewReceiver[T any](reader io.Reader, coverFn func(data []byte) (T, bool), opts ...any) *Receiver[T] {
	if coverFn == nil {
		coverFn = func(data []byte) (T, bool) {
			var out T
			err := json.Unmarshal(data, &out)
			return out, err == nil
		}
	}

	r := &Receiver[T]{
		dataCh:    make(chan T),
		errCh:     make(chan error),
		runCh:     make(chan struct{}),
		parser:    NewParser(reader),
		coverFn:   coverFn,
		dataEvent: "message",
	}
	for _, opt := range opts {
		switch v := opt.(type) {
		case ReceiveDataEvent:
			r.dataEvent = string(v)
		}
	}
	return r
}

type Receiver[T any] struct {
	dataCh chan T
	errCh  chan error
	runCh  chan struct{}

	parser    *Parser
	coverFn   func(data []byte) (T, bool)
	dataEvent string
}

func (r *Receiver[T]) IsClosed() bool {
	select {
	case <-r.runCh:
		return false
	default:
	}
	return true
}

func (r *Receiver[T]) Data() <-chan T {
	return r.dataCh
}

func (r *Receiver[T]) Err() <-chan error {
	return r.errCh
}

func (r *Receiver[T]) Run(ctx context.Context) {
	if !r.IsClosed() {
		return
	}
	close(r.runCh)

	err := r.parser.ReadLoop(func(message *Message, err error) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if err != nil {
			return err
		}

		if message.Event == "error" {
			if strings.ToUpper(message.Data) == io.EOF.Error() {
				return io.EOF
			}
			r.errCh <- errors.New(message.Data)
			return nil
		}
		if message.Event != r.dataEvent {
			return nil
		}

		data, ok := r.coverFn([]byte(message.Data))
		if ok {
			r.dataCh <- data
		}
		return nil
	})
	if err != nil {
		r.errCh <- err
	}
	close(r.errCh)
	r.runCh = make(chan struct{})
}
