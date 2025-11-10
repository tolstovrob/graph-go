/*
 * This a CLI service for my graph implementation. It is build with tview and
 * represents TUI CLI.
 *
 * Author: github.com/tolstovrob
 */

package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tolstovrob/graph-go/graph"
)

func (cli *CLIService) showGraphOptions() {
	form := tview.NewForm()

	form.AddCheckbox("Directed Graph", cli.graph.Options.IsDirected, func(checked bool) {
		cli.graph.UpdateGraph(graph.WithGraphDirected(checked))
	})
	form.AddCheckbox("Multi Graph", cli.graph.Options.IsMulti, func(checked bool) {
		cli.graph.UpdateGraph(graph.WithGraphMulti(checked))
	})

	form.AddButton("Save", func() {
		cli.updateStatus("Graph options updated", Success)
		cli.pages.SwitchToPage("main")
	})
	form.AddButton("Back", func() {
		cli.pages.SwitchToPage("main")
	})

	form.SetBorder(true).SetTitle(" Graph Options ")
	cli.pages.AddAndSwitchToPage("graph_options", form, true)
}

func (cli *CLIService) showGraphInfo() {
	info := cli.getDetailedGraphInfo()
	cli.showScrollableModal("Graph Information", info, "main")
}

func (cli *CLIService) showJSONOperations() {
	modal := tview.NewModal().
		SetText("JSON Operations").
		AddButtons([]string{"Save to JSON", "Load from JSON", "Show JSON", "Back"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonLabel {
			case "Save to JSON":
				cli.showSaveJSONForm()
			case "Load from JSON":
				cli.showLoadJSONForm()
			case "Show JSON":
				cli.showJSONView()
			case "Back":
				cli.pages.SwitchToPage("main")
			}
		})

	cli.pages.AddAndSwitchToPage("json_operations", modal, true)
}

func (cli *CLIService) getDetailedGraphInfo() string {
	var info strings.Builder

	info.WriteString("GRAPH OPTIONS\n")
	info.WriteString(strings.Repeat("─", 50) + "\n")
	info.WriteString(fmt.Sprintf("Directed: %v\n", cli.graph.Options.IsDirected))
	info.WriteString(fmt.Sprintf("Multi-graph: %v\n", cli.graph.Options.IsMulti))
	info.WriteString(fmt.Sprintf("Graph Type: %s%s\n\n",
		map[bool]string{true: "Directed", false: "Undirected"}[cli.graph.Options.IsDirected],
		map[bool]string{true: " Multi", false: ""}[cli.graph.Options.IsMulti]))

	info.WriteString("GRAPH ANALYSIS\n")
	info.WriteString(strings.Repeat("─", 50) + "\n")
	info.WriteString(fmt.Sprintf("Connected: %v\n", cli.graph.IsConnected()))
	info.WriteString(fmt.Sprintf("Connected components: %d\n", cli.graph.GetConnectedComponents()))
	info.WriteString(fmt.Sprintf("Has cycles: %v\n", cli.graph.HasCycle()))
	info.WriteString(fmt.Sprintf("Is tree: %v\n", cli.graph.IsTree()))
	info.WriteString("\n")

	info.WriteString("STATISTICS\n")
	info.WriteString(strings.Repeat("─", 50) + "\n")
	info.WriteString(fmt.Sprintf("Total Nodes: %d\n", len(cli.graph.Nodes)))
	info.WriteString(fmt.Sprintf("Total Edges: %d\n\n", len(cli.graph.Edges)))

	info.WriteString("NODES LIST\n")
	info.WriteString(strings.Repeat("─", 50) + "\n")
	if len(cli.graph.Nodes) == 0 {
		info.WriteString("No nodes in graph\n")
	} else {
		keys := make([]graph.TKey, 0, len(cli.graph.Nodes))
		for key := range cli.graph.Nodes {
			keys = append(keys, key)
		}
		sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

		for _, key := range keys {
			node := cli.graph.Nodes[key]
			degree := len(cli.graph.AdjacencyMap[key])
			info.WriteString(fmt.Sprintf("Key: %4d | Label: %-20s | Degree: %d\n",
				key, node.Label, degree))
		}
	}
	info.WriteString("\n")

	info.WriteString("EDGES LIST\n")
	info.WriteString(strings.Repeat("─", 50) + "\n")
	if len(cli.graph.Edges) == 0 {
		info.WriteString("No edges in graph\n")
	} else {
		keys := make([]graph.TKey, 0, len(cli.graph.Edges))
		for key := range cli.graph.Edges {
			keys = append(keys, key)
		}
		sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

		for _, key := range keys {
			edge := cli.graph.Edges[key]
			info.WriteString(fmt.Sprintf("Key: %4d | %4d → %4d | Weight: %4d | Label: %s\n",
				key, edge.Source, edge.Destination, edge.Weight, edge.Label))
		}
	}
	info.WriteString("\n")

	info.WriteString("ADJACENCY LIST\n")
	info.WriteString(strings.Repeat("─", 50) + "\n")
	if len(cli.graph.AdjacencyMap) == 0 {
		info.WriteString("Empty adjacency list\n")
	} else {
		keys := make([]graph.TKey, 0, len(cli.graph.AdjacencyMap))
		for key := range cli.graph.AdjacencyMap {
			keys = append(keys, key)
		}
		sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

		for _, key := range keys {
			neighbors := cli.graph.AdjacencyMap[key]
			info.WriteString(fmt.Sprintf("%4d → [", key))

			sort.Slice(neighbors, func(i, j int) bool { return neighbors[i] < neighbors[j] })

			for i, neighbor := range neighbors {
				if i > 0 {
					info.WriteString(", ")
				}
				info.WriteString(fmt.Sprintf("%d", neighbor))
			}
			info.WriteString(fmt.Sprintf("] (degree: %d)\n", len(neighbors)))
		}
	}

	return info.String()
}

