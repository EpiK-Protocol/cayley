package graphlog

import (
	"bytes"
	"sort"

	"github.com/cayleygraph/quad"
	"github.com/epik-protocol/epik-gateway-backend/graph"
	"github.com/epik-protocol/epik-gateway-backend/graph/refs"
)

type Op interface {
	isOp()
}

var (
	_ Op = NodeUpdate{}
	_ Op = QuadUpdate{}
)

type NodeUpdate struct {
	Hash   refs.ValueHash
	Val    quad.Value
	RefInc int
}

func (NodeUpdate) isOp() {}

type QuadUpdate struct {
	Ind  int
	Quad refs.QuadHash
	Del  bool
}

func (QuadUpdate) isOp() {}

type Deltas struct {
	IncNode []NodeUpdate
	DecNode []NodeUpdate
	QuadAdd []QuadUpdate
	QuadDel []QuadUpdate
}

func InsertQuads(in []quad.Quad) *Deltas {
	hnodes := make(map[refs.ValueHash]*NodeUpdate, len(in)*2)
	quadAdd := make([]QuadUpdate, 0, len(in))
	for i, qd := range in {
		var q refs.QuadHash
		for _, dir := range quad.Directions {
			v := qd.Get(dir)
			if v == nil {
				continue
			}
			h := refs.HashOf(v)
			q.Set(dir, h)
			n := hnodes[h]
			if n == nil {
				n = &NodeUpdate{Hash: h, Val: v}
				hnodes[h] = n
			}
			n.RefInc++
		}
		quadAdd = append(quadAdd, QuadUpdate{Ind: i, Quad: q})
	}
	incNodes := make([]NodeUpdate, 0, len(hnodes))
	for _, n := range hnodes {
		incNodes = append(incNodes, *n)
	}
	hnodes = nil
	sort.Slice(incNodes, func(i, j int) bool {
		return bytes.Compare(incNodes[i].Hash[:], incNodes[j].Hash[:]) < 0
	})
	return &Deltas{
		IncNode: incNodes,
		QuadAdd: quadAdd,
	}
}

func SplitDeltas(in []graph.Delta) *Deltas {
	hnodes := make(map[refs.ValueHash]*NodeUpdate, len(in)*2)
	quadAdd := make([]QuadUpdate, 0, len(in))
	quadDel := make([]QuadUpdate, 0, len(in)/2)
	var nadd, ndel int
	for i, d := range in {
		dn := 0
		switch d.Action {
		case graph.Add:
			dn = +1
			nadd++
		case graph.Delete:
			dn = -1
			ndel++
		default:
			panic("unknown action")
		}
		var q refs.QuadHash
		if d.Action == graph.Add || len(d.Cid) == 0 {
			// nodes deleted by cid will be loaded at ApplyDeltas
			for _, dir := range quad.Directions {
				v := d.Quad.Get(dir)
				if v == nil {
					continue
				}
				h := refs.HashOf(v)
				q.Set(dir, h)
				n := hnodes[h]
				if n == nil {
					n = &NodeUpdate{Hash: h, Val: v}
					hnodes[h] = n
				}
				n.RefInc += dn
			}
		}
		u := QuadUpdate{Ind: i, Quad: q, Del: d.Action == graph.Delete}
		if !u.Del {
			quadAdd = append(quadAdd, u)
		} else {
			quadDel = append(quadDel, u)
		}
	}
	incNodes := make([]NodeUpdate, 0, nadd)
	decNodes := make([]NodeUpdate, 0, ndel)
	for _, n := range hnodes {
		if n.RefInc >= 0 {
			incNodes = append(incNodes, *n)
		} else {
			decNodes = append(decNodes, *n)
		}
	}
	sort.Slice(incNodes, func(i, j int) bool {
		return bytes.Compare(incNodes[i].Hash[:], incNodes[j].Hash[:]) < 0
	})
	sort.Slice(decNodes, func(i, j int) bool {
		return bytes.Compare(decNodes[i].Hash[:], decNodes[j].Hash[:]) < 0
	})
	hnodes = nil
	return &Deltas{
		IncNode: incNodes, DecNode: decNodes,
		QuadAdd: quadAdd, QuadDel: quadDel,
	}
}

func SplitEpikDeltas(in []graph.Delta) *Deltas {
	hnodes := make(map[refs.ValueHash]*NodeUpdate, len(in)*2)
	quadAdd := make([]QuadUpdate, 0, len(in))
	quadDel := make([]QuadUpdate, 0, len(in)/2)
	for i, d := range in {
		switch d.Action {
		case graph.Add:
			var q refs.QuadHash
			for _, dir := range quad.Directions {
				v := d.Quad.Get(dir)
				if v == nil {
					continue
				}
				h := refs.HashOf(v)
				q.Set(dir, h)
				n := hnodes[h]
				if n == nil {
					n = &NodeUpdate{Hash: h, Val: v}
					hnodes[h] = n
				}
				n.RefInc++
			}
			quadAdd = append(quadAdd, QuadUpdate{Ind: i, Quad: q})
		case graph.Delete:
			quadDel = append(quadDel, QuadUpdate{Ind: i, Del: true})
		default:
			panic("unknown action")
		}
	}
	incNodes := make([]NodeUpdate, 0, len(hnodes))
	for _, n := range hnodes {
		incNodes = append(incNodes, *n)
	}
	sort.Slice(incNodes, func(i, j int) bool {
		return bytes.Compare(incNodes[i].Hash[:], incNodes[j].Hash[:]) < 0
	})
	hnodes = nil
	return &Deltas{
		IncNode: incNodes,
		QuadAdd: quadAdd,
		QuadDel: quadDel,
	}
}
