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

package database

import (
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
	page, _ := strconv.Atoi(
		r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(
		r.URL.Query().Get("size"))
	return PaginateDirect(page, size)
}

func PaginateDirect(page int, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page < 1 {
			page = 1
		}
		switch {
		case size > 100:
			size = 100
		case size < 1:
			size = 10
		}
		return db.Offset((page - 1) * size).Limit(size)
	}
}
