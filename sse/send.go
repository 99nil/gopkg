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

	for _, opt := range opts {
		switch v := opt.(type) {
		case Buffered:
			if v {
				h.Set("X-Accel-Buffering", "yes")
			} else {
				h.Set("X-Accel-Buffering", "no")
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

func (s *Sender) Close() {
	s.SendError("", "eof")
}

func (s *Sender) Ping() {
	s.Send("", "", "", "", "ping")
}

func (s *Sender) Send(id, event, data, retry, comment string) {
	s.SendMessage(&Message{
		ID:      id,
		Data:    data,
		Event:   event,
		Retry:   retry,
		Comment: comment,
	})
}

func (s *Sender) SendError(id string, vs ...any) {
	ms := make([]*Message, 0, len(vs))
	for _, v := range vs {
		m := &Message{ID: id, Event: "error"}
		switch vv := v.(type) {
		case string:
			m.Data = vv
		case error:
			m.Data = vv.Error()
		default:
			m.Data = fmt.Sprintf("%v", vv)
		}
		ms = append(ms, m)
	}
	s.SendMessage(ms...)
}

func (s *Sender) SendMessage(messages ...*Message) {
	var hasClose bool
	for _, m := range messages {
		msg := m.String()
		if msg == "" {
			continue
		}
		if m.IsCloseMsg() {
			hasClose = true
			_, _ = io.WriteString(s.rw, msg)
			break
		}
		_, _ = io.WriteString(s.rw, msg)
	}

	s.rw.Flush()
	if hasClose {
		close(s.closeCh)
	}
}

func SendLoop[T any](ctx context.Context, s *Sender, dataCh <-chan T, errCh <-chan error) error {
	if s == nil {
		return errors.New("sender is nil")
	}
	if dataCh == nil {
		return errors.New("data chan is nil")
	}
	if errCh == nil {
		errCh = make(chan error)
	}

	s.Ping()
	pingTicker := time.NewTicker(30 * time.Second)
	timeoutTicker := time.NewTicker(24 * time.Hour)
L:
	for {
		pingTicker.Reset(30 * time.Second)
		timeoutTicker.Reset(24 * time.Hour)

		select {
		case <-ctx.Done():
			break L
		case <-timeoutTicker.C:
			break L
		case <-pingTicker.C:
			s.Ping()
		case data := <-dataCh:
			sendData(s, data)
		case err := <-errCh:
			if err == io.EOF {
				for {
					select {
					case data := <-dataCh:
						sendData(s, data)
					default:
						break L
					}
				}
			}
			s.SendError("", err)
		}
	}
	s.Close()
	return nil
}

func sendData(s *Sender, data any) {
	var dataStr string
	switch dv := data.(type) {
	case string:
		dataStr = dv
	case fmt.Stringer:
		dataStr = dv.String()
	default:
		dataBytes, _ := json.Marshal(dv)
		if len(dataBytes) == 0 {
			return
		}
		dataStr = string(dataBytes)
	}

	s.SendMessage(&Message{
		Event: "data",
		Data:  dataStr,
	})
}
