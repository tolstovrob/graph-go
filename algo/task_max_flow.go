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

type FlowEdge struct {
	Source      graph.TKey    `json:"source"`
	Destination graph.TKey    `json:"destination"`
	Capacity    graph.TWeight `json:"capacity"`
	Flow        graph.TWeight `json:"flow"`
}

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

	if _, err := gr.GetNodeByKey(source); err != nil {
		return nil, fmt.Errorf("source node %d does not exist", source)
	}

	if _, err := gr.GetNodeByKey(sink); err != nil {
		return nil, fmt.Errorf("sink node %d does not exist", sink)
	}

	if source == sink {
		return nil, fmt.Errorf("source and sink cannot be the same node")
	}

	// Create residual graph
	residualGraph := createResidualGraph(gr)

	// Initialize flow
	maxFlow := graph.TWeight(0)
	flowMap := make(map[graph.TKey]map[graph.TKey]graph.TWeight)

	// Initialize flow map
	for u := range gr.Nodes {
		flowMap[u] = make(map[graph.TKey]graph.TWeight)
		for v := range gr.Nodes {
			flowMap[u][v] = 0
		}
	}

	// Edmonds-Karp algorithm
	for {
		// Find augmenting path using BFS
		path, parent := findAugmentingPath(residualGraph, source, sink)
		if path == nil {
			break
		}

		// Find minimum residual capacity along the path
		pathFlow := graph.TWeight(1 << 30) // Large number
		v := sink
		for v != source {
			u := parent[v]
			residualCapacity := residualGraph[u][v]
			if residualCapacity < pathFlow {
				pathFlow = residualCapacity
			}
			v = u
		}

		// Update residual capacities and flow
		v = sink
		for v != source {
			u := parent[v]

			// Update residual graph
			residualGraph[u][v] -= pathFlow
			residualGraph[v][u] += pathFlow

			// Update flow
			if _, exists := gr.Edges[getEdgeKey(gr, u, v)]; exists {
				// Forward edge
				flowMap[u][v] += pathFlow
			} else {
				// Backward edge (subtract flow)
				flowMap[v][u] -= pathFlow
			}

			v = u
		}

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

// createResidualGraph creates the residual graph from the original graph
func createResidualGraph(gr *graph.Graph) map[graph.TKey]map[graph.TKey]graph.TWeight {
	residual := make(map[graph.TKey]map[graph.TKey]graph.TWeight)

	// Initialize residual graph
	for u := range gr.Nodes {
		residual[u] = make(map[graph.TKey]graph.TWeight)
		for v := range gr.Nodes {
			residual[u][v] = 0
		}
	}

	// Fill with capacities from original edges
	for _, edge := range gr.Edges {
		// Use weight as capacity, if weight is 0, assume capacity 1
		capacity := edge.Weight
		if capacity <= 0 {
			capacity = 1
		}
		residual[edge.Source][edge.Destination] = capacity
	}

	return residual
}

// findAugmentingPath finds an augmenting path using BFS
func findAugmentingPath(residualGraph map[graph.TKey]map[graph.TKey]graph.TWeight, source, sink graph.TKey) ([]graph.TKey, map[graph.TKey]graph.TKey) {
	visited := make(map[graph.TKey]bool)
	parent := make(map[graph.TKey]graph.TKey)
	queue := []graph.TKey{source}
	visited[source] = true

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]

		for v, capacity := range residualGraph[u] {
			if !visited[v] && capacity > 0 {
				parent[v] = u
				visited[v] = true
				queue = append(queue, v)

				if v == sink {
					// Reconstruct path
					path := []graph.TKey{}
					curr := sink
					for curr != source {
						path = append([]graph.TKey{curr}, path...)
						curr = parent[curr]
					}
					path = append([]graph.TKey{source}, path...)
					return path, parent
				}
			}
		}
	}

	return nil, parent
}

// getEdgeKey finds the key of an edge between two nodes
func getEdgeKey(gr *graph.Graph, u, v graph.TKey) graph.TKey {
	for key, edge := range gr.Edges {
		if edge.Source == u && edge.Destination == v {
			return key
		}
	}
	return 0
}

// buildFlowEdges builds the list of flow edges from the flow map
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
	sb.WriteString(fmt.Sprintf("Flow efficiency: %.1f%%\n", float64(totalFlow)*100/float64(totalCapacity)))

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
