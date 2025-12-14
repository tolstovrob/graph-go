= Список смежности Iб
== Условие
Построить граф, полученный из исходного удалением висячих вершин.

== Код

```go
/*
 * This package contains algorithms and tasks for my SSU course
 */

package algo

import "github.com/tolstovrob/graph-go/graph"

/*
 * Task: Build graph obtained by removing pendant vertices from original graph
 */

func RemovePendantVertices(gr *graph.Graph) (*graph.Graph, error) {
	if gr.Nodes == nil {
		return nil, graph.ThrowNodesListIsNil()
	}

	newGraph := gr.Copy()

	removed := true
	for removed {
		removed = false
		for key := range newGraph.Nodes {
			degree := len(newGraph.AdjacencyMap[key])
			if degree == 1 {
				newGraph.RemoveNodeByKey(key)
				removed = true
				break // Restart iteration since map changed
			}
		}
	}

	return newGraph, nil
}

```
