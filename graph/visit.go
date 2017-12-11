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

package graph

import (
	"fmt"

	"github.com/wallix/awless/cloud/graph"
	"github.com/wallix/awless/cloud/properties"
	"github.com/wallix/awless/cloud/rdf"
	tstore "github.com/wallix/triplestore"
)

type Visitor interface {
	Visit(*Graph) error
}

type visitEachFunc func(res cloudgraph.Resource, depth int) error

func VisitorCollectFunc(collect *[]cloudgraph.Resource) visitEachFunc {
	return func(res cloudgraph.Resource, depth int) error {
		*collect = append(*collect, res)
		return nil
	}
}

type ParentsVisitor struct {
	From        cloudgraph.Resource
	Each        visitEachFunc
	IncludeFrom bool
}

func (v *ParentsVisitor) Visit(g cloudgraph.GraphAPI) error {
	startNode, foreach, err := prepareRDFVisit(g, v.From, v.Each, v.IncludeFrom)
	if err != nil {
		return err
	}
	rdfG, ok := g.(*Graph)
	if !ok {
		return fmt.Errorf("graph is not a RDF graph and thus can not visited with ChildrenVisitor")
	}
	return tstore.NewTree(rdfG.store.Snapshot(), rdf.ParentOf).TraverseAncestors(startNode, foreach)
}

type ChildrenVisitor struct {
	From        cloudgraph.Resource
	Each        visitEachFunc
	IncludeFrom bool
}

func (v *ChildrenVisitor) Visit(g cloudgraph.GraphAPI) error {
	startNode, foreach, err := prepareRDFVisit(g, v.From, v.Each, v.IncludeFrom)
	if err != nil {
		return err
	}
	rdfG, ok := g.(*Graph)
	if !ok {
		return fmt.Errorf("graph is not a RDF graph and thus can not visited with ChildrenVisitor")
	}
	return tstore.NewTree(rdfG.store.Snapshot(), rdf.ParentOf).TraverseDFS(startNode, foreach)
}

type SiblingsVisitor struct {
	From        *Resource
	Each        visitEachFunc
	IncludeFrom bool
}

func (v *SiblingsVisitor) Visit(g *Graph) error {
	startNode, foreach, err := prepareRDFVisit(g, v.From, v.Each, v.IncludeFrom)
	if err != nil {
		return err
	}

	return tstore.NewTree(g.store.Snapshot(), rdf.ParentOf).TraverseSiblings(startNode, resolveResourceType, foreach)
}

func prepareRDFVisit(g cloudgraph.GraphAPI, root cloudgraph.Resource, each visitEachFunc, includeRoot bool) (string, func(g tstore.RDFGraph, n string, i int) error, error) {
	rootNode := root.Id()

	foreach := func(rdfG tstore.RDFGraph, n string, i int) error {
		rT, err := resolveResourceType(rdfG, n)
		if err != nil {
			return err
		}
		res, err := g.FindOne(cloudgraph.NewQuery(rT).Property(properties.ID, n))
		if err != nil {
			return err
		}
		if includeRoot || !root.Same(res) {
			if err := each(res, i); err != nil {
				return err
			}
		}
		return nil
	}
	return rootNode, foreach, nil
}
