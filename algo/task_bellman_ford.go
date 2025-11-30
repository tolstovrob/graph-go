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
 * Task: Find all negative cycles using Bellman-Ford algorithm
 *
 * Negative Cycle: A cycle in a graph where the sum of edge weights is negative
 * Bellman-Ford Algorithm: Single-source shortest path algorithm that can detect negative cycles
 * Time Complexity: O(V * E) for each source vertex, O(V^2 * E) total
 */

// NegativeCycle represents a single negative cycle found in the graph
type NegativeCycle struct {
	Vertices    []graph.TKey  `json:"vertices"`     // Ordered list of vertices in the cycle
	Edges       []graph.TKey  `json:"edges"`        // Ordered list of edges in the cycle
	TotalWeight graph.TWeight `json:"total_weight"` // Sum of all edge weights in the cycle
}

// NegativeCyclesResult contains the complete result of negative cycle detection
type NegativeCyclesResult struct {
	Cycles            []NegativeCycle `json:"cycles"`              // List of all unique negative cycles found
	HasNegativeCycles bool            `json:"has_negative_cycles"` // Whether any negative cycles exist
	TotalCycles       int             `json:"total_cycles"`        // Count of unique negative cycles
	Message           string          `json:"message"`             // Status message describing the result
}

// FindNegativeCycles finds all negative cycles in the graph using Bellman-Ford algorithm
// This is the main entry point for negative cycle detection
func FindNegativeCycles(gr *graph.Graph) (*NegativeCyclesResult, error) {
	if gr.Nodes == nil {
		return nil, graph.ThrowNodesListIsNil()
	}

	// Handle empty graph case - no cycles possible
	if len(gr.Nodes) == 0 {
		return &NegativeCyclesResult{
			Cycles:            []NegativeCycle{},
			HasNegativeCycles: false,
			TotalCycles:       0,
			Message:           "Graph is empty",
		}, nil
	}

	// Bellman-Ford requires directed graphs for negative cycle detection
	if !gr.Options.IsDirected {
		return &NegativeCyclesResult{
			Cycles:            []NegativeCycle{},
			HasNegativeCycles: false,
			TotalCycles:       0,
			Message:           "Bellman-Ford algorithm for negative cycles requires directed graph",
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

// findAllNegativeCycles executes Bellman-Ford from each vertex to find all negative cycles
// This is the core algorithm implementation
func findAllNegativeCycles(gr *graph.Graph) []NegativeCycle {
	keys := getSortedNodeKeys(gr.Nodes)    // Get sorted vertices for consistent processing
	allCycles := []NegativeCycle{}         // Store all found cycles
	visitedCycles := make(map[string]bool) // Track seen cycles to avoid duplicates

	// Try each vertex as a potential starting point for cycle detection
	for _, start := range keys {
		dist := make(map[graph.TKey]graph.TWeight)  // Shortest distance estimates
		prev := make(map[graph.TKey]graph.TKey)     // Predecessor vertices for path reconstruction
		edgePrev := make(map[graph.TKey]graph.TKey) // Predecessor edges for cycle tracing

		// Initialize with large values representing infinity
		infinity := graph.TWeight(1 << 30)
		for _, key := range keys {
			dist[key] = infinity
		}
		dist[start] = 0 // Distance to start vertex is 0

		// Relaxation phase: |V| - 1 iterations of edge relaxation
		for i := 0; i < len(keys)-1; i++ {
			changed := false // Track if any distances were updated
			for _, edge := range gr.Edges {
				u, v, w := edge.Source, edge.Destination, edge.Weight
				// If we found a shorter path through u to v, update
				if dist[u] != infinity && dist[u]+w < dist[v] {
					dist[v] = dist[u] + w
					prev[v] = u
					edgePrev[v] = edge.Key
					changed = true
				}
			}
			// Early termination if no improvements in this iteration
			if !changed {
				break
			}
		}

		// Negative cycle detection phase: check if we can still relax edges
		for _, edge := range gr.Edges {
			u, v, w := edge.Source, edge.Destination, edge.Weight
			// If we can still improve after |V|-1 iterations, negative cycle exists
			if dist[u] != infinity && dist[u]+w < dist[v] {
				// Trace and reconstruct the actual cycle
				cycle := traceCycle(gr, u, v, prev, edgePrev)
				if cycle != nil {
					// Normalize cycle representation and check for duplicates
					normalized := normalizeCycle(*cycle)
					cycleKey := generateCycleKey(normalized)
					if !visitedCycles[cycleKey] && normalized.TotalWeight < 0 {
						visitedCycles[cycleKey] = true
						allCycles = append(allCycles, normalized)
					}
				}
			}
		}
	}

	return allCycles
}

// traceCycle traces back from a negatively-weighted edge to find the actual cycle
// Uses Floyd's cycle-finding algorithm (tortoise and hare)
func traceCycle(gr *graph.Graph, u, v graph.TKey, prev, edgePrev map[graph.TKey]graph.TKey) *NegativeCycle {
	// Use two pointers to detect cycle: slow moves 1 step, fast moves 2 steps
	slow, fast := v, v

	// Find meeting point inside the cycle
	for i := 0; i < len(gr.Nodes); i++ {
		if prev[fast] == 0 || prev[prev[fast]] == 0 {
			return nil // Invalid pointers, no cycle found
		}
		slow = prev[slow]
		fast = prev[prev[fast]]
		if slow == fast {
			break // Cycle detected
		}
	}

	if slow != fast {
		return nil // No cycle found
	}

	// Find the start node of the cycle
	cycleStart := findCycleStart(slow, prev)
	if cycleStart == 0 {
		return nil
	}

	// Reconstruct the complete cycle
	return reconstructCycle(gr, cycleStart, prev, edgePrev)
}

// findCycleStart finds the starting node of a cycle using cycle length detection
func findCycleStart(meetingPoint graph.TKey, prev map[graph.TKey]graph.TKey) graph.TKey {
	// Determine cycle length by traversing until we return to meeting point
	cycleLength := 0
	current := meetingPoint
	for {
		current = prev[current]
		cycleLength++
		if current == meetingPoint {
			break
		}
	}

	// Use two pointers to find cycle start: one starts cycleLength steps ahead
	ptr1 := meetingPoint
	for i := 0; i < cycleLength; i++ {
		ptr1 = prev[ptr1]
	}

	// Move both pointers until they meet at cycle start
	ptr2 := meetingPoint
	for ptr1 != ptr2 {
		ptr1 = prev[ptr1]
		ptr2 = prev[ptr2]
	}

	return ptr1
}

// reconstructCycle builds the complete cycle from predecessor information
func reconstructCycle(gr *graph.Graph, start graph.TKey, prev, edgePrev map[graph.TKey]graph.TKey) *NegativeCycle {
	cycleVertices := []graph.TKey{start}
	cycleEdges := []graph.TKey{}
	totalWeight := graph.TWeight(0)
	visited := make(map[graph.TKey]bool)
	visited[start] = true

	// Reconstruct vertices in the cycle by following predecessor links
	current := prev[start]
	for current != start {
		if current == 0 || visited[current] {
			return nil // Invalid cycle
		}
		visited[current] = true
		cycleVertices = append([]graph.TKey{current}, cycleVertices...)
		current = prev[current]

		// Safety check to prevent infinite loops
		if len(cycleVertices) > len(gr.Nodes) {
			return nil
		}
	}

	// Reconstruct edges and calculate total weight
	for i := 0; i < len(cycleVertices); i++ {
		from := cycleVertices[i]
		to := cycleVertices[(i+1)%len(cycleVertices)]

		edge := findEdgeBetween(gr, from, to)
		if edge == nil {
			return nil
		}

		cycleEdges = append(cycleEdges, edge.Key)
		totalWeight += edge.Weight
	}

	if len(cycleVertices) < 2 {
		return nil
	}

	return &NegativeCycle{
		Vertices:    cycleVertices,
		Edges:       cycleEdges,
		TotalWeight: totalWeight,
	}
}

// normalizeCycle rotates the cycle to start with the smallest vertex for consistent comparison
func normalizeCycle(cycle NegativeCycle) NegativeCycle {
	if len(cycle.Vertices) == 0 {
		return cycle
	}

	// Find the vertex with smallest key value
	minIndex := 0
	minVertex := cycle.Vertices[0]
	for i, vertex := range cycle.Vertices {
		if vertex < minVertex {
			minVertex = vertex
			minIndex = i
		}
	}

	// Rotate cycle to start with smallest vertex
	normalizedVertices := make([]graph.TKey, len(cycle.Vertices))
	copy(normalizedVertices, cycle.Vertices[minIndex:])
	copy(normalizedVertices[len(cycle.Vertices)-minIndex:], cycle.Vertices[:minIndex])

	return NegativeCycle{
		Vertices:    normalizedVertices,
		Edges:       cycle.Edges, // Edge order doesn't affect cycle identity
		TotalWeight: cycle.TotalWeight,
	}
}

// generateCycleKey creates a unique string identifier for a cycle
func generateCycleKey(cycle NegativeCycle) string {
	vertices := make([]string, len(cycle.Vertices))
	for i, v := range cycle.Vertices {
		vertices[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(vertices, "-")
}

// Helper functions for graph operations

// getSortedNodeKeys returns vertices sorted by key for consistent processing
func getSortedNodeKeys(nodes map[graph.TKey]*graph.Node) []graph.TKey {
	keys := make([]graph.TKey, 0, len(nodes))
	for key := range nodes {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}

// findEdgeBetween finds the edge between two vertices in the graph
func findEdgeBetween(gr *graph.Graph, from, to graph.TKey) *graph.Edge {
	for _, edge := range gr.Edges {
		if edge.Source == from && edge.Destination == to {
			return edge
		}
	}
	return nil
}

// FormatNegativeCyclesResult creates a human-readable formatted output
func (result *NegativeCyclesResult) FormatNegativeCyclesResult(gr *graph.Graph) string {
	var sb strings.Builder

	sb.WriteString("NEGATIVE CYCLES ANALYSIS\n\n")
	sb.WriteString("Algorithm: Bellman-Ford\n")
	sb.WriteString(fmt.Sprintf("Total vertices: %d\n", len(gr.Nodes)))
	sb.WriteString(fmt.Sprintf("Total edges: %d\n", len(gr.Edges)))
	sb.WriteString(fmt.Sprintf("Graph directed: %v\n", gr.Options.IsDirected))
	sb.WriteString(fmt.Sprintf("Found negative cycles: %v\n", result.HasNegativeCycles))
	sb.WriteString(fmt.Sprintf("Total unique cycles: %d\n\n", result.TotalCycles))

	if !result.HasNegativeCycles {
		sb.WriteString("No negative cycles found in the graph.\n")
		if !gr.Options.IsDirected {
			sb.WriteString("Note: Bellman-Ford for negative cycles requires directed graphs.\n")
		}
		return sb.String()
	}

	// Display each found cycle with detailed information
	for i, cycle := range result.Cycles {
		sb.WriteString(fmt.Sprintf("NEGATIVE CYCLE %d:\n", i+1))
		sb.WriteString(strings.Repeat("─", 40) + "\n")
		sb.WriteString(fmt.Sprintf("Total weight: %d\n", cycle.TotalWeight))
		sb.WriteString(fmt.Sprintf("Length: %d vertices, %d edges\n\n", len(cycle.Vertices), len(cycle.Edges)))

		sb.WriteString("CYCLE PATH:\n")
		// Display each edge in the cycle
		for j := 0; j < len(cycle.Vertices); j++ {
			current := cycle.Vertices[j]
			next := cycle.Vertices[(j+1)%len(cycle.Vertices)]

			currentNode, _ := gr.GetNodeByKey(current)
			nextNode, _ := gr.GetNodeByKey(next)

			var currentLabel, nextLabel string
			if currentNode != nil && currentNode.Label != "" {
				currentLabel = fmt.Sprintf(" (%s)", currentNode.Label)
			}
			if nextNode != nil && nextNode.Label != "" {
				nextLabel = fmt.Sprintf(" (%s)", nextNode.Label)
			}

			edge := findEdgeBetween(gr, current, next)
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
