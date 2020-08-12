// Copyright 2014 The Cayley Authors. All rights reserved.
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

package memstore

import (
	"context"
	"reflect"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cayleygraph/quad"
	"github.com/epik-protocol/epik-gateway-backend/graph"
	"github.com/epik-protocol/epik-gateway-backend/graph/graphtest"
	"github.com/epik-protocol/epik-gateway-backend/graph/iterator"
	"github.com/epik-protocol/epik-gateway-backend/graph/refs"
	"github.com/epik-protocol/epik-gateway-backend/query/shape"
	"github.com/epik-protocol/epik-gateway-backend/writer"
)

// This is a simple test graph.
//
//    +---+                        +---+
//    | A |-------               ->| F |<--
//    +---+       \------>+---+-/  +---+   \--+---+
//                 ------>|#B#|      |        | E |
//    +---+-------/      >+---+      |        +---+
//    | C |             /            v
//    +---+           -/           +---+
//      ----    +---+/             |#G#|
//          \-->|#D#|------------->+---+
//              +---+
//
var simpleGraph = []quad.Quad{
	quad.MakeRaw("A", "follows", "B", ""),
	quad.MakeRaw("C", "follows", "B", ""),
	quad.MakeRaw("C", "follows", "D", ""),
	quad.MakeRaw("D", "follows", "B", ""),
	quad.MakeRaw("B", "follows", "F", ""),
	quad.MakeRaw("F", "follows", "G", ""),
	quad.MakeRaw("D", "follows", "G", ""),
	quad.MakeRaw("E", "follows", "F", ""),
	quad.MakeRaw("B", "status", "cool", "status_graph"),
	quad.MakeRaw("D", "status", "cool", "status_graph"),
	quad.MakeRaw("G", "status", "cool", "status_graph"),
}

func makeTestStore(data []quad.Quad) (*QuadStore, graph.QuadWriter, []pair) {
	seen := make(map[string]struct{})
	qs := New()
	var (
		val int64
		ind []pair
	)
	writer, _ := writer.NewSingleReplication(qs, nil)
	for _, t := range data {
		for _, dir := range quad.Directions {
			qp := t.GetString(dir)
			if _, ok := seen[qp]; !ok && qp != "" {
				val++
				ind = append(ind, pair{qp, val})
				seen[qp] = struct{}{}
			}
		}

		writer.AddQuad(t)
		val++
	}
	return qs, writer, ind
}

func TestMemstore(t *testing.T) {
	graphtest.TestAll(t, func(t testing.TB) (graph.QuadStore, graph.Options, func()) {
		return New(), nil, func() {}
	}, &graphtest.Config{
		AlwaysRunIntegration: true,
	})
}

func BenchmarkMemstore(b *testing.B) {
	graphtest.BenchmarkAll(b, func(t testing.TB) (graph.QuadStore, graph.Options, func()) {
		return New(), nil, func() {}
	}, &graphtest.Config{
		AlwaysRunIntegration: true,
	})
}

type pair struct {
	query string
	value int64
}

func TestMemstoreValueOf(t *testing.T) {
	qs, _, index := makeTestStore(simpleGraph)
	exp := graph.Stats{
		Nodes: refs.Size{Value: 11, Exact: true},
		Quads: refs.Size{Value: 11, Exact: true},
	}
	st, err := qs.Stats(context.Background(), true)
	require.NoError(t, err)
	require.Equal(t, exp, st, "Unexpected quadstore size")

	for _, test := range index {
		v := qs.ValueOf(quad.Raw(test.query))
		switch v := v.(type) {
		default:
			t.Errorf("ValueOf(%q) returned unexpected type, got:%T expected int64", test.query, v)
		case bnode:
			require.Equal(t, test.value, int64(v))
		}
	}
}

