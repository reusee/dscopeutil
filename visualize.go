package dscopeutil

import (
	"io"
	"reflect"
	"sync"

	"github.com/emicklei/dot"
	"github.com/reusee/dscope"
)

func Visualize(s dscope.Scope, w io.Writer) error {
	s = s.Fork(dscope.DebugDefs)
	var debugInfo dscope.DebugInfo
	s.Assign(&debugInfo)

	edges := make(map[[2]reflect.Type]bool)
	for _, info := range debugInfo.Values {
		for _, defType := range info.DefTypes {
			if defType.Kind() != reflect.Func {
				continue
			}
			for i := 0; i < defType.NumIn(); i++ {
				in := defType.In(i)
				for j := 0; j < defType.NumOut(); j++ {
					out := defType.Out(j)
					edges[[2]reflect.Type{out, in}] = true
				}
			}
		}
	}

	g := dot.NewGraph(dot.Directed)
	var nodes sync.Map
	getNode := func(t reflect.Type) dot.Node {
		if v, ok := nodes.Load(t); ok {
			return v.(dot.Node)
		}
		node := g.Node(t.String())
		nodes.Store(t, node)
		return node
	}

	for edge := range edges {
		g.Edge(
			getNode(edge[0]),
			getNode(edge[1]),
		)
	}

	_, err := w.Write([]byte(g.String()))
	if err != nil {
		return we(err)
	}

	return nil
}
