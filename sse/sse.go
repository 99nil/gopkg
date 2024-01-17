// Copyright © 2023 zc2638 <zc2638@qq.com>.
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
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"strconv"
	"strings"
	"time"
)

const (
	headerID    = "id:"
	headerData  = "data:"
	headerEvent = "event:"
	headerRetry = "retry:"
)

const defaultBufferSize = 4096

var CloseMessage = &Message{Data: io.EOF.Error(), Event: "error"}

type Message struct {
	ID      string
	Data    string
	Event   string
	Retry   time.Duration // ms
	Comment string
}

func (m *Message) IsClose() bool {
	return m.Event == "error" && m.Data == io.EOF.Error()
}

func (m *Message) String() string {
	var sb strings.Builder

	if m.ID != "" {
		sb.WriteString("id: ")
		sb.WriteString(m.ID)
		sb.WriteString("\n")
	}
	if m.Event != "" {
		sb.WriteString("event: ")
		sb.WriteString(m.Event)
		sb.WriteString("\n")
	}
	if m.Data != "" {
		sb.WriteString("data: ")
		sb.WriteString(m.Data)
		sb.WriteString("\n")
	}
	if m.Retry > 0 {
		sb.WriteString("retry: ")
		sb.WriteString(strconv.FormatInt(m.Retry.Milliseconds(), 10))
		sb.WriteString("\n")
	}
	if m.Comment != "" {
		sb.WriteString(": ")
		sb.WriteString(m.Comment)
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	return sb.String()
}

type Parser struct {
	scanner *bufio.Scanner
}

func NewParser(r io.Reader) *Parser {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, defaultBufferSize), defaultBufferSize)
	scanner.Split(func(data []byte, atEOF bool) (int, []byte, error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i, nlen := containsDoubleNewline(data); i >= 0 {
			return i + nlen, data[0:i], nil
		}
		if atEOF {
			return len(data), data, nil
		}
		return 0, nil, nil
	})
	return &Parser{scanner: scanner}
}

func (p *Parser) ReadLoop(callback func(event *Message, err error) error) error {
	if callback == nil {
		return errors.New("callback must be defined")
	}
	for {
		message, err := p.Read()
		if err := callback(message, err); err != nil {
			return err
		}
	}
}

func (p *Parser) Read() (*Message, error) {
	scanner := p.scanner
	if scanner.Scan() {
		msg := scanner.Text()
		if len(msg) < 1 {
			return nil, errors.New("event message was empty")
		}
		return processEventMsg(msg), nil
	}
	if err := scanner.Err(); err != nil && !errors.Is(err, context.Canceled) {
		return nil, err
	}
	return nil, io.EOF
}

func processEventMsg(msg string) *Message {
	message := new(Message)
	for _, line := range strings.FieldsFunc(msg, func(r rune) bool { return r == '\n' || r == '\r' }) {
		switch {
		case strings.HasPrefix(line, headerID):
			message.ID = trimPrefix(line, headerID)
		case strings.HasPrefix(line, headerData):
			data := trimPrefix(line, headerData)
			if message.Data != "" {
				message.Data += "\n"
			}
			message.Data += data
		case line == "data":
			// The spec says that a line that simply contains the string "data"
			// should be treated as a data field with an empty body.
			message.Data += "\n"
		case strings.HasPrefix(line, headerEvent):
			message.Event = trimPrefix(line, headerEvent)
		case strings.HasPrefix(line, headerRetry):
			v := trimPrefix(line, headerRetry)
			d, _ := strconv.ParseInt(v, 10, 64)
			message.Retry = time.Duration(d) * time.Millisecond
		default:
		}
	}
	return message
}

func trimPrefix(data string, prefix string) string {
	size := len(prefix)
	if len(data) <= size {
		return data
	}

	data = data[size:]
	// Remove optional leading whitespace
	if len(data) > 0 && data[0] == 32 {
		data = data[1:]
	}
	// Remove trailing new line
	if len(data) > 0 && data[len(data)-1] == 10 {
		data = data[:len(data)-1]
	}
	return data
}

func containsDoubleNewline(data []byte) (int, int) {
	if index := bytes.Index(data, []byte("\r\n\r\n")); index > -1 {
		return index, 4
	}
	if index := bytes.Index(data, []byte("\r\r")); index > -1 {
		return index, 2
	}
	if index := bytes.Index(data, []byte("\n\n")); index > -1 {
		return index, 2
	}
	if index := bytes.Index(data, []byte("\r\n\n")); index > -1 {
		return index, 3
	}
	if index := bytes.Index(data, []byte("\n\r\n")); index > -1 {
		return index, 3
	}
	return -1, 2
}
