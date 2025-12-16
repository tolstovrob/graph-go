= Потоки V

== Условие
Максимальный поток. Алгоритм Эдмонса Карпа

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
 * Task: Find maximum flow using Edmonds-Karp algorithm (BFS-based Ford-Fulkerson)
 */

// FlowEdge represents an edge with flow information
type FlowEdge struct {
	Source      graph.TKey    `json:"source"`
	Destination graph.TKey    `json:"destination"`
	Capacity    graph.TWeight `json:"capacity"`
	Flow        graph.TWeight `json:"flow"`
}

// MaxFlowResult contains the result of maximum flow calculation
type MaxFlowResult struct {
	MaxFlowValue graph.TWeight `json:"max_flow_value"`
	Source       graph.TKey    `json:"source"`
	Sink         graph.TKey    `json:"sink"`
	FlowEdges    []FlowEdge    `json:"flow_edges"`
	MinCut       []graph.TKey  `json:"min_cut"`
	Message      string        `json:"message"`
}

// FindMaxFlow finds maximum flow from source to sink using Edmonds-Karp algorithm
func FindMaxFlow(gr *graph.Graph, source, sink graph.TKey) (*MaxFlowResult, error) {
	if gr.Nodes == nil {
		return nil, graph.ThrowNodesListIsNil()
	}

	// Validate source and sink nodes
	if _, err := gr.GetNodeByKey(source); err != nil {
		return nil, fmt.Errorf("source node %d does not exist", source)
	}

	if _, err := gr.GetNodeByKey(sink); err != nil {
		return nil, fmt.Errorf("sink node %d does not exist", sink)
	}

	if source == sink {
		return nil, fmt.Errorf("source and sink cannot be the same node")
	}

	// Create residual graph and initialize flow
	residualGraph := createResidualGraph(gr)
	flowMap := initializeFlowMap(gr)

	maxFlow := graph.TWeight(0)

	// Edmonds-Karp algorithm: repeatedly find augmenting paths using BFS
	for {
		// Find augmenting path in residual graph
		path, parent := findAugmentingPath(residualGraph, source, sink)
		if path == nil {
			break // No more augmenting paths
		}

		// Find bottleneck capacity in the path
		pathFlow := findBottleneckCapacity(residualGraph, path, parent, sink)

		// Update residual capacities and flow along the path
		updateResidualGraph(residualGraph, flowMap, path, parent, pathFlow, sink, gr)

		maxFlow += pathFlow
	}

	// Build result
	flowEdges := buildFlowEdges(gr, flowMap)
	minCut := findMinCut(residualGraph, source)

	return &MaxFlowResult{
		MaxFlowValue: maxFlow,
		Source:       source,
		Sink:         sink,
		FlowEdges:    flowEdges,
		MinCut:       minCut,
		Message:      fmt.Sprintf("Maximum flow from %d to %d is %d", source, sink, maxFlow),
	}, nil
}

// createResidualGraph creates the residual graph from original graph
func createResidualGraph(gr *graph.Graph) map[graph.TKey]map[graph.TKey]graph.TWeight {
	residual := make(map[graph.TKey]map[graph.TKey]graph.TWeight)

	// Initialize residual capacities
	for u := range gr.Nodes {
		residual[u] = make(map[graph.TKey]graph.TWeight)
		for v := range gr.Nodes {
			residual[u][v] = 0
		}
	}

	// Set initial capacities from original edges
	for _, edge := range gr.Edges {
		capacity := edge.Weight
		if capacity <= 0 {
			capacity = 1 // Default capacity for zero/negative weights
		}
		residual[edge.Source][edge.Destination] = capacity
	}

	return residual
}

// initializeFlowMap creates initial flow map with zero flow
func initializeFlowMap(gr *graph.Graph) map[graph.TKey]map[graph.TKey]graph.TWeight {
	flowMap := make(map[graph.TKey]map[graph.TKey]graph.TWeight)
	for u := range gr.Nodes {
		flowMap[u] = make(map[graph.TKey]graph.TWeight)
		for v := range gr.Nodes {
			flowMap[u][v] = 0
		}
	}
	return flowMap
}

