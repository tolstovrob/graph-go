/*
 * This a CLI service for my graph implementation. It is build with tview and
 * represents TUI CLI.
 *
 * Author: github.com/tolstovrob
 */

package cli

import (
	"fmt"
	"strconv"

	"github.com/rivo/tview"
	"github.com/tolstovrob/graph-go/algo"
	"github.com/tolstovrob/graph-go/graph"
)

func (cli *CLIService) showAlgorithmsMenu() {
	list := tview.NewList().
		AddItem("In-Degree less than", "Find nodes with in-degree less than target", '1', cli.showInDegreeLessThanForm).
		AddItem("In-nodes in directed", "Find nodes, that are in-nodes for target in directed graph", '2', cli.showIncomingNeighborsForm).
		AddItem("Remove pendant", "Remove all pendant nodes. Destructive action", '3', cli.showRemovePendantVertices).
		AddItem("Vertex to Tree", "Check if removing a vertex makes graph a tree", '4', cli.showVertexToTreeCheck).
		AddItem("Connected Components", "Count and analyze connected components", '5', cli.showConnectedComponentsAnalysis).
		AddItem("Minimum Spanning Tree", "Find MST using Prim's algorithm", '6', cli.showMSTPrim).
		AddItem("All Pairs Shortest Path", "Find shortest paths between all vertices", '7', cli.showAllPairsShortestPath).
		AddItem("Eccentricity and Radius", "Find eccentricity of vertices and graph radius", '8', cli.showEccentricityAndRadius).
		AddItem("Negative Cycles", "Find all negative cycles using Bellman-Ford", '9', cli.showNegativeCycles).
		AddItem("Maximum Flow", "Find maximum flow from source to sink", '0', cli.showMaxFlowForm).
		AddItem("Back to Main Menu", "Return to main menu", 'q', func() {
			cli.pages.SwitchToPage("main")
		})

	list.SetBorder(true).SetTitle(" Graph Algorithms ")
	cli.pages.AddAndSwitchToPage("algorithms_menu", list, true)
}
func (cli *CLIService) showInDegreeLessThanForm() {
	form := tview.NewForm()
	var targetKey string

	form.AddInputField("Target Node Key", "", 10, nil, func(text string) {
		targetKey = text
	})
	form.AddButton("Run Algorithm", func() {
		keyVal, err := strconv.ParseUint(targetKey, 10, 64)
		if err != nil {
			cli.updateStatus("Error: Invalid key format", Error)
			return
		}

		if _, err := cli.graph.GetNodeByKey(graph.TKey(keyVal)); err != nil {
			cli.updateStatus(fmt.Sprintf("Error: Node %d does not exist", keyVal), Error)
			return
		}

		result := algo.InDegreeLessThan(cli.graph, graph.TKey(keyVal))

		var resultText string
		if len(result) == 0 {
			resultText = fmt.Sprintf("No nodes found with in-degree less than target node %d", keyVal)
		} else {
			resultText = fmt.Sprintf("Nodes with in-degree less than target node %d:\n\n", keyVal)
			for i, nodeKey := range result {
				node, _ := cli.graph.GetNodeByKey(nodeKey)
				if node != nil && node.Label != "" {
					resultText += fmt.Sprintf("%d. Node %d (Label: %s)\n", i+1, nodeKey, node.Label)
				} else {
					resultText += fmt.Sprintf("%d. Node %d\n", i+1, nodeKey)
				}
			}
			resultText += fmt.Sprintf("\nTotal: %d nodes", len(result))
		}

		cli.showScrollableModal("Algorithm Result", resultText, "algorithms_menu")
		cli.updateStatus("Algorithm completed successfully", Success)
	})

	form.AddButton("Cancel", func() {
		cli.pages.SwitchToPage("algorithms_menu")
	})

	form.SetBorder(true).SetTitle(" In-Degree Less Than Algorithm ")
	cli.pages.AddAndSwitchToPage("in_degree_algorithm", form, true)
}

func (cli *CLIService) showIncomingNeighborsForm() {
	form := tview.NewForm()
	var targetKey string

	form.AddInputField("Target Vertex Key", "", 10, nil, func(text string) {
		targetKey = text
	})
	form.AddButton("Find Neighbors", func() {
		keyVal, err := strconv.ParseUint(targetKey, 10, 64)
		if err != nil {
			cli.updateStatus("Error: Invalid key format", Error)
			return
		}

		if _, err := cli.graph.GetNodeByKey(graph.TKey(keyVal)); err != nil {
			cli.updateStatus(fmt.Sprintf("Error: Node %d does not exist", keyVal), Error)
			return
		}

		result, err := algo.InNodesInDirected(cli.graph, graph.TKey(keyVal))

		if err != nil {
			cli.updateStatus(fmt.Sprintf("Error: %v", err), Error)
			return
		}

		var resultText string

		resultText = fmt.Sprintf("Incoming neighbors for vertex %d (directed graph):\n\n", keyVal)

		if len(result) == 0 {
			resultText += "No incoming neighbors found"
		} else {
			for i, neighborKey := range result {
				node, _ := cli.graph.GetNodeByKey(neighborKey)
				if node != nil && node.Label != "" {
					resultText += fmt.Sprintf("%d. Node %d (Label: %s)\n", i+1, neighborKey, node.Label)
				} else {
					resultText += fmt.Sprintf("%d. Node %d\n", i+1, neighborKey)
				}
			}
			resultText += fmt.Sprintf("\nTotal: %d incoming neighbors", len(result))
		}

		cli.showScrollableModal("Incoming Neighbors", resultText, "algorithms_menu")
		cli.updateStatus("Incoming neighbors found successfully", Success)
	})

	form.AddButton("Cancel", func() {
		cli.pages.SwitchToPage("algorithms_menu")
	})

	form.SetBorder(true).SetTitle(" Find Incoming Neighbors ")
	cli.pages.AddAndSwitchToPage("incoming_neighbors", form, true)
}

func (cli *CLIService) showRemovePendantVertices() {
	modal := tview.NewModal().
		SetText("This will create a new graph without pendant vertices (degree 1).\n\nOriginal graph will be replaced. Continue?").
		AddButtons([]string{"Yes, Remove Pendant Vertices", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonLabel {
			case "Yes, Remove Pendant Vertices":
				cli.executeRemovePendantVertices()
			case "Cancel":
				cli.pages.SwitchToPage("algorithms_menu")
			}
		})

	cli.pages.AddAndSwitchToPage("remove_pendant_confirm", modal, true)
}

func (cli *CLIService) executeRemovePendantVertices() {
	newGraph, err := algo.RemovePendantVertices(cli.graph)

	if err != nil {
		cli.updateStatus(fmt.Sprintf("Error: %v", err), Error)
		cli.pages.SwitchToPage("algorithms_menu")
		return
	}

	originalNodes := len(cli.graph.Nodes)
	originalEdges := len(cli.graph.Edges)

	cli.graph = newGraph

	newNodes := len(cli.graph.Nodes)
	newEdges := len(cli.graph.Edges)
	removedNodes := originalNodes - newNodes
	removedEdges := originalEdges - newEdges

	resultText := "Pendant Vertices Removal Results:\n\n"
	resultText += fmt.Sprintf("Original graph: %d nodes, %d edges\n", originalNodes, originalEdges)
	resultText += fmt.Sprintf("New graph:      %d nodes, %d edges\n\n", newNodes, newEdges)
	resultText += fmt.Sprintf("Removed:        %d nodes, %d edges\n\n", removedNodes, removedEdges)

	if removedNodes == 0 {
		resultText += "No pendant vertices found in the graph."
	} else {
		resultText += "Graph successfully updated without pendant vertices."
	}

	cli.showScrollableModal("Pendant Vertices Removal", resultText, "algorithms_menu")
	cli.updateStatus(fmt.Sprintf("Removed %d pendant vertices", removedNodes), Success)
}

func (cli *CLIService) showVertexToTreeCheck() {
	if len(cli.graph.Nodes) > 20 {
		modal := tview.NewModal().
			SetText(fmt.Sprintf("Graph has %d vertices. This operation may take some time. Continue?", len(cli.graph.Nodes))).
			AddButtons([]string{"Continue", "Cancel"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				switch buttonLabel {
				case "Continue":
					cli.executeVertexToTreeCheck()
				case "Cancel":
					cli.pages.SwitchToPage("algorithms_menu")
				}
			})
		cli.pages.AddAndSwitchToPage("vertex_tree_confirm", modal, true)
	} else {
		cli.executeVertexToTreeCheck()
	}
}

func (cli *CLIService) executeVertexToTreeCheck() {
	cli.updateStatus("Checking vertices... This may take a while for large graphs", Default)

	go func() {
		result, candidates, err := algo.CanRemoveVertexToMakeTree(cli.graph)

		cli.app.QueueUpdateDraw(func() {
			var resultText string
			if err != nil {
				resultText = fmt.Sprintf("Error: %v", err)
				cli.updateStatus("Algorithm failed", Error)
			} else if result {
				resultText = fmt.Sprintf("SUCCESS: Graph can become a tree by removing %d vertex(es):\n\n", len(candidates))
				for i, vertex := range candidates {
					node, _ := cli.graph.GetNodeByKey(vertex)
					if node != nil && node.Label != "" {
						resultText += fmt.Sprintf("%d. Node %d (Label: %s)\n", i+1, vertex, node.Label)
					} else {
						resultText += fmt.Sprintf("%d. Node %d\n", i+1, vertex)
					}
				}
				resultText += fmt.Sprintf("\nTotal: %d candidate vertices", len(candidates))
				cli.updateStatus(fmt.Sprintf("Found %d candidate vertices", len(candidates)), Success)
			} else {
				resultText = "No such vertex exists - removing any vertex cannot make this graph a tree"
				cli.updateStatus("No candidate vertices found", Default)
			}

			cli.showScrollableModal("Vertex to Tree Check", resultText, "algorithms_menu")
		})
	}()
}

func (cli *CLIService) showConnectedComponentsAnalysis() {
	analysis, err := algo.AnalyzeConnectedComponents(cli.graph)

	var resultText string
	if err != nil {
		resultText = fmt.Sprintf("Error: %v", err)
		cli.updateStatus("Analysis failed", Error)
	} else {
		resultText = fmt.Sprintf("CONNECTED COMPONENTS ANALYSIS\n\n")
		resultText += fmt.Sprintf("Total components: %d\n", analysis.TotalComponents)
		resultText += fmt.Sprintf("Graph is connected: %v\n", analysis.IsConnected)

		if analysis.TotalComponents > 0 {
			resultText += fmt.Sprintf("\nCOMPONENT SIZES:\n")
			for i, size := range analysis.ComponentSizes {
				resultText += fmt.Sprintf("Component %d: %d vertices\n", i+1, size)
			}

			resultText += fmt.Sprintf("\nSTATISTICS:\n")
			resultText += fmt.Sprintf("Largest component: %d vertices\n", analysis.LargestComponent)
			resultText += fmt.Sprintf("Smallest component: %d vertices\n", analysis.SmallestComponent)

			if analysis.TotalComponents > 1 {
				resultText += fmt.Sprintf("Isolated vertices: %d\n", countIsolatedVertices(analysis.ComponentSizes))
			}
		}

		cli.updateStatus(fmt.Sprintf("Found %d connected components", analysis.TotalComponents), Success)
	}

	cli.showScrollableModal("Connected Components Analysis", resultText, "algorithms_menu")
}

func countIsolatedVertices(sizes []int) int {
	count := 0
	for _, size := range sizes {
		if size == 1 {
			count++
		}
	}
	return count
}

func (cli *CLIService) showMSTPrim() {
	cli.updateStatus("Finding Minimum Spanning Tree using Prim's algorithm...", Default)

	go func() {
		result, err := algo.FindMSTPrim(cli.graph)

		cli.app.QueueUpdateDraw(func() {
			var resultText string
			if err != nil {
				resultText = fmt.Sprintf("Error: %v", err)
				cli.updateStatus("MST calculation failed", Error)
			} else if !result.IsPossible {
				resultText = "MINIMUM SPANNING TREE ANALYSIS\n\n"
				resultText += "MST is NOT possible for this graph\n\n"
				resultText += "Reason: Graph is not connected\n"
				resultText += "Prim's algorithm requires the graph to be connected to find a spanning tree."
				cli.updateStatus("Graph is not connected - MST not possible", Error)
			} else {
				resultText = fmt.Sprintf("MINIMUM SPANNING TREE (Prim's Algorithm)\n\n")
				resultText += fmt.Sprintf("Total weight: %d\n", result.TotalWeight)
				resultText += fmt.Sprintf("Number of edges in MST: %d\n", len(result.Edges))
				resultText += fmt.Sprintf("Theoretical minimum edges: %d\n\n", len(cli.graph.Nodes)-1)

				resultText += fmt.Sprintf("MST EDGES:\n")
				resultText += fmt.Sprintf("%-8s %-8s %-8s %-12s %s\n", "Key", "From", "To", "Weight", "Label")
				resultText += fmt.Sprintf("%-8s %-8s %-8s %-12s %s\n", "────", "────", "──", "──────", "─────")

				for _, edge := range result.Edges {
					srcNode, _ := cli.graph.GetNodeByKey(edge.Source)
					dstNode, _ := cli.graph.GetNodeByKey(edge.Destination)

					srcLabel := fmt.Sprintf("%d", edge.Source)
					if srcNode != nil && srcNode.Label != "" {
						srcLabel = fmt.Sprintf("%d(%s)", edge.Source, srcNode.Label)
					}

					dstLabel := fmt.Sprintf("%d", edge.Destination)
					if dstNode != nil && dstNode.Label != "" {
						dstLabel = fmt.Sprintf("%d(%s)", edge.Destination, dstNode.Label)
					}

					resultText += fmt.Sprintf("%-8d %-8s %-8s %-12d %s\n",
						edge.Key, srcLabel, dstLabel, edge.Weight, edge.Label)
				}

				resultText += fmt.Sprintf("\nGRAPH INFORMATION:\n")
				resultText += fmt.Sprintf("Original graph: %d nodes, %d edges\n", len(cli.graph.Nodes), len(cli.graph.Edges))
				resultText += fmt.Sprintf("MST covers: %d nodes, %d edges\n", len(cli.graph.Nodes), len(result.Edges))

				if len(result.Edges) != len(cli.graph.Nodes)-1 {
					resultText += fmt.Sprintf("\nWarning: MST has %d edges but expected %d for %d nodes\n",
						len(result.Edges), len(cli.graph.Nodes)-1, len(cli.graph.Nodes))
				}

				cli.updateStatus(fmt.Sprintf("MST found with total weight %d", result.TotalWeight), Success)
			}

			cli.showScrollableModal("Minimum Spanning Tree", resultText, "algorithms_menu")
		})
	}()
}

func (cli *CLIService) showAllPairsShortestPath() {
	cli.updateStatus("Computing shortest paths between all pairs of vertices...", Default)

	go func() {
		result, err := algo.FindAllPairsShortestPath(cli.graph)

		cli.app.QueueUpdateDraw(func() {
			var resultText string
			if err != nil {
				resultText = fmt.Sprintf("Error: %v", err)
				cli.updateStatus("Shortest path computation failed", Error)
			} else if !result.IsValid {
				resultText = result.Message
				cli.updateStatus("Invalid graph for shortest paths", Error)
			} else {
				resultText = result.FormatDistanceMatrix(cli.graph)
				cli.updateStatus("All-pairs shortest paths computed successfully", Success)
			}

			cli.showScrollableModal("All Pairs Shortest Path", resultText, "algorithms_menu")
		})
	}()
}

func (cli *CLIService) showEccentricityAndRadius() {
	cli.updateStatus("Calculating eccentricity and radius using Dijkstra's algorithm...", Default)

	go func() {
		result, err := algo.FindEccentricityAndRadius(cli.graph)

		cli.app.QueueUpdateDraw(func() {
			var resultText string
			if err != nil {
				resultText = fmt.Sprintf("Error: %v", err)
				cli.updateStatus("Eccentricity calculation failed", Error)
			} else {
				resultText = result.FormatEccentricityResult(cli.graph)
				cli.updateStatus("Eccentricity and radius calculated successfully", Success)
			}

			cli.showScrollableModal("Eccentricity and Radius", resultText, "algorithms_menu")
		})
	}()
}

func (cli *CLIService) showNegativeCycles() {
	cli.updateStatus("Searching for negative cycles using Bellman-Ford algorithm...", Default)

	go func() {
		result, err := algo.FindNegativeCycles(cli.graph)

		cli.app.QueueUpdateDraw(func() {
			var resultText string
			if err != nil {
				resultText = fmt.Sprintf("Error: %v", err)
				cli.updateStatus("Negative cycles search failed", Error)
			} else {
				resultText = result.FormatNegativeCyclesResult(cli.graph)
				if result.HasNegativeCycles {
					cli.updateStatus(fmt.Sprintf("Found %d negative cycle(s)", result.TotalCycles), Error)
				} else {
					cli.updateStatus("No negative cycles found", Success)
				}
			}

			cli.showScrollableModal("Negative Cycles", resultText, "algorithms_menu")
		})
	}()
}

func (cli *CLIService) showMaxFlowForm() {
	form := tview.NewForm()
	var sourceKey, sinkKey string

	form.AddInputField("Source Node Key", "", 10, nil, func(text string) {
		sourceKey = text
	})
	form.AddInputField("Sink Node Key", "", 10, nil, func(text string) {
		sinkKey = text
	})
	form.AddButton("Find Max Flow", func() {
		sourceVal, err := strconv.ParseUint(sourceKey, 10, 64)
		if err != nil {
			cli.updateStatus("Error: Invalid source key format", Error)
			return
		}

		sinkVal, err := strconv.ParseUint(sinkKey, 10, 64)
		if err != nil {
			cli.updateStatus("Error: Invalid sink key format", Error)
			return
		}

		result, err := algo.FindMaxFlow(cli.graph, graph.TKey(sourceVal), graph.TKey(sinkVal))

		var resultText string
		if err != nil {
			resultText = fmt.Sprintf("Error: %v", err)
			cli.updateStatus("Max flow calculation failed", Error)
		} else {
			resultText = result.FormatMaxFlowResult(cli.graph)
			cli.updateStatus(fmt.Sprintf("Max flow: %d from %d to %d", result.MaxFlowValue, sourceVal, sinkVal), Success)
		}

		cli.showScrollableModal("Maximum Flow", resultText, "algorithms_menu")
	})
	form.AddButton("Cancel", func() {
		cli.pages.SwitchToPage("algorithms_menu")
	})

	form.SetBorder(true).SetTitle(" Find Maximum Flow ")
	cli.pages.AddAndSwitchToPage("max_flow", form, true)
}
