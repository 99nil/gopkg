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

package ctr

import (
	"encoding/json"
	"io"
)

var DefaultDecoder = func(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}

type Validator interface {
	Validate() error
}

func Decode(r io.Reader, v any) error {
	if err := DefaultDecoder(r, v); err != nil {
		return err
	}
	if vv, ok := v.(Validator); ok {
		return vv.Validate()
	}
	return nil
}

func CustomDecode(r io.Reader, v any, decoder func(r io.Reader, v any) error) error {
	if decoder == nil {
		decoder = DefaultDecoder
	}
	if err := decoder(r, v); err != nil {
		return err
	}
	if vv, ok := v.(Validator); ok {
		return vv.Validate()
	}
	return nil
}