// findAugmentingPath finds a path from source to sink using BFS
func findAugmentingPath(residualGraph map[graph.TKey]map[graph.TKey]graph.TWeight, source, sink graph.TKey) ([]graph.TKey, map[graph.TKey]graph.TKey) {
	visited := make(map[graph.TKey]bool)
	parent := make(map[graph.TKey]graph.TKey)
	queue := []graph.TKey{source}
	visited[source] = true

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]

		// Check all neighbors with positive residual capacity
		for v, capacity := range residualGraph[u] {
			if !visited[v] && capacity > 0 {
				parent[v] = u
				visited[v] = true
				queue = append(queue, v)

				// If we reached sink, reconstruct path
				if v == sink {
					return reconstructPath(parent, source, sink), parent
				}
			}
		}
	}

	return nil, parent
}

// reconstructPath builds the path from parent pointers
func reconstructPath(parent map[graph.TKey]graph.TKey, source, sink graph.TKey) []graph.TKey {
	path := []graph.TKey{}
	current := sink

	for current != source {
		path = append([]graph.TKey{current}, path...)
		current = parent[current]
	}
	path = append([]graph.TKey{source}, path...)

	return path
}

// findBottleneckCapacity finds the minimum residual capacity along the path
func findBottleneckCapacity(residualGraph map[graph.TKey]map[graph.TKey]graph.TWeight, path []graph.TKey, parent map[graph.TKey]graph.TKey, sink graph.TKey) graph.TWeight {
	bottleneck := graph.TWeight(1 << 30) // Large number
	v := sink

	for v != path[0] { // While not at source
		u := parent[v]
		if residualGraph[u][v] < bottleneck {
			bottleneck = residualGraph[u][v]
		}
		v = u
	}

	return bottleneck
}

// updateResidualGraph updates residual capacities and flow after augmenting path
func updateResidualGraph(residualGraph map[graph.TKey]map[graph.TKey]graph.TWeight, flowMap map[graph.TKey]map[graph.TKey]graph.TWeight, path []graph.TKey, parent map[graph.TKey]graph.TKey, pathFlow graph.TWeight, sink graph.TKey, gr *graph.Graph) {
	v := sink

	for v != path[0] { // While not at source
		u := parent[v]

		// Update residual capacities
		residualGraph[u][v] -= pathFlow
		residualGraph[v][u] += pathFlow

		// Update flow
		if hasOriginalEdge(gr, u, v) {
			flowMap[u][v] += pathFlow
		} else {
			// Backward edge - subtract flow
			flowMap[v][u] -= pathFlow
		}

		v = u
	}
}

// hasOriginalEdge checks if an edge exists in the original graph
func hasOriginalEdge(gr *graph.Graph, u, v graph.TKey) bool {
	for _, edge := range gr.Edges {
		if edge.Source == u && edge.Destination == v {
			return true
		}
	}
	return false
}

// buildFlowEdges creates the list of flow edges from flow map
func buildFlowEdges(gr *graph.Graph, flowMap map[graph.TKey]map[graph.TKey]graph.TWeight) []FlowEdge {
	flowEdges := []FlowEdge{}

	for u := range flowMap {
		for v := range flowMap[u] {
			flow := flowMap[u][v]
			if flow > 0 {
				// Find original capacity
				capacity := graph.TWeight(1)
				for _, edge := range gr.Edges {
					if edge.Source == u && edge.Destination == v {
						if edge.Weight > 0 {
							capacity = edge.Weight
						}
						break
					}
				}

				flowEdges = append(flowEdges, FlowEdge{
					Source:      u,
					Destination: v,
					Capacity:    capacity,
					Flow:        flow,
				})
			}
		}
	}

	// Sort for consistent output
	sort.Slice(flowEdges, func(i, j int) bool {
		if flowEdges[i].Source == flowEdges[j].Source {
			return flowEdges[i].Destination < flowEdges[j].Destination
		}
		return flowEdges[i].Source < flowEdges[j].Source
	})

	return flowEdges
}

