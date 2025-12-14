= Список смежности III
== Условие
Дан взвешенный неориентированный граф из N вершин и M ребер. Требуется найти в нем каркас минимального веса. Алгоритм Прима


== Код

```go
/*
 * This package contains algorithms and tasks for my SSU course
 */

package algo

import "github.com/tolstovrob/graph-go/graph"

/*
 * Task: Find Minimum Spanning Tree using Prim's algorithm
 *
 * Minimum Spanning Tree (MST) - a tree that connects all vertices with minimum total edge weight
 * Prim's Algorithm - greedy algorithm that grows the MST one vertex at a time
 */

// MSTResult represents the result of Minimum Spanning Tree calculation
type MSTResult struct {
	TotalWeight graph.TWeight // Total weight of all edges in MST
	Edges       []*graph.Edge // List of edges that form the MST
	IsPossible  bool          // Whether MST construction is possible (graph must be connected)
}

// FindMSTPrim finds Minimum Spanning Tree using Prim's algorithm
// Time Complexity: O(V^2) for this implementation, can be optimized to O(E log V) with priority queue
// Space Complexity: O(V + E)
func FindMSTPrim(gr *graph.Graph) (*MSTResult, error) {
	// Input validation: check if graph nodes exist
	if gr.Nodes == nil {
		return nil, graph.ThrowNodesListIsNil()
	}

	// Base case: empty graph is trivially a tree
	if len(gr.Nodes) == 0 {
		return &MSTResult{
			TotalWeight: 0,
			Edges:       []*graph.Edge{},
			IsPossible:  true,
		}, nil
	}

	// MST Requirement: graph must be connected
	// A disconnected graph cannot have a spanning tree that connects all vertices
	if !gr.IsConnected() {
		return &MSTResult{
			TotalWeight: 0,
			Edges:       []*graph.Edge{},
			IsPossible:  false, // MST not possible for disconnected graphs
		}, nil
	}

	// Prim's algorithm works on undirected graphs
	// If graph is directed, create an undirected copy for MST calculation
	if gr.Options.IsDirected {
		tempGraph := gr.Copy()
		tempGraph.UpdateGraph(graph.WithGraphDirected(false))
		return findMSTPrimInternal(tempGraph)
	}

	// For undirected graphs, proceed with normal MST calculation
	return findMSTPrimInternal(gr)
}

// findMSTPrimInternal implements the core Prim's algorithm logic
// Algorithm Strategy: Grow the MST by repeatedly adding the cheapest edge
// that connects a vertex in MST to a vertex outside MST
func findMSTPrimInternal(gr *graph.Graph) (*MSTResult, error) {
	// Handle empty graph case
	if len(gr.Nodes) == 0 {
		return &MSTResult{
			TotalWeight: 0,
			Edges:       []*graph.Edge{},
			IsPossible:  true,
		}, nil
	}

	// Data Structures for Algorithm:
	// inMST: tracks which vertices are already included in the MST
	// mstEdges: stores the edges that form the MST
	inMST := make(map[graph.TKey]bool)
	var mstEdges []*graph.Edge

	// Step 1: Initialize with any vertex
	// Prim's algorithm can start from any vertex - choice doesn't affect result
	var firstVertex graph.TKey
	for vertex := range gr.Nodes {
		firstVertex = vertex
		break // Pick the first available vertex
		// I do love working with queues in go
	}
	inMST[firstVertex] = true // Mark first vertex as included

	// Step 2: Repeat until all vertices are in MST
	// MST must contain exactly V-1 edges for V vertices
	for len(inMST) < len(gr.Nodes) {
		var bestEdge *graph.Edge             // The edge with minimum weight
		bestWeight := graph.TWeight(1 << 30) // Initialize with "infinity"

		// Step 2.1: Find minimum weight edge connecting MST to non-MST vertices
		// Strategy: Check all edges from vertices inside MST to their neighbors outside MST
		for u := range inMST {
			// Explore all neighbors of vertex u (which is already in MST)
			for _, neighbor := range gr.AdjacencyMap[u] {
				// Only consider neighbors that are NOT yet in MST
				if !inMST[neighbor] {
					// Find the actual edge between u and neighbor
					edge := getEdgeBetweenReliable(gr, u, neighbor)

					// If edge exists and has lower weight than current best, update best edge
					if edge != nil && edge.Weight < bestWeight {
						bestWeight = edge.Weight
						bestEdge = edge
					}
				}
			}
		}

		if bestEdge == nil {
			return &MSTResult{
				TotalWeight: 0,
				Edges:       []*graph.Edge{},
				IsPossible:  false,
			}, nil
		}

		mstEdges = append(mstEdges, bestEdge)

		if inMST[bestEdge.Source] {
			inMST[bestEdge.Destination] = true
		} else {
			inMST[bestEdge.Source] = true
		}

		// At this point, MST has grown by one vertex and one edge
		// The algorithm maintains the invariant that MST is always a tree
	}

	// Step 3: Calculate total weight of MST
	totalWeight := graph.TWeight(0)
	for _, edge := range mstEdges {
		totalWeight += edge.Weight
	}

	return &MSTResult{
		TotalWeight: totalWeight,
		Edges:       mstEdges,
		IsPossible:  true,
	}, nil
}

// getEdgeBetweenReliable finds an edge between two vertices in the graph
// Handles both directed and undirected graphs correctly
func getEdgeBetweenReliable(gr *graph.Graph, u, v graph.TKey) *graph.Edge {
	// First, search for edge u → v (forward direction)
	for _, edge := range gr.Edges {
		if edge.Source == u && edge.Destination == v {
			return edge
		}
	}

	// For undirected graphs, also search for edge v → u (reverse direction)
	// In undirected graphs, edge A-B is the same as B-A
	if !gr.Options.IsDirected {
		for _, edge := range gr.Edges {
			if edge.Source == v && edge.Destination == u {
				return edge
			}
		}
	}

	// No edge found between u and v
	return nil
}

```