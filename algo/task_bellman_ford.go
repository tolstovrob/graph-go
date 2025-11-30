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
 * Task: Find all negative cycles using exhaustive search
 */

type NegativeCycle struct {
	Vertices    []graph.TKey  `json:"vertices"`
	Edges       []graph.TKey  `json:"edges"`
	TotalWeight graph.TWeight `json:"total_weight"`
}

type NegativeCyclesResult struct {
	Cycles            []NegativeCycle `json:"cycles"`
	HasNegativeCycles bool            `json:"has_negative_cycles"`
	TotalCycles       int             `json:"total_cycles"`
	Message           string          `json:"message"`
}

// FindNegativeCycles finds all negative cycles in the graph
func FindNegativeCycles(gr *graph.Graph) (*NegativeCyclesResult, error) {
	if gr.Nodes == nil {
		return nil, graph.ThrowNodesListIsNil()
	}

	if len(gr.Nodes) == 0 {
		return &NegativeCyclesResult{
			Cycles:            []NegativeCycle{},
			HasNegativeCycles: false,
			TotalCycles:       0,
			Message:           "Graph is empty",
		}, nil
	}

	if !gr.Options.IsDirected {
		return &NegativeCyclesResult{
			Cycles:            []NegativeCycle{},
			HasNegativeCycles: false,
			TotalCycles:       0,
			Message:           "Negative cycle detection requires directed graph",
		}, nil
	}

	allCycles := findAllNegativeCycles(gr)

	return &NegativeCyclesResult{
		Cycles:            allCycles,
		HasNegativeCycles: len(allCycles) > 0,
		TotalCycles:       len(allCycles),
		Message:           fmt.Sprintf("Found %d negative cycle(s)", len(allCycles)),
	}, nil
}

// findAllNegativeCycles finds all negative cycles by checking all possible cycles
func findAllNegativeCycles(gr *graph.Graph) []NegativeCycle {
	keys := getSortedNodeKeys(gr.Nodes)
	allCycles := []NegativeCycle{}

	// Check all possible cycles of length 2, 3, 4
	for _, u := range keys {
		for _, v := range keys {
			if u == v {
				continue
			}

			// Check 2-cycles (u->v->u)
			edge1 := findEdge(gr, u, v)
			edge2 := findEdge(gr, v, u)
			if edge1 != nil && edge2 != nil {
				totalWeight := edge1.Weight + edge2.Weight
				if totalWeight < 0 {
					cycle := NegativeCycle{
						Vertices:    []graph.TKey{u, v, u},
						Edges:       []graph.TKey{edge1.Key, edge2.Key},
						TotalWeight: totalWeight,
					}
					if !containsCycle(allCycles, cycle) {
						allCycles = append(allCycles, cycle)
					}
				}
			}

			// Check 3-cycles (u->v->w->u)
			for _, w := range keys {
				if u == w || v == w {
					continue
				}

				edge1 := findEdge(gr, u, v)
				edge2 := findEdge(gr, v, w)
				edge3 := findEdge(gr, w, u)

				if edge1 != nil && edge2 != nil && edge3 != nil {
					totalWeight := edge1.Weight + edge2.Weight + edge3.Weight
					if totalWeight < 0 {
						cycle := NegativeCycle{
							Vertices:    []graph.TKey{u, v, w, u},
							Edges:       []graph.TKey{edge1.Key, edge2.Key, edge3.Key},
							TotalWeight: totalWeight,
						}
						if !containsCycle(allCycles, cycle) {
							allCycles = append(allCycles, cycle)
						}
					}
				}
			}

			// Check 4-cycles (u->v->w->x->u)
			for _, w := range keys {
				for _, x := range keys {
					if u == w || u == x || v == w || v == x || w == x {
						continue
					}

					edge1 := findEdge(gr, u, v)
					edge2 := findEdge(gr, v, w)
					edge3 := findEdge(gr, w, x)
					edge4 := findEdge(gr, x, u)

					if edge1 != nil && edge2 != nil && edge3 != nil && edge4 != nil {
						totalWeight := edge1.Weight + edge2.Weight + edge3.Weight + edge4.Weight
						if totalWeight < 0 {
							cycle := NegativeCycle{
								Vertices:    []graph.TKey{u, v, w, x, u},
								Edges:       []graph.TKey{edge1.Key, edge2.Key, edge3.Key, edge4.Key},
								TotalWeight: totalWeight,
							}
							if !containsCycle(allCycles, cycle) {
								allCycles = append(allCycles, cycle)
							}
						}
					}
				}
			}
		}
	}

	return allCycles
}

// findEdge finds an edge between two vertices
func findEdge(gr *graph.Graph, from, to graph.TKey) *graph.Edge {
	for _, edge := range gr.Edges {
		if edge.Source == from && edge.Destination == to {
			return edge
		}
	}
	return nil
}

