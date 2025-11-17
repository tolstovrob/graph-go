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
 */

type AllPairsShortestPath struct {
	Distances map[graph.TKey]map[graph.TKey]int64
	Next      map[graph.TKey]map[graph.TKey]graph.TKey
	IsValid   bool
	Message   string
}

func FindAllPairsShortestPath(gr *graph.Graph) (*AllPairsShortestPath, error) {
	if gr.Nodes == nil {
		return nil, graph.ThrowNodesListIsNil()
	}

	return findFloydWarshall(gr)
}

func findFloydWarshall(gr *graph.Graph) (*AllPairsShortestPath, error) {
	if len(gr.Nodes) == 0 {
		return &AllPairsShortestPath{
			Distances: make(map[graph.TKey]map[graph.TKey]int64),
			Next:      make(map[graph.TKey]map[graph.TKey]graph.TKey),
			IsValid:   true,
			Message:   "Graph is empty",
		}, nil
	}

	// Get sorted node keys for consistent processing
	keys := getSortedKeys(gr.Nodes)

	// Initialize distance and next matrices
	dist := make(map[graph.TKey]map[graph.TKey]int64)
	next := make(map[graph.TKey]map[graph.TKey]graph.TKey)

	infinity := int64(1 << 30)

	// Initialize matrices
	for _, i := range keys {
		dist[i] = make(map[graph.TKey]int64)
		next[i] = make(map[graph.TKey]graph.TKey)

		for _, j := range keys {
			if i == j {
				dist[i][j] = 0 // Distance to self is 0
			} else {
				dist[i][j] = infinity // Unknown distance
			}
			next[i][j] = 0 // No next node initially
		}
	}

	// Initialize with direct edges
	for _, edge := range gr.Edges {
		weight := int64(edge.Weight)
		if weight < 0 {
			return &AllPairsShortestPath{
				IsValid: false,
				Message: "Graph contains negative weights",
			}, nil
		}

		// Set direct edge distance
		if weight < dist[edge.Source][edge.Destination] {
			dist[edge.Source][edge.Destination] = weight
			next[edge.Source][edge.Destination] = edge.Destination
		}

		// For undirected graphs, set reverse edge
		if !gr.Options.IsDirected {
			if weight < dist[edge.Destination][edge.Source] {
				dist[edge.Destination][edge.Source] = weight
				next[edge.Destination][edge.Source] = edge.Source
			}
		}
	}

	// Floyd-Warshall algorithm
	for _, k := range keys {
		for _, i := range keys {
			// Skip if no path from i to k
			if dist[i][k] == infinity {
				continue
			}

			for _, j := range keys {
				// Skip if no path from k to j
				if dist[k][j] == infinity {
					continue
				}

				// Check if path through k is better
				if dist[i][k]+dist[k][j] < dist[i][j] {
					dist[i][j] = dist[i][k] + dist[k][j]
					next[i][j] = next[i][k]
				}
			}
		}
	}

	// Check for negative cycles (unnecessary, but anyway it is better to check than not to)
	for _, k := range keys {
		if dist[k][k] < 0 {
			return &AllPairsShortestPath{
				IsValid: false,
				Message: "Graph contains negative weight cycles",
			}, nil
		}
	}

	return &AllPairsShortestPath{
		Distances: dist,
		Next:      next,
		IsValid:   true,
		Message:   fmt.Sprintf("Computed shortest paths for %d vertices", len(keys)),
	}, nil
}

func getSortedKeys(nodes map[graph.TKey]*graph.Node) []graph.TKey {
	keys := make([]graph.TKey, 0, len(nodes))
	for key := range nodes {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}

func (apsp *AllPairsShortestPath) GetPath(start, end graph.TKey) []graph.TKey {
	if apsp.Next[start][end] == 0 {
		return nil
	}

	path := []graph.TKey{start}
	current := start

	for current != end {
		current = apsp.Next[current][end]
		path = append(path, current)
	}

	return path
}

func (apsp *AllPairsShortestPath) FormatDistanceMatrix(gr *graph.Graph) string {
	if !apsp.IsValid {
		return apsp.Message
	}

	var sb strings.Builder
	keys := getSortedKeys(gr.Nodes)

	sb.WriteString("SHORTEST PATH DISTANCES BETWEEN ALL PAIRS OF VERTICES\n\n")
	sb.WriteString("Algorithm: Floyd-Warshall\n")
	sb.WriteString(fmt.Sprintf("Total vertices: %d\n", len(keys)))
	sb.WriteString(fmt.Sprintf("Total edges: %d\n", len(gr.Edges)))
	sb.WriteString(fmt.Sprintf("Directed: %v\n\n", gr.Options.IsDirected))

	// Header row
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

	sb.WriteString(strings.Repeat("─", 8+len(keys)*12) + "\n")

	// Distance matrix
	infinity := int64(1 << 30)
	for _, i := range keys {
		node, _ := gr.GetNodeByKey(i)
		if node != nil && node.Label != "" {
			sb.WriteString(fmt.Sprintf("%-8s", fmt.Sprintf("%d(%s)", i, node.Label)))
		} else {
			sb.WriteString(fmt.Sprintf("%-8d", i))
		}

		for _, j := range keys {
			dist := apsp.Distances[i][j]
			if dist == infinity {
				sb.WriteString(fmt.Sprintf("%-12s", "∞"))
			} else if i == j {
				sb.WriteString(fmt.Sprintf("%-12s", "0"))
			} else {
				sb.WriteString(fmt.Sprintf("%-12d", dist))
			}
		}
		sb.WriteString("\n")
	}

	// Add connectivity information
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
