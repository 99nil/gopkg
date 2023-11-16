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

package printer

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

const (
	tabWriterMinWidth = 6
	tabWriterWidth    = 4
	tabWriterPadding  = 3
	tabWriterPadChar  = ' '
)

func New() *tabwriter.Writer {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, tabWriterMinWidth, tabWriterWidth, tabWriterPadding, tabWriterPadChar, 0)
	return w
}

func NewTab(header ...string) *Tab {
	return &Tab{Header: header}
}

type Tab struct {
	Header []string
	Data   [][]any
}

func (t *Tab) Add(vs ...string) {
	var data []any
	for _, v := range vs {
		data = append(data, v)
	}
	t.Data = append(t.Data, data)
}

func (t *Tab) Print() {
	var num int
	if t.Header != nil {
		num = len(t.Header)
	} else if len(t.Data) > 0 {
		num = len(t.Data[0])
	} else {
		return
	}
	formats := make([]string, 0, num)
	for i := 0; i < num; i++ {
		formats = append(formats, "%s")
	}
	format := strings.Join(formats, "\t") + "\n"
	w := New()

	if t.Header != nil {
		var header []any
		for _, h := range t.Header {
			header = append(header, h)
		}
		fmt.Fprintf(w, format, header...)
	}
	for _, v := range t.Data {
		fmt.Fprintf(w, format, v...)
	}
	w.Flush()
	t.Header = nil
	t.Data = nil
}
