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
	"fmt"
	"io"
	"net/http"
	"time"
)

type Buffered bool

type ResponseWriter interface {
	http.ResponseWriter
	http.Flusher
}

func NewSender(w http.ResponseWriter, opts ...any) (*Sender, error) {
	rw, ok := w.(ResponseWriter)
	if !ok {
		return nil, errors.New("response writer not support flush")
	}

	h := rw.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")
	h.Set("Connection", "keep-alive")
	h.Set("X-Accel-Buffering", "no")

	for _, opt := range opts {
		switch v := opt.(type) {
		case Buffered:
			if v {
				h.Set("X-Accel-Buffering", "yes")
			}
		}
	}

	return &Sender{
		rw:      rw,
		closeCh: make(chan struct{}),
	}, nil
}

type Sender struct {
	rw      ResponseWriter
	closeCh chan struct{}
}

func (s *Sender) WaitForClose() <-chan struct{} {
	return s.closeCh
}

func (s *Sender) IsClosed() bool {
	select {
	case <-s.closeCh:
		return true
	default:
	}
	return false
}

func (s *Sender) Close() {
	s.SendError(io.EOF)
}

func (s *Sender) Ping() {
	s.Send(&Message{Comment: "ping"})
}

func (s *Sender) SendComment(comment string) {
	if comment == "" {
		return
	}
	s.Send(&Message{Comment: comment})
}

func (s *Sender) SendError(v any) {
	if v == nil {
		return
	}

	m := &Message{Event: "error"}
	switch vv := v.(type) {
	case string:
		m.Data = vv
	case error:
		m.Data = vv.Error()
	default:
		m.Data = fmt.Sprintf("%v", vv)
	}
	s.Send(m)
}

func (s *Sender) Send(messages ...*Message) {
	if s.IsClosed() {
		return
	}

	for _, m := range messages {
		if m == nil {
			continue
		}

		msg := m.String()
		if msg == "" {
			continue
		}
		if m.IsClose() {
			close(s.closeCh)
			_, _ = io.WriteString(s.rw, msg)
			break
		}
		_, _ = io.WriteString(s.rw, msg)
	}
	s.rw.Flush()
}

func SendLoop[T any](
	ctx context.Context,
	s *Sender,
	dataCh <-chan T,
	coverFn func(data T) ([]*Message, error),
	pingInterval time.Duration,
	timeout time.Duration,
) error {
	return SendLoopWithErr[T](ctx, s, dataCh, nil, coverFn, pingInterval, timeout)
}

func SendLoopWithErr[T any](
	ctx context.Context,
	s *Sender,
	dataCh <-chan T,
	errCh <-chan error,
	coverFn func(data T) ([]*Message, error),
	pingInterval time.Duration,
	timeout time.Duration,
) error {
	if s == nil {
		return errors.New("sender is nil")
	}
	if dataCh == nil {
		return errors.New("data chan is nil")
	}
	if s.IsClosed() {
		return nil
	}
	if coverFn == nil {
		coverFn = func(data T) ([]*Message, error) {
			return sendCoverFunc(data)
		}
	}
	if errCh == nil {
		errCh = make(chan error)
	}

	if pingInterval == 0 {
		pingInterval = 30 * time.Second
	}
	if timeout == 0 {
		timeout = 24 * time.Hour
	}
	pingTicker := time.NewTicker(pingInterval)
	timeoutTicker := time.NewTimer(timeout)
	s.Ping()

L:
	for {
		select {
		case <-s.WaitForClose():
			return nil
		case <-ctx.Done():
			break L
		case <-timeoutTicker.C:
			break L
		case <-pingTicker.C:
			s.Ping()
		case data, ok := <-dataCh:
			msgs, err := coverFn(data)
			if err != nil {
				s.SendError(err)
			} else {
				s.Send(msgs...)
			}
			if !ok {
				break L
			}
		case err := <-errCh:
			if err == io.EOF {
				for {
					select {
					case <-s.WaitForClose():
						return nil
					case <-ctx.Done():
						break L
					case data, ok := <-dataCh:
						msgs, coverErr := coverFn(data)
						if coverErr != nil {
							s.SendError(coverErr)
						} else {
							s.Send(msgs...)
						}
						if !ok {
							break L
						}
					default:
						break L
					}
				}
			}
			if err != nil {
				s.SendError(err)
			}
		}
	}
	s.Close()
	return nil
}

func sendCoverFunc(data any) ([]*Message, error) {
	if data == nil {
		return nil, nil
	}

	var dataStr string
	switch dv := data.(type) {
	case Message:
		return []*Message{&dv}, nil
	case *Message:
		return []*Message{dv}, nil
	case string:
		dataStr = dv
	case fmt.Stringer:
		dataStr = dv.String()
	default:
		dataBytes, _ := json.Marshal(dv)
		if len(dataBytes) == 0 {
			return nil, nil
		}
		dataStr = string(dataBytes)
	}

	if data == "" {
		return nil, nil
	}
	return []*Message{{Event: "data", Data: dataStr}}, nil
}
