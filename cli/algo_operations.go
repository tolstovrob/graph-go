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
