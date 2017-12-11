/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cloudgraph

import (
	"io"
	"sync"
)

type LazyGraph struct {
	LoadingFunc func() GraphAPI
	once        sync.Once
	api         GraphAPI
}

func (g *LazyGraph) load() {
	g.once.Do(func() {
		g.api = g.LoadingFunc()
	})
}

func (g *LazyGraph) Find(q Query) ([]Resource, error) {
	g.load()
	return g.api.Find(q)
}

func (g *LazyGraph) FindOne(q Query) (Resource, error) {
	g.load()
	return g.api.FindOne(q)
}

func (g *LazyGraph) MarshalTo(w io.Writer) error {
	g.load()
	return g.api.MarshalTo(w)
}
