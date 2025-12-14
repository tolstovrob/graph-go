= Список смежности Ia
== Условие
Вывести те вершины, полустепень захода которых меньше, чем у заданной вершины.

Для данной вершины орграфа вывести все «заходящие» соседние вершины.

== Код

```go
/*
 * This package contains algorithms and tasks for my SSU course
 */

package algo

import "github.com/tolstovrob/graph-go/graph"

/*
 * Task 1: Get all nodes, for which degree is greater then half-degree of entry
 */

func InDegreeLessThan(gr *graph.Graph, targetKey graph.TKey) []graph.TKey {
	// Find all entries
	inDegree := make(map[graph.TKey]int)
	for _, edge := range gr.Edges {
		inDegree[edge.Destination]++
		if !gr.Options.IsDirected {
			inDegree[edge.Source]++
		}
	}
	targetInDegree := inDegree[targetKey]

	// Find all nodes satisfies task objective
	var result []graph.TKey
	for nodeKey := range gr.Nodes {
		if inDegree[nodeKey] < targetInDegree {
			result = append(result, nodeKey)
		}
	}
	return result
}

/*
 * Task 2: For directed graph node output all in-nodes
 */

func InNodesInDirected(gr *graph.Graph, targetKey graph.TKey) ([]graph.TKey, error) {
	// Check if graph is directed
	if !gr.Options.IsDirected {
		return nil, graph.ThrowGraphNotDirected()
	}

	// If graph directed, find all entries
	inNodes := []graph.TKey{}
	for _, edge := range gr.Edges {
		if edge.Destination == targetKey {
			inNodes = append(inNodes, edge.Source)
		}
	}
	return inNodes, nil
}

```
