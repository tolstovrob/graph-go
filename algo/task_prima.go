/*
 * This package contains algorithms and tasks for my SSU course
 */

package algo

import "github.com/tolstovrob/graph-go/graph"

/*
 * Task: Find Minimum Spanning Tree using Prim's algorithm
 */

type MSTResult struct {
	TotalWeight graph.TWeight
	Edges       []*graph.Edge
	IsPossible  bool
}

func FindMSTPrim(gr *graph.Graph) (*MSTResult, error) {
	if gr.Nodes == nil {
		return nil, graph.ThrowNodesListIsNil()
	}

	// Check if graph is connected
	if !gr.IsConnected() {
		return &MSTResult{
			TotalWeight: 0,
			Edges:       []*graph.Edge{},
			IsPossible:  false,
		}, nil
	}

	// For Prim's algorithm, we need an undirected graph
	if gr.Options.IsDirected {
		// Create undirected copy for MST calculation
		tempGraph := gr.Copy()
		tempGraph.UpdateGraph(graph.WithGraphDirected(false))
		return findMSTPrimInternal(tempGraph)
	}

	return findMSTPrimInternal(gr)
}

func findMSTPrimInternal(gr *graph.Graph) (*MSTResult, error) {
	if len(gr.Nodes) == 0 {
		return &MSTResult{
			TotalWeight: 0,
			Edges:       []*graph.Edge{},
			IsPossible:  true,
		}, nil
	}

	inMST := make(map[graph.TKey]bool)
	minEdge := make(map[graph.TKey]*graph.Edge)
	minWeight := make(map[graph.TKey]graph.TWeight)

	for key := range gr.Nodes {
		minWeight[key] = ^graph.TWeight(0) // Max value
	}

	var startKey graph.TKey
	for key := range gr.Nodes {
		startKey = key
		break
	}
	minWeight[startKey] = 0

	for range gr.Nodes {
		// Find vertex with minimum weight not yet in MST
		currentKey := findMinKey(minWeight, inMST)
		if currentKey == 0 { // 0 indicates no valid key found
			break
		}

		inMST[currentKey] = true

		// Update adjacent vertices
		for _, neighbor := range gr.AdjacencyMap[currentKey] {
			if !inMST[neighbor] {
				// Find edge weight between currentKey and neighbor
				weight := getEdgeWeight(gr, currentKey, neighbor)
				if weight < minWeight[neighbor] {
					minWeight[neighbor] = weight
					minEdge[neighbor] = getEdgeBetween(gr, currentKey, neighbor)
				}
			}
		}
	}

	result := &MSTResult{
		Edges:      []*graph.Edge{},
		IsPossible: true,
	}

	// Collect MST edges (skip the starting node)
	for key, edge := range minEdge {
		if key != startKey && edge != nil {
			result.Edges = append(result.Edges, edge)
			result.TotalWeight += edge.Weight
		}
	}

	return result, nil
}

func findMinKey(weights map[graph.TKey]graph.TWeight, inMST map[graph.TKey]bool) graph.TKey {
	minWeight := ^graph.TWeight(0)
	var minKey graph.TKey

	for key, weight := range weights {
		if !inMST[key] && weight < minWeight {
			minWeight = weight
			minKey = key
		}
	}

	return minKey
}

func getEdgeWeight(gr *graph.Graph, u, v graph.TKey) graph.TWeight {
	for _, edge := range gr.Edges {
		if (edge.Source == u && edge.Destination == v) ||
			(!gr.Options.IsDirected && edge.Source == v && edge.Destination == u) {
			return edge.Weight
		}
	}
	return ^graph.TWeight(0) // Max weight if no edge found
}

func getEdgeBetween(gr *graph.Graph, u, v graph.TKey) *graph.Edge {
	for _, edge := range gr.Edges {
		if (edge.Source == u && edge.Destination == v) ||
			(!gr.Options.IsDirected && edge.Source == v && edge.Destination == u) {
			return edge
		}
	}
	return nil
}