func TestIteratorsAndNextResultOrderA(t *testing.T) {
	ctx := context.TODO()
	qs, _, _ := makeTestStore(simpleGraph)

	fixed := iterator.NewFixed()
	fixed.Add(qs.ValueOf(quad.Raw("C")))

	fixed2 := iterator.NewFixed()
	fixed2.Add(qs.ValueOf(quad.Raw("follows")))

	all := qs.NodesAllIterator()

	const allTag = "all"
	innerAnd := iterator.NewAnd(
		graph.NewLinksTo(qs, fixed2, quad.Predicate),
		graph.NewLinksTo(qs, iterator.Tag(all, allTag), quad.Object),
	)

	hasa := graph.NewHasA(qs, innerAnd, quad.Subject)
	outerAnd := iterator.NewAnd(fixed, hasa).Iterate()

	if !outerAnd.Next(ctx) {
		t.Error("Expected one matching subtree")
	}
	val := outerAnd.Result()
	if qs.NameOf(val) != quad.Raw("C") {
		t.Errorf("Matching subtree should be %s, got %s", "barak", qs.NameOf(val))
	}

	var (
		got    []string
		expect = []string{"B", "D"}
	)
	for {
		m := make(map[string]graph.Ref, 1)
		outerAnd.TagResults(m)
		got = append(got, quad.ToString(qs.NameOf(m[allTag])))
		if !outerAnd.NextPath(ctx) {
			break
		}
	}
	sort.Strings(got)

	if !reflect.DeepEqual(got, expect) {
		t.Errorf("Unexpected result, got:%q expect:%q", got, expect)
	}

	if outerAnd.Next(ctx) {
		t.Error("More than one possible top level output?")
	}
}

func TestLinksToOptimization(t *testing.T) {
	qs, _, _ := makeTestStore(simpleGraph)

	lto := shape.BuildIterator(context.TODO(), qs, shape.Quads{
		{Dir: quad.Object, Values: shape.Lookup{quad.Raw("cool")}},
	})

	newIt, changed := lto.Optimize(context.TODO())
	if changed {
		t.Errorf("unexpected optimization step")
	}
	if _, ok := newIt.(*Iterator); !ok {
		t.Fatal("Didn't swap out to LLRB")
	}
}

func TestRemoveQuad(t *testing.T) {
	ctx := context.TODO()
	qs, w, _ := makeTestStore(simpleGraph)

	err := w.RemoveQuad(quad.Make(
		"E",
		"follows",
		"F",
		nil,
	))

	if err != nil {
		t.Error("Couldn't remove quad", err)
	}

	fixed := iterator.NewFixed()
	fixed.Add(qs.ValueOf(quad.Raw("E")))

	fixed2 := iterator.NewFixed()
	fixed2.Add(qs.ValueOf(quad.Raw("follows")))

	innerAnd := iterator.NewAnd(
		graph.NewLinksTo(qs, fixed, quad.Subject),
		graph.NewLinksTo(qs, fixed2, quad.Predicate),
	)

	hasa := graph.NewHasA(qs, innerAnd, quad.Object)

	newIt, _ := hasa.Optimize(ctx)
	if newIt.Iterate().Next(ctx) {
		t.Error("E should not have any followers.")
	}
}

func TestTransaction(t *testing.T) {
	qs, w, _ := makeTestStore(simpleGraph)
	st, err := qs.Stats(context.Background(), true)
	require.NoError(t, err)

	tx := graph.NewTransaction()
	tx.AddQuad(quad.Make(
		"E",
		"follows",
		"G",
		nil))
	tx.RemoveQuad(quad.Make(
		"Non",
		"existent",
		"quad",
		nil))

	err = w.ApplyTransaction(tx)
	if err == nil {
		t.Error("Able to remove a non-existent quad")
	}
	st2, err := qs.Stats(context.Background(), true)
	require.NoError(t, err)
	require.Equal(t, st, st2, "Appended a new quad in a failed transaction")
}

