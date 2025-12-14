= Список смежности III
== Условие
Проверить, можно ли из графа удалить какую-либо вершину так, чтобы получилось дерево.

== Код

```go
/*
 * This package contains algorithms and tasks for my SSU course
 */

package algo

import "github.com/tolstovrob/graph-go/graph"

/*
 * Task: Count connected components and analyze their sizes
 */

type ComponentAnalysis struct {
	TotalComponents   int
	ComponentSizes    []int
	IsConnected       bool
	LargestComponent  int
	SmallestComponent int
}

func AnalyzeConnectedComponents(gr *graph.Graph) (*ComponentAnalysis, error) {
	if gr.Nodes == nil {
		return nil, graph.ThrowNodesListIsNil()
	}

	totalComponents := gr.GetConnectedComponents()
	componentSizes := gr.GetComponentSizes()

	var largest, smallest int
	if len(componentSizes) > 0 {
		largest = componentSizes[0]
		smallest = componentSizes[0]
		for _, size := range componentSizes {
			if size > largest {
				largest = size
			}
			if size < smallest {
				smallest = size
			}
		}
	}

	return &ComponentAnalysis{
		TotalComponents:   totalComponents,
		ComponentSizes:    componentSizes,
		IsConnected:       totalComponents == 1,
		LargestComponent:  largest,
		SmallestComponent: smallest,
	}, nil
}

```

Методы для проверки на дерево описаны в структуре графа.