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
 */

type EccentricityResult struct {
	Eccentricities     map[graph.TKey]int64 `json:"eccentricities"`
	Radius             int64                `json:"radius"`
	Diameter           int64                `json:"diameter"`
	CenterVertices     []graph.TKey         `json:"center_vertices"`
	PeripheralVertices []graph.TKey         `json:"peripheral_vertices"`
	IsConnected        bool                 `json:"is_connected"`
	Message            string               `json:"message"`
}

// FindEccentricityAndRadius calculates eccentricity for all vertices and graph radius
func FindEccentricityAndRadius(gr *graph.Graph) (*EccentricityResult, error) {
	if gr.Nodes == nil {
		return nil, graph.ThrowNodesListIsNil()
	}

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

	// Check for negative weights
	for _, edge := range gr.Edges {
		if edge.Weight < 0 {
			return nil, fmt.Errorf("Dijkstra's algorithm cannot handle negative weights. Edge %d has weight %d", edge.Key, edge.Weight)
		}
	}

	eccentricities := make(map[graph.TKey]int64)

	// Calculate eccentricity for each vertex
	for vertex := range gr.Nodes {
		distances, err := dijkstra(gr, vertex)
		if err != nil {
			return nil, err
		}

		// Eccentricity is the maximum distance from this vertex to any other reachable vertex
		eccentricity := int64(0)
		for _, dist := range distances {
			if dist > eccentricity && dist != math.MaxInt64 {
				eccentricity = dist
			}
		}

		// If any vertex is unreachable, eccentricity is infinite (represented as MaxInt64)
		for _, dist := range distances {
			if dist == math.MaxInt64 {
				eccentricity = math.MaxInt64
				break
			}
		}

		eccentricities[vertex] = eccentricity
	}

	// Calculate radius (minimum eccentricity) and diameter (maximum eccentricity)
	radius := int64(math.MaxInt64)
	diameter := int64(0)
	var centerVertices, peripheralVertices []graph.TKey

	for _, ecc := range eccentricities {
		if ecc != math.MaxInt64 {
			if ecc < radius {
				radius = ecc
			}
			if ecc > diameter {
				diameter = ecc
			}
		}
	}

	// Find center vertices (vertices with eccentricity = radius)
	// and peripheral vertices (vertices with eccentricity = diameter)
	for vertex, ecc := range eccentricities {
		if ecc == radius {
			centerVertices = append(centerVertices, vertex)
		}
		if ecc == diameter {
			peripheralVertices = append(peripheralVertices, vertex)
		}
	}

	// If graph is disconnected, radius and diameter are infinite
	if radius == math.MaxInt64 {
		radius = -1 // Represent disconnected case
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
func dijkstra(gr *graph.Graph, source graph.TKey) (map[graph.TKey]int64, error) {
	distances := make(map[graph.TKey]int64)
	visited := make(map[graph.TKey]bool)

	// Initialize distances
	for vertex := range gr.Nodes {
		distances[vertex] = math.MaxInt64
	}
	distances[source] = 0

	for len(visited) < len(gr.Nodes) {
		// Find vertex with minimum distance
		minVertex := graph.TKey(0)
		minDist := math.MaxInt64

		for vertex, dist := range distances {
			if !visited[vertex] && dist < int64(minDist) {
				minDist = int(dist)
				minVertex = vertex
			}
		}

		// If no more reachable vertices, break
		if minDist == math.MaxInt64 {
			break
		}

		visited[minVertex] = true

		// Update distances to neighbors
		for _, neighbor := range gr.AdjacencyMap[minVertex] {
			if !visited[neighbor] {
				edgeWeight := getEdgeWeight(gr, minVertex, neighbor)
				if edgeWeight == math.MaxInt64 {
					continue
				}

				newDist := graph.TWeight(distances[minVertex]) + graph.TWeight(edgeWeight)
				if newDist < graph.TWeight(distances[neighbor]) {
					distances[neighbor] = int64(newDist)
				}
			}
		}
	}

	return distances, nil
}

// DECLARED floyd --- func getEdgeWeight(gr *graph.Graph, u, v graph.TKey) int64 {

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

	for _, vertex := range vertices {
		ecc := result.Eccentricities[vertex]
		node, _ := gr.GetNodeByKey(vertex)

		if node != nil && node.Label != "" {
			sb.WriteString(fmt.Sprintf("Vertex %d (%s): %s\n", vertex, node.Label, formatDistance(ecc)))
		} else {
			sb.WriteString(fmt.Sprintf("Vertex %d: %s\n", vertex, formatDistance(ecc)))
		}
	}

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

func formatDistance(d int64) string {
	if d == math.MaxInt64 {
		return "∞ (unreachable)"
	}
	if d == -1 {
		return "∞ (graph disconnected)"
	}
	return fmt.Sprintf("%d", d)
}