// getSortedNodeKeys returns sorted node keys
func getSortedNodeKeys(nodes map[graph.TKey]*graph.Node) []graph.TKey {
	keys := make([]graph.TKey, 0, len(nodes))
	for key := range nodes {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}

// containsCycle checks if a cycle already exists in the list
func containsCycle(cycles []NegativeCycle, newCycle NegativeCycle) bool {
	for _, cycle := range cycles {
		if areCyclesEqual(cycle, newCycle) {
			return true
		}
	}
	return false
}

// areCyclesEqual checks if two cycles are the same
func areCyclesEqual(c1, c2 NegativeCycle) bool {
	if len(c1.Vertices) != len(c2.Vertices) {
		return false
	}

	// Normalize cycles by rotating to start with smallest vertex
	norm1 := normalizeCycle(c1)
	norm2 := normalizeCycle(c2)

	for i := range norm1.Vertices {
		if norm1.Vertices[i] != norm2.Vertices[i] {
			return false
		}
	}

	return true
}

// normalizeCycle rotates the cycle to start with the smallest vertex
func normalizeCycle(cycle NegativeCycle) NegativeCycle {
	if len(cycle.Vertices) == 0 {
		return cycle
	}

	minIndex := 0
	minVertex := cycle.Vertices[0]
	for i, vertex := range cycle.Vertices {
		if vertex < minVertex {
			minVertex = vertex
			minIndex = i
		}
	}

	normalizedVertices := make([]graph.TKey, len(cycle.Vertices))
	copy(normalizedVertices, cycle.Vertices[minIndex:])
	copy(normalizedVertices[len(cycle.Vertices)-minIndex:], cycle.Vertices[:minIndex])

	return NegativeCycle{
		Vertices:    normalizedVertices,
		Edges:       cycle.Edges, // Edges order doesn't matter for comparison
		TotalWeight: cycle.TotalWeight,
	}
}

// FormatNegativeCyclesResult creates a formatted string representation
func (result *NegativeCyclesResult) FormatNegativeCyclesResult(gr *graph.Graph) string {
	var sb strings.Builder

	sb.WriteString("NEGATIVE CYCLES ANALYSIS\n\n")
	sb.WriteString("Algorithm: Exhaustive Cycle Search\n")
	sb.WriteString(fmt.Sprintf("Total vertices: %d\n", len(gr.Nodes)))
	sb.WriteString(fmt.Sprintf("Total edges: %d\n", len(gr.Edges)))
	sb.WriteString(fmt.Sprintf("Graph directed: %v\n", gr.Options.IsDirected))
	sb.WriteString(fmt.Sprintf("Found negative cycles: %v\n", result.HasNegativeCycles))
	sb.WriteString(fmt.Sprintf("Total unique cycles: %d\n\n", result.TotalCycles))

	if !result.HasNegativeCycles {
		sb.WriteString("No negative cycles found in the graph.\n")
		if !gr.Options.IsDirected {
			sb.WriteString("Note: Negative cycle detection requires directed graphs.\n")
		}
		return sb.String()
	}

	for i, cycle := range result.Cycles {
		sb.WriteString(fmt.Sprintf("NEGATIVE CYCLE %d:\n", i+1))
		sb.WriteString(strings.Repeat("─", 40) + "\n")
		sb.WriteString(fmt.Sprintf("Total weight: %d\n", cycle.TotalWeight))
		sb.WriteString(fmt.Sprintf("Length: %d vertices, %d edges\n\n", len(cycle.Vertices), len(cycle.Edges)))

		sb.WriteString("CYCLE PATH:\n")
		for j := 0; j < len(cycle.Vertices)-1; j++ {
			current := cycle.Vertices[j]
			next := cycle.Vertices[j+1]

			currentNode, _ := gr.GetNodeByKey(current)
			nextNode, _ := gr.GetNodeByKey(next)

			var currentLabel, nextLabel string
			if currentNode != nil && currentNode.Label != "" {
				currentLabel = fmt.Sprintf(" (%s)", currentNode.Label)
			}
			if nextNode != nil && nextNode.Label != "" {
				nextLabel = fmt.Sprintf(" (%s)", nextNode.Label)
			}

			edge := findEdge(gr, current, next)
			var edgeLabel string
			if edge != nil && edge.Label != "" {
				edgeLabel = edge.Label
			}

			sb.WriteString(fmt.Sprintf("  %d%s → %d%s", current, currentLabel, next, nextLabel))
			sb.WriteString(fmt.Sprintf(" [Weight: %d", edge.Weight))
			if edgeLabel != "" {
				sb.WriteString(fmt.Sprintf(", Edge: %s", edgeLabel))
			}
			sb.WriteString("]\n")
		}

		sb.WriteString(fmt.Sprintf("\nCycle completes with total weight: %d\n", cycle.TotalWeight))

		if i < len(result.Cycles)-1 {
			sb.WriteString("\n" + strings.Repeat("═", 50) + "\n\n")
		}
	}

	return sb.String()
}
