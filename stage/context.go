// Package stage

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

package stage

import (
	"context"
	"sync"
)

type Context interface {
	context.Context
	Ctx() context.Context
	WithCtx(ctx context.Context)
	WithValue(key, value interface{})
}

type valueCtx struct {
	context.Context
	set sync.Map
}

func NewCtx(ctx context.Context) Context {
	if ctx == nil {
		panic("cannot create stage context from nil context")
	}
	return &valueCtx{Context: ctx}
}

func (c *valueCtx) Ctx() context.Context {
	return c.Context
}

func (c *valueCtx) WithCtx(ctx context.Context) {
	c.Context = ctx
}

func (c *valueCtx) clone() Context {
	vc := &valueCtx{Context: c.Context}
	c.set.Range(func(key, value interface{}) bool {
		vc.set.Store(key, value)
		return true
	})
	return vc
}

func (c *valueCtx) combine(c2 Context) {
	if vc, ok := c2.(*valueCtx); ok {
		vc.set.Range(func(key, value interface{}) bool {
			c.set.Store(key, value)
			return true
		})
	}
}

func (c *valueCtx) WithValue(key, value interface{}) {
	c.set.Store(key, value)
}

func (c *valueCtx) Value(key interface{}) interface{} {
	value, ok := c.set.Load(key)
	if ok {
		return value
	}
	return c.Context.Value(key)
}

type contextKey struct {
	name string
}

func (c *contextKey) String() string {
	return "context value " + c.name
}

var (
	NameKey = &contextKey{name: "name"}
)

func ContextName(ctx Context) string {
	return ctx.Value(NameKey).(string)
}
