= Список смежности IVa
== Условие
Эксцентриситет вершины — максимальное расстояние из всех минимальных расстояний от других вершин до данной вершины. Найти радиус графа — минимальный из эксцентриситетов его вершин. Алгоритм Дейкстры.

== Код

```go
/*
 * This package contains algorithms and tasks for my SSU course
 */

package algo

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tolstovrob/graph-go/graph"
)

/*
 * Task: Find shortest paths between all pairs of vertices using Floyd-Warshall algorithm
 *
 * Floyd-Warshall Algorithm - dynamic programming algorithm that finds shortest paths
 * between all pairs of vertices in a weighted graph
 * Can handle negative weights but not negative cycles
 */

// AllPairsShortestPath represents the result of Floyd-Warshall algorithm
type AllPairsShortestPath struct {
	Distances map[graph.TKey]map[graph.TKey]int64      // Shortest distance between every pair
	Next      map[graph.TKey]map[graph.TKey]graph.TKey // Next vertex in shortest path
	IsValid   bool                                     // Whether result is valid (no negative cycles)
	Message   string                                   // Status message about computation
}

// FindAllPairsShortestPath finds shortest paths between all vertex pairs using Floyd-Warshall
// Time Complexity: O(V^3) where V is number of vertices
// Can handle: directed/undirected graphs, negative weights (but not negative cycles)
func FindAllPairsShortestPath(gr *graph.Graph) (*AllPairsShortestPath, error) {
	if gr.Nodes == nil {
		return nil, graph.ThrowNodesListIsNil()
	}

	// Handle empty graph case
	if len(gr.Nodes) == 0 {
		return &AllPairsShortestPath{
			Distances: make(map[graph.TKey]map[graph.TKey]int64),
			Next:      make(map[graph.TKey]map[graph.TKey]graph.TKey),
			IsValid:   true,
			Message:   "Graph is empty",
		}, nil
	}

	return findFloydWarshall(gr)
}

// findFloydWarshall implements the core Floyd-Warshall algorithm
// Algorithm Strategy: Dynamic Programming - gradually improve shortest path estimates
// by considering each vertex as an intermediate point
func findFloydWarshall(gr *graph.Graph) (*AllPairsShortestPath, error) {
	// Step 1: Get sorted node keys for consistent processing
	// This ensures we always process vertices in the same order
	keys := getSortedKeys(gr.Nodes)

	// Step 2: Initialize distance and next matrices
	// dist[i][j] = shortest distance from i to j
	// next[i][j] = next vertex after i in shortest path to j
	dist := make(map[graph.TKey]map[graph.TKey]int64)
	next := make(map[graph.TKey]map[graph.TKey]graph.TKey)

	infinity := int64(1 << 30)

	// Step 3: Initialize matrices with base cases
	for _, i := range keys {
		dist[i] = make(map[graph.TKey]int64)
		next[i] = make(map[graph.TKey]graph.TKey)

		for _, j := range keys {
			if i == j {
				dist[i][j] = 0
			} else {
				dist[i][j] = infinity
			}
			next[i][j] = 0
		}
	}

	// Step 4: Initialize with direct edges
	// Set distances for edges that exist directly in the graph
	for _, edge := range gr.Edges {
		weight := int64(edge.Weight)

		// Set direct edge distance if it's better than current value
		if weight < dist[edge.Source][edge.Destination] {
			dist[edge.Source][edge.Destination] = weight
			next[edge.Source][edge.Destination] = edge.Destination
		}

		// For undirected graphs, set reverse edge as well
		if !gr.Options.IsDirected {
			if weight < dist[edge.Destination][edge.Source] {
				dist[edge.Destination][edge.Source] = weight
				next[edge.Destination][edge.Source] = edge.Source
			}
		}
	}

	// Step 5: Floyd-Warshall Algorithm
	// Consider each vertex as an intermediate point
	for _, k := range keys {
		for _, i := range keys { // Source vertex
			if dist[i][k] == infinity {
				continue
			}

			for _, j := range keys { // Destination vertex
				if dist[k][j] == infinity {
					continue
				}

				if dist[i][k]+dist[k][j] < dist[i][j] {
					dist[i][j] = dist[i][k] + dist[k][j]
					next[i][j] = next[i][k] // Path goes through k next
				}
			}
		}
	}

	// Step 6: Check for negative weight cycles
	// Negative cycle exists if any dist[i][i] < 0 (distance to self becomes negative)
	for _, k := range keys {
		if dist[k][k] < 0 {
			return &AllPairsShortestPath{
				IsValid: false,
				Message: "Graph contains negative weight cycles",
			}, nil
		}
	}

	// Step 7: Return successful result
	return &AllPairsShortestPath{
		Distances: dist,
		Next:      next,
		IsValid:   true,
		Message:   fmt.Sprintf("Computed shortest paths for %d vertices", len(keys)),
	}, nil
}

// GetPath reconstructs the shortest path from start to end using the next matrix
// Returns the sequence of vertices in the shortest path
func (apsp *AllPairsShortestPath) GetPath(start, end graph.TKey) []graph.TKey {
	// Check if no path exists
	if apsp.Next[start][end] == 0 {
		return nil
	}

	// Reconstruct path by following next pointers
	path := []graph.TKey{start}
	current := start

	for current != end {
		current = apsp.Next[current][end]
		path = append(path, current)
	}

	return path
}

// FormatDistanceMatrix creates a formatted string representation of the distance matrix
// Useful for displaying results in CLI
func (apsp *AllPairsShortestPath) FormatDistanceMatrix(gr *graph.Graph) string {
	if !apsp.IsValid {
		return apsp.Message
	}

	var sb strings.Builder
	keys := getSortedKeys(gr.Nodes)

	// Header information
	sb.WriteString("SHORTEST PATH DISTANCES BETWEEN ALL PAIRS OF VERTICES\n\n")
	sb.WriteString("Algorithm: Floyd-Warshall\n")
	sb.WriteString(fmt.Sprintf("Total vertices: %d\n", len(keys)))
	sb.WriteString(fmt.Sprintf("Total edges: %d\n", len(gr.Edges)))
	sb.WriteString(fmt.Sprintf("Directed: %v\n\n", gr.Options.IsDirected))

	// Table header row
	sb.WriteString(fmt.Sprintf("%-8s", "From\\To"))
	for _, j := range keys {
		node, _ := gr.GetNodeByKey(j)
		if node != nil && node.Label != "" {
			sb.WriteString(fmt.Sprintf("%-12s", fmt.Sprintf("%d(%s)", j, node.Label)))
		} else {
			sb.WriteString(fmt.Sprintf("%-12d", j))
		}
	}
	sb.WriteString("\n")

	// Table separator
	sb.WriteString(strings.Repeat("─", 8+len(keys)*12) + "\n")

	// Distance matrix rows
	infinity := int64(1 << 30)
	for _, i := range keys {
		node, _ := gr.GetNodeByKey(i)
		if node != nil && node.Label != "" {
			sb.WriteString(fmt.Sprintf("%-8s", fmt.Sprintf("%d(%s)", i, node.Label)))
		} else {
			sb.WriteString(fmt.Sprintf("%-8d", i))
		}

		// Distance values for this row
		for _, j := range keys {
			dist := apsp.Distances[i][j]
			if dist == infinity {
				sb.WriteString(fmt.Sprintf("%-12s", "inf"))
			} else if i == j {
				sb.WriteString(fmt.Sprintf("%-12s", "0"))
			} else {
				sb.WriteString(fmt.Sprintf("%-12d", dist))
			}
		}
		sb.WriteString("\n")
	}

	// Connectivity analysis
	sb.WriteString("\nCONNECTIVITY:\n")
	sb.WriteString(strings.Repeat("─", 50) + "\n")

	reachablePairs := 0
	totalPairs := len(keys) * (len(keys) - 1)

	for _, i := range keys {
		for _, j := range keys {
			if i != j && apsp.Distances[i][j] < infinity {
				reachablePairs++
			}
		}
	}

	sb.WriteString(fmt.Sprintf("Reachable vertex pairs: %d/%d (%.1f%%)\n",
		reachablePairs, totalPairs, float64(reachablePairs)/float64(totalPairs)*100))

	return sb.String()
}

// getSortedKeys returns sorted node keys for consistent processing
// Important for Floyd-Warshall to maintain consistent vertex ordering
func getSortedKeys(nodes map[graph.TKey]*graph.Node) []graph.TKey {
	keys := make([]graph.TKey, 0, len(nodes))
	for key := range nodes {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}

```