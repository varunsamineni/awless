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

import "io"

type GraphAPI interface {
	Find(Query) ([]Resource, error)
	FindOne(Query) (Resource, error)
	Accept(v Visitor) error
	MarshalTo(w io.Writer) error
}

type Visitor interface {
	Visit(GraphAPI) error
}

type Resource interface {
	Type() string
	Id() string
	Properties() map[string]interface{}
	Property(string) (interface{}, bool)
	Meta(string) (interface{}, bool)
	Relations(string) []Resource
	Same(Resource) bool
}

func NewQuery(resourceType ...string) Query {
	return Query{ResourceType: resourceType}
}

func (q Query) Property(name string, value interface{}) Query {
	q.PropertyValues = append(q.PropertyValues, propertyValue{Name: name, Value: value})
	return q
}

type Query struct {
	ResourceType   []string
	PropertyValues []propertyValue
}

type propertyValue struct {
	Name  string
	Value interface{}
}
