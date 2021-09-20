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
	"errors"
	"fmt"

	"github.com/99nil/go/cycle"
	"github.com/99nil/go/sets"
)

type InstanceFunc func(ctx context.Context) error

type Instance struct {
	name string
	// Whether to enable asynchronous processing.
	async bool
	// Stage subset.
	cs []*Instance
	// Calling method before executing cs.
	pre InstanceFunc
	// Calling method after executing cs.
	sub InstanceFunc
	// The names of other stages that need to be relied upon before execution.
	relies []string
}

func New(name string) *Instance {
	return &Instance{name: name}
}

func (ins *Instance) rename() {
	ns := sets.NewString()
	for _, c := range ins.cs {
		i := 0
		for ns.Has(c.name) {
			c.name = fmt.Sprintf("%s_%d", c.name, i)
			i++
		}
		ns.Add(c.name)
	}
}

// Add adds subsets
func (ins *Instance) Add(cs ...*Instance) *Instance {
	for _, c := range cs {
		ins.cs = append(ins.cs, c)
	}
	return ins
}

// SetAsync sets whether the current stage is executed asynchronously.
func (ins *Instance) SetAsync(async bool) *Instance {
	ins.async = async
	return ins
}

// SetPreFunc sets the execution method before executing the subset.
func (ins *Instance) SetPreFunc(f InstanceFunc) *Instance {
	ins.pre = f
	return ins
}

// SetSubFunc sets the execution method after executing the subset.
func (ins *Instance) SetSubFunc(f InstanceFunc) *Instance {
	ins.sub = f
	return ins
}

// SetRely sets the names of other stages that the current stage needs to depend on.
func (ins *Instance) SetRely(names ...string) *Instance {
	ins.relies = make([]string, 0, len(names))
	ins.relies = append(ins.relies, names...)
	return ins
}

// getChildNames gets the names of all subsets.
func (ins *Instance) getChildNames() []string {
	csLen := len(ins.cs)
	if csLen == 0 {
		return nil
	}
	res := make([]string, 0, csLen)
	for _, c := range ins.cs {
		res = append(res, c.name)
	}
	return res
}

// hasLoop checks whether there is a circular dependency.
func (ins *Instance) hasLoop() bool {
	graph := cycle.New()
	for _, c := range ins.cs {
		graph.Add(c.name, c.relies...)
	}
	return graph.DetectCycles()
}

func (ins *Instance) Run(ctx context.Context) error {
	if ins.hasLoop() {
		return errors.New("dependency cycle detected")
	}
	// TODO Check for non-existent dependencies.

	ctx = contextWithName(ctx, ins.name)
	if ins.pre != nil {
		if err := ins.pre(ctx); err != nil {
			return err
		}
	}

	var err error
	if ins.async {
		err = ins.runAsync(ctx)
	} else {
		err = ins.runSync(ctx)
	}
	if err != nil {
		return err
	}

	if ins.sub != nil {
		if err := ins.sub(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (ins *Instance) runSync(ctx context.Context) error {
	doneSet := sets.NewString()
	pending := ins.cs[:]
	for len(pending) > 0 {
		wait := make([]*Instance, 0, len(pending))

		for _, c := range pending {
			if doneSet.Has(c.name) {
				continue
			}
			if len(c.relies) != 0 {
				// Check whether the dependency has completed running.
				if !doneSet.HasAll(c.relies...) {
					wait = append(wait, c)
					continue
				}
			}
			if err := c.Run(ctx); err != nil {
				return err
			}
			doneSet.Add(c.name)
		}

		pending = wait[:]
	}
	return nil
}

func (ins *Instance) runAsync(ctx context.Context) error {
	doneSet := sets.NewString()
	pending := ins.cs[:]
	for len(pending) > 0 {
		wait := make([]*Instance, 0, len(pending))

		errCh := make(chan error)
		for _, c := range pending {
			if doneSet.Has(c.name) {
				continue
			}
			if len(c.relies) != 0 {
				// Check whether the dependency has completed running.
				if !doneSet.HasAll(c.relies...) {
					wait = append(wait, c)
					continue
				}
			}

			go func(c *Instance) {
				if err := c.Run(ctx); err != nil {
					errCh <- err
				}
			}(c)
			doneSet.Add(c.name)
		}
		if err := <-errCh; err != nil {
			return err
		}
		pending = wait[:]
	}
	return nil
}
