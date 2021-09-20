// Package cycle

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

package cycle

// Graph is a directed acyclic graph.
type Graph map[string]*Endpoint

// Endpoint is an endpoint in the graph.
type Endpoint struct {
	Name   string
	relies []string
}

// New creates a new directed acyclic graph that can
// determinate if a stage has dependencies.
func New() Graph {
	return make(map[string]*Endpoint)
}

// Add establishes a dependency between two endpoints in the graph.
func (g Graph) Add(name string, relies ...string) *Endpoint {
	endpoint := &Endpoint{
		Name:   name,
		relies: relies,
	}
	g[name] = endpoint
	return endpoint
}

// Get returns the endpoint from the graph.
func (g Graph) Get(name string) (*Endpoint, bool) {
	endpoint, ok := g[name]
	return endpoint, ok
}

// Dependencies returns the direct dependencies accounting for
// skipped dependencies.
func (g Graph) Dependencies(name string) []string {
	endpoint := g[name]
	return g.dependencies(endpoint)
}

// dependencies returns the list of dependencies for the
// endpoint taking into account skipped dependencies.
func (g Graph) dependencies(parent *Endpoint) []string {
	if parent == nil {
		return nil
	}

	combined := make([]string, 0, len(parent.relies))
	for _, name := range parent.relies {
		endpoint, ok := g[name]
		if !ok {
			continue
		}
		combined = append(combined, endpoint.Name)
	}
	return combined
}

// Ancestors returns the all dependencies of the endpoint.
func (g Graph) Ancestors(name string) []*Endpoint {
	endpoint := g[name]
	return g.ancestors(endpoint)
}

// ancestors returns the list of all dependencies for the endpoint.
func (g Graph) ancestors(parent *Endpoint) []*Endpoint {
	if parent == nil {
		return nil
	}

	combined := make([]*Endpoint, 0, len(parent.relies))
	for _, name := range parent.relies {
		endpoint, ok := g[name]
		if !ok {
			continue
		}
		combined = append(combined, g.ancestors(endpoint)...)
	}
	return combined
}

// DetectCycles returns true if cycles are detected in the graph.
func (g Graph) DetectCycles() bool {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for name := range g {
		if !visited[name] {
			if g.detectCycles(name, visited, recStack) {
				return true
			}
		}
	}
	return false
}

// detectCycles returns true if the endpoint is cyclical.
func (g Graph) detectCycles(name string, visited, recStack map[string]bool) bool {
	visited[name] = true
	recStack[name] = true

	endpoint, ok := g[name]
	if !ok {
		return false
	}
	for _, v := range endpoint.relies {
		// only check cycles on an endpoint one time
		if !visited[v] {
			if g.detectCycles(v, visited, recStack) {
				return true
			}
			// if we've visited this endpoint in this recursion
			// stack, then we have a cycle
		} else if recStack[v] {
			return true
		}

	}
	recStack[name] = false
	return false
}