func (cli *CLIService) showSaveJSONForm() {
	form := tview.NewForm()
	var filename string

	form.AddInputField("Filename", "graph.json", 30, nil, func(text string) {
		filename = text
	})
	form.AddButton("Save", func() {
		if filename == "" {
			cli.updateStatus("Error: Filename cannot be empty", Error)
			return
		}

		jsonData, err := cli.graph.ToJSON()
		if err != nil {
			cli.updateStatus(fmt.Sprintf("Error generating JSON: %v", err), Error)
			return
		}

		err = os.WriteFile(filename, []byte(jsonData), 0644)
		if err != nil {
			cli.updateStatus(fmt.Sprintf("Error writing file: %v", err), Error)
			return
		}

		cli.updateStatus(fmt.Sprintf("Graph saved to %s successfully", filename), Success)
		cli.pages.SwitchToPage("main")
	})
	form.AddButton("Cancel", func() {
		cli.pages.SwitchToPage("json_operations")
	})

	form.SetBorder(true).SetTitle(" Save Graph to JSON ")
	cli.pages.AddAndSwitchToPage("save_json", form, true)
}

func (cli *CLIService) showLoadJSONForm() {
	form := tview.NewForm()
	var filename string

	form.AddInputField("Filename", "graph.json", 30, nil, func(text string) {
		filename = text
	})
	form.AddButton("Load", func() {
		if filename == "" {
			cli.updateStatus("Error: Filename cannot be empty", Error)
			return
		}

		data, err := os.ReadFile(filename)
		if err != nil {
			cli.updateStatus(fmt.Sprintf("Error reading file: %v", err), Error)
			return
		}

		newGraph := graph.MakeGraph()
		if err := newGraph.FromJSON(string(data)); err != nil {
			cli.updateStatus(fmt.Sprintf("Error parsing JSON: %v", err), Error)
			return
		}

		cli.graph = newGraph
		cli.updateStatus(fmt.Sprintf("Graph loaded from %s successfully", filename), Success)
		cli.pages.SwitchToPage("main")
	})
	form.AddButton("Cancel", func() {
		cli.pages.SwitchToPage("json_operations")
	})

	form.SetBorder(true).SetTitle(" Load Graph from JSON ")
	cli.pages.AddAndSwitchToPage("load_json", form, true)
}

func (cli *CLIService) showJSONView() {
	jsonData, err := cli.graph.ToJSON()
	if err != nil {
		cli.updateStatus(fmt.Sprintf("Error generating JSON: %v", err), Error)
		return
	}

	var formattedJSON map[string]any
	if err := json.Unmarshal([]byte(jsonData), &formattedJSON); err != nil {
		cli.updateStatus(fmt.Sprintf("Error formatting JSON: %v", err), Error)
		return
	}

	prettyJSON, err := json.MarshalIndent(formattedJSON, "", "  ")
	if err != nil {
		cli.updateStatus(fmt.Sprintf("Error formatting JSON: %v", err), Error)
		return
	}

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			cli.app.Draw()
		}).
		SetText(string(prettyJSON))

	textView.SetBorder(true).SetTitle(" Graph JSON - Use arrow keys to scroll, Q to go back ")

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(textView, 0, 1, true)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Rune() == 'Q' {
			cli.pages.SwitchToPage("json_operations")
			return nil
		}
		return event
	})

	cli.pages.AddAndSwitchToPage("json_view", flex, true)
}