// findMinCut finds the minimum cut (reachable nodes from source in residual graph)
func findMinCut(residualGraph map[graph.TKey]map[graph.TKey]graph.TWeight, source graph.TKey) []graph.TKey {
	visited := make(map[graph.TKey]bool)
	queue := []graph.TKey{source}
	visited[source] = true

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]

		for v, capacity := range residualGraph[u] {
			if !visited[v] && capacity > 0 {
				visited[v] = true
				queue = append(queue, v)
			}
		}
	}

	// Convert visited map to sorted slice
	minCut := []graph.TKey{}
	for node := range visited {
		minCut = append(minCut, node)
	}
	sort.Slice(minCut, func(i, j int) bool { return minCut[i] < minCut[j] })

	return minCut
}

// FormatMaxFlowResult creates a formatted string representation
func (result *MaxFlowResult) FormatMaxFlowResult(gr *graph.Graph) string {
	var sb strings.Builder

	sb.WriteString("MAXIMUM FLOW ANALYSIS\n\n")
	sb.WriteString("Algorithm: Edmonds-Karp (BFS-based Ford-Fulkerson)\n")
	sb.WriteString(fmt.Sprintf("Source: %d", result.Source))
	if node, _ := gr.GetNodeByKey(result.Source); node != nil && node.Label != "" {
		sb.WriteString(fmt.Sprintf(" (%s)", node.Label))
	}
	sb.WriteString(fmt.Sprintf("\nSink: %d", result.Sink))
	if node, _ := gr.GetNodeByKey(result.Sink); node != nil && node.Label != "" {
		sb.WriteString(fmt.Sprintf(" (%s)", node.Label))
	}
	sb.WriteString(fmt.Sprintf("\nMaximum Flow Value: %d\n\n", result.MaxFlowValue))

	sb.WriteString("FLOW DISTRIBUTION:\n")
	sb.WriteString(strings.Repeat("─", 60) + "\n")
	sb.WriteString(fmt.Sprintf("%-8s %-8s %-12s %-12s %-12s\n", "From", "To", "Capacity", "Flow", "Utilization"))
	sb.WriteString(fmt.Sprintf("%-8s %-8s %-12s %-12s %-12s\n", "────", "──", "────────", "────", "────────────"))

	totalCapacity := graph.TWeight(0)
	totalFlow := graph.TWeight(0)

	for _, edge := range result.FlowEdges {
		fromLabel := fmt.Sprintf("%d", edge.Source)
		if node, _ := gr.GetNodeByKey(edge.Source); node != nil && node.Label != "" {
			fromLabel = fmt.Sprintf("%d(%s)", edge.Source, node.Label)
		}

		toLabel := fmt.Sprintf("%d", edge.Destination)
		if node, _ := gr.GetNodeByKey(edge.Destination); node != nil && node.Label != "" {
			toLabel = fmt.Sprintf("%d(%s)", edge.Destination, node.Label)
		}

		utilization := "0%"
		if edge.Capacity > 0 {
			utilization = fmt.Sprintf("%.1f%%", float64(edge.Flow)*100/float64(edge.Capacity))
		}

		sb.WriteString(fmt.Sprintf("%-8s %-8s %-12d %-12d %-12s\n",
			fromLabel, toLabel, edge.Capacity, edge.Flow, utilization))

		totalCapacity += edge.Capacity
		totalFlow += edge.Flow
	}

	sb.WriteString("\nSUMMARY:\n")
	sb.WriteString(strings.Repeat("─", 40) + "\n")
	sb.WriteString(fmt.Sprintf("Total capacity: %d\n", totalCapacity))
	sb.WriteString(fmt.Sprintf("Total flow: %d\n", totalFlow))
	if totalCapacity > 0 {
		sb.WriteString(fmt.Sprintf("Flow efficiency: %.1f%%\n", float64(totalFlow)*100/float64(totalCapacity)))
	}

	sb.WriteString(fmt.Sprintf("\nMINIMUM CUT (%d nodes):\n", len(result.MinCut)))
	for i, node := range result.MinCut {
		nodeLabel := fmt.Sprintf("%d", node)
		if nodeObj, _ := gr.GetNodeByKey(node); nodeObj != nil && nodeObj.Label != "" {
			nodeLabel = fmt.Sprintf("%d (%s)", node, nodeObj.Label)
		}
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, nodeLabel))
	}

	return sb.String()
}
```

#image("/assets/image-19.png")
