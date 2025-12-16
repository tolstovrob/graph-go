= Взвешенный граф IVa
== Условие
Вывести длины кратчайших путей для всех пар вершин. Алгоритм Флойда.

== Код

```go
/*
 * This package contains algorithms and tasks for my SSU course
 */

package algo

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/tolstovrob/graph-go/graph"
)

/*
 * Task: Find eccentricity of vertices and graph radius using Dijkstra's algorithm
 *
 * Eccentricity of a vertex = maximum shortest path distance from that vertex to all other vertices
 * Graph Radius = minimum eccentricity among all vertices
 * Graph Diameter = maximum eccentricity among all vertices
 */

// EccentricityResult represents the result of eccentricity and radius calculation
type EccentricityResult struct {
	Eccentricities     map[graph.TKey]int64 // Eccentricity value for each vertex
	Radius             int64                // Graph radius (minimum eccentricity)
	Diameter           int64                // Graph diameter (maximum eccentricity)
	CenterVertices     []graph.TKey         // Vertices with eccentricity = radius
	PeripheralVertices []graph.TKey         // Vertices with eccentricity = diameter
	IsConnected        bool                 // Whether graph is connected
	Message            string               // Status message
}

// FindEccentricityAndRadius calculates eccentricity for all vertices and graph radius
// Time Complexity: O(V * E log V) with Dijkstra for each vertex
func FindEccentricityAndRadius(gr *graph.Graph) (*EccentricityResult, error) {
	if gr.Nodes == nil {
		return nil, graph.ThrowNodesListIsNil()
	}

	// Handle empty graph case
	if len(gr.Nodes) == 0 {
		return &EccentricityResult{
			Eccentricities:     make(map[graph.TKey]int64),
			Radius:             0,
			Diameter:           0,
			CenterVertices:     []graph.TKey{},
			PeripheralVertices: []graph.TKey{},
			IsConnected:        true,
			Message:            "Graph is empty",
		}, nil
	}

	// Step 1: Check for negative weights - Dijkstra cannot handle them
	for _, edge := range gr.Edges {
		if edge.Weight < 0 {
			return nil, fmt.Errorf("Dijkstra's algorithm cannot handle negative weights. Edge %d has weight %d", edge.Key, edge.Weight)
		}
	}

	// For each vertex, find maximum distance to all other vertices
	eccentricities := make(map[graph.TKey]int64)
	for vertex := range gr.Nodes {
		// Run Dijkstra from current vertex to find shortest paths to all others
		distances, err := dijkstra(gr, vertex)
		if err != nil {
			return nil, err
		}

		// Eccentricity = maximum distance to any reachable vertex
		eccentricity := int64(0)
		for _, dist := range distances {
			if dist > eccentricity && dist != math.MaxInt64 {
				eccentricity = dist
			}
		}

		// If any vertex is unreachable, graph is disconnected
		// In disconnected graphs, eccentricity is considered infinite
		for _, dist := range distances {
			if dist == math.MaxInt64 {
				eccentricity = math.MaxInt64
				break
			}
		}

		eccentricities[vertex] = eccentricity
	}

	// Step 3: Calculate radius and diameter
	radius := int64(math.MaxInt64) // Start with "infinity"
	diameter := int64(0)           // Start with 0
	var centerVertices, peripheralVertices []graph.TKey

	// Find minimum and maximum eccentricity values
	for _, ecc := range eccentricities {
		if ecc != math.MaxInt64 { // Only consider reachable vertices
			if ecc < radius {
				radius = ecc
			}
			if ecc > diameter {
				diameter = ecc
			}
		}
	}

	// Step 4: Find center and peripheral vertices
	for vertex, ecc := range eccentricities {
		if ecc == radius {
			centerVertices = append(centerVertices, vertex)
		}
		if ecc == diameter {
			peripheralVertices = append(peripheralVertices, vertex)
		}
	}

	// Step 5: Handle disconnected graphs
	if radius == math.MaxInt64 {
		radius = -1 // Special value for disconnected graphs
		diameter = -1
	}

	return &EccentricityResult{
		Eccentricities:     eccentricities,
		Radius:             radius,
		Diameter:           diameter,
		CenterVertices:     centerVertices,
		PeripheralVertices: peripheralVertices,
		IsConnected:        radius != -1,
		Message:            fmt.Sprintf("Found eccentricities for %d vertices", len(gr.Nodes)),
	}, nil
}

// dijkstra implements Dijkstra's algorithm for single-source shortest paths
// Returns distances from source vertex to all other vertices
func dijkstra(gr *graph.Graph, source graph.TKey) (map[graph.TKey]int64, error) {
	distances := make(map[graph.TKey]int64)
	visited := make(map[graph.TKey]bool)

	// Initialize all distances to "infinity"
	for vertex := range gr.Nodes {
		distances[vertex] = math.MaxInt64
	}
	distances[source] = 0 // Distance to self is 0

	// Main Dijkstra loop - process all vertices
	for len(visited) < len(gr.Nodes) {
		// Find unvisited vertex with minimum distance
		minVertex := graph.TKey(0)
		minDist := math.MaxInt64

		for vertex, dist := range distances {
			if !visited[vertex] && dist < int64(minDist) {
				minDist = int(dist)
				minVertex = vertex
			}
		}

		// If no more reachable vertices, stop
		if minDist == math.MaxInt64 {
			break
		}

		// Mark vertex as visited
		visited[minVertex] = true

		// Update distances to all neighbors
		for _, neighbor := range gr.AdjacencyMap[minVertex] {
			if !visited[neighbor] {
				edgeWeight := getEdgeWeight(gr, minVertex, neighbor)
				if edgeWeight == math.MaxInt64 {
					continue // No edge exists
				}

				// Relaxation step: update distance if shorter path found
				newDist := distances[minVertex] + edgeWeight
				if newDist < distances[neighbor] {
					distances[neighbor] = newDist
				}
			}
		}
	}

	return distances, nil
}

// getEdgeWeight finds the weight of an edge between two vertices
func getEdgeWeight(gr *graph.Graph, u, v graph.TKey) int64 {
	for _, edge := range gr.Edges {
		if (edge.Source == u && edge.Destination == v) ||
			(!gr.Options.IsDirected && edge.Source == v && edge.Destination == u) {
			return int64(edge.Weight)
		}
	}
	return math.MaxInt64 // No edge found
}

// FormatEccentricityResult creates a formatted string representation
func (result *EccentricityResult) FormatEccentricityResult(gr *graph.Graph) string {
	var sb strings.Builder

	sb.WriteString("ECCENTRICITY AND RADIUS ANALYSIS\n\n")
	sb.WriteString("Algorithm: Dijkstra's Algorithm\n")
	sb.WriteString(fmt.Sprintf("Total vertices: %d\n", len(gr.Nodes)))
	sb.WriteString(fmt.Sprintf("Graph connected: %v\n", result.IsConnected))
	sb.WriteString(fmt.Sprintf("Radius: %s\n", formatDistance(result.Radius)))
	sb.WriteString(fmt.Sprintf("Diameter: %s\n", formatDistance(result.Diameter)))
	sb.WriteString(fmt.Sprintf("Center vertices: %d\n", len(result.CenterVertices)))
	sb.WriteString(fmt.Sprintf("Peripheral vertices: %d\n\n", len(result.PeripheralVertices)))

	sb.WriteString("ECCENTRICITIES BY VERTEX:\n")
	sb.WriteString(strings.Repeat("─", 50) + "\n")

	// Sort vertices for consistent output
	vertices := make([]graph.TKey, 0, len(result.Eccentricities))
	for v := range result.Eccentricities {
		vertices = append(vertices, v)
	}
	sort.Slice(vertices, func(i, j int) bool { return vertices[i] < vertices[j] })

	// Display eccentricity for each vertex
	for _, vertex := range vertices {
		ecc := result.Eccentricities[vertex]
		node, _ := gr.GetNodeByKey(vertex)

		if node != nil && node.Label != "" {
			sb.WriteString(fmt.Sprintf("Vertex %d (%s): %s\n", vertex, node.Label, formatDistance(ecc)))
		} else {
			sb.WriteString(fmt.Sprintf("Vertex %d: %s\n", vertex, formatDistance(ecc)))
		}
	}

	// Display center vertices (vertices with minimum eccentricity)
	if len(result.CenterVertices) > 0 {
		sb.WriteString("\nCENTER VERTICES (eccentricity = radius):\n")
		for i, vertex := range result.CenterVertices {
			node, _ := gr.GetNodeByKey(vertex)
			if node != nil && node.Label != "" {
				sb.WriteString(fmt.Sprintf("%d. Vertex %d (%s)\n", i+1, vertex, node.Label))
			} else {
				sb.WriteString(fmt.Sprintf("%d. Vertex %d\n", i+1, vertex))
			}
		}
	}

	// Display peripheral vertices (vertices with maximum eccentricity)
	if len(result.PeripheralVertices) > 0 && result.Diameter != -1 {
		sb.WriteString("\nPERIPHERAL VERTICES (eccentricity = diameter):\n")
		for i, vertex := range result.PeripheralVertices {
			node, _ := gr.GetNodeByKey(vertex)
			if node != nil && node.Label != "" {
				sb.WriteString(fmt.Sprintf("%d. Vertex %d (%s)\n", i+1, vertex, node.Label))
			} else {
				sb.WriteString(fmt.Sprintf("%d. Vertex %d\n", i+1, vertex))
			}
		}
	}

	return sb.String()
}

// formatDistance converts numerical distance to readable string
func formatDistance(d int64) string {
	if d == math.MaxInt64 {
		return "inf (unreachable)"
	}
	if d == -1 {
		return "inf (graph disconnected)"
	}
	return fmt.Sprintf("%d", d)
}
```

#image("/assets/image-15.png")
