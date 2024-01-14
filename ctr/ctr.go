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
	"encoding/json"
	"errors"
	"net/http"
)

// Bytes writes the Bytes message to the response.
func Bytes(w http.ResponseWriter, bytes []byte) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(bytes)
}

// Str writes the string message to the response.
func Str(w http.ResponseWriter, str string) {
	Bytes(w, []byte(str))
}

// Success writes ok message to the response.
func Success(w http.ResponseWriter) {
	Str(w, "success")
}

// NoContent writes no content to the response.
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// Redirect replies to the request with a redirect to url.
func Redirect(w http.ResponseWriter, r *http.Request, url string, code int) {
	http.Redirect(w, r, url, code)
}

// Found replies to the request with a redirect to url.
func Found(w http.ResponseWriter, r *http.Request, url string) {
	Redirect(w, r, url, http.StatusFound)
}

// JSON writes the json-encoded error message to the response
// with a 400 bad request status code.
func JSON(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	_ = enc.Encode(v)
}

// OK writes the json-encoded data to the response.
func OK(w http.ResponseWriter, v any) {
	JSON(w, v, http.StatusOK)
}

// ErrorCode writes the json-encoded error message to the response.
func ErrorCode(w http.ResponseWriter, status int, v ...any) {
	var errStr string
	if len(v) > 0 {
		ev := v[0]
		if hookError != nil {
			switch vv := ev.(type) {
			case error:
				ev = hookError(vv)
			case string:
				ev = hookError(errors.New(vv))
			}
		}
		switch vv := ev.(type) {
		case error:
			errStr = vv.Error()
		case string:
			errStr = vv
		}
	}
	if errStr == "" {
		w.WriteHeader(status)
		return
	}

	logger.Error(errStr)
	JSON(w, errStr, status)
}

// InternalError writes the json-encoded error message to the response
// with a 500 internal server error.
func InternalError(w http.ResponseWriter, v ...any) {
	ErrorCode(w, http.StatusInternalServerError, v...)
}

// NotImplemented writes the json-encoded error message to the
// response with a 501 not found status code.
func NotImplemented(w http.ResponseWriter, v ...any) {
	ErrorCode(w, http.StatusNotImplemented, v...)
}

// NotFound writes the json-encoded error message to the response
// with a 404 not found status code.
func NotFound(w http.ResponseWriter, v ...any) {
	ErrorCode(w, http.StatusNotFound, v...)
}

// Unauthorized writes the json-encoded error message to the response
// with a 401 unauthorized status code.
func Unauthorized(w http.ResponseWriter, v ...any) {
	ErrorCode(w, http.StatusUnauthorized, v...)
}

// Forbidden writes the json-encoded error message to the response
// with a 403 forbidden status code.
func Forbidden(w http.ResponseWriter, v ...any) {
	ErrorCode(w, http.StatusForbidden, v...)
}

// BadRequest writes the json-encoded error message to the response
// with a 400 bad request status code.
func BadRequest(w http.ResponseWriter, v ...any) {
	ErrorCode(w, http.StatusBadRequest, v...)
}
