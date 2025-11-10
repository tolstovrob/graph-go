/*
 * This package contains algorithms and tasks for my SSU course
 */

package algo

import "github.com/tolstovrob/graph-go/graph"

/*
 * Task: Check if there exists a vertex that can be removed to make the graph a tree
 */

func CanRemoveVertexToMakeTree(gr *graph.Graph) (bool, []graph.TKey, error) {
	if gr.Nodes == nil {
		return false, nil, graph.ThrowNodesListIsNil()
	}

	var candidates []graph.TKey

	// For each vertex, check if removing it makes the graph a tree
	for key := range gr.Nodes {
		// Create a copy and remove the vertex
		tempGraph := gr.Copy()
		if err := tempGraph.RemoveNodeByKey(key); err != nil {
			continue
		}

		// Check if the resulting graph is a tree
		if tempGraph.IsTree() {
			candidates = append(candidates, key)
		}
	}

	return len(candidates) > 0, candidates, nil
}