func TestApplyDeltas(t *testing.T) {

	var epoch1Deltas = []graph.Delta{
		{ // duplicated quad
			Cid:    "e6e5b003b3ce13d939e6",
			Quad:   quad.MakeRaw("A", "follows", "B", ""),
			Action: graph.Add,
		},
		{ // new quad
			Cid:    "5ea85b3acc794d9ed651",
			Quad:   quad.MakeRaw("A", "follows", "C", ""),
			Action: graph.Add,
		},
	}

	qs, _, _ := makeTestStore(simpleGraph)
	exp := graph.Stats{
		Nodes: refs.Size{Value: 11, Exact: true},
		Quads: refs.Size{Value: 11, Exact: true},
	}
	st, err := qs.Stats(context.Background(), true)
	require.NoError(t, err)
	require.Equal(t, exp, st, "Unexpected quadstore size")
	require.Equal(t, len(qs.cqIndex.index), 0)

	err = qs.ApplyDeltas(0, epoch1Deltas, graph.IgnoreOpts{IgnoreDup: true, IgnoreMissing: true})
	rerr := err.(*graph.DeltaError)
	require.EqualError(t, rerr.Err, graph.ErrInvalidCid.Error())

	// add epoch1Deltas
	err = qs.ApplyDeltas(1, epoch1Deltas, graph.IgnoreOpts{IgnoreDup: true, IgnoreMissing: true})
	require.NoError(t, err)
	exp.Quads.Value = 12
	exp.Epoch = 1
	st, err = qs.Stats(context.Background(), true)
	require.NoError(t, err)
	require.Equal(t, exp, st, "Unexpected quadstore size")
	require.Equal(t, len(qs.cqIndex.index), len(epoch1Deltas))

	// delete by cid ---- epoch1Deltas[1]
	err = qs.ApplyDeltas(1, []graph.Delta{
		{
			Cid:    epoch1Deltas[1].Cid,
			Action: graph.Delete,
		},
	}, graph.IgnoreOpts{IgnoreDup: true, IgnoreMissing: true})
	require.NoError(t, err)
	require.Equal(t, len(qs.cqIndex.index), 1)
	_, ok := qs.cqIndex.index[epoch1Deltas[0].Cid]
	require.True(t, ok, "Unexpected cid")

	exp.Quads.Value = 11
	st, err = qs.Stats(context.Background(), true)
	require.NoError(t, err)
	require.Equal(t, exp, st, "Unexpected quadstore size")

	// delete by quad, same with epoch1Deltas[0]
	err = qs.ApplyDeltas(0, []graph.Delta{
		{
			Quad:   simpleGraph[0],
			Action: graph.Delete,
		},
	}, graph.IgnoreOpts{IgnoreDup: true, IgnoreMissing: true})
	require.NoError(t, err)
	exp = graph.Stats{
		Nodes: refs.Size{Value: 10, Exact: true}, // "A" removed
		Quads: refs.Size{Value: 10, Exact: true}, // simpleGraph[0] removed
		Epoch: 1,
	}
	st, err = qs.Stats(context.Background(), true)
	require.NoError(t, err)
	require.Equal(t, exp, st, "Unexpected quadstore size")
	require.Equal(t, len(qs.cqIndex.index), 1)
	_, ok = qs.cqIndex.index[epoch1Deltas[0].Cid]
	require.True(t, ok, "Unexpected cid")

	// delete invalid cid index
	err = qs.ApplyDeltas(2, []graph.Delta{
		{
			Cid:    epoch1Deltas[0].Cid,
			Action: graph.Delete,
		},
	}, graph.IgnoreOpts{IgnoreDup: true, IgnoreMissing: false})
	require.NoError(t, err)
	require.Equal(t, len(qs.cqIndex.index), 0)
	exp.Epoch = 2
	st, err = qs.Stats(context.Background(), true)
	require.NoError(t, err)
	require.Equal(t, exp, st, "Unexpected quadstore size")
}
