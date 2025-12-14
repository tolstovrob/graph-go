= Минимальные требования для класса Граф
== Условие
Для решения всех задач курса необходимо создать класс (или иерархию классов - на усмотрение разработчика), содержащий:

1. Структуру для хранения списка смежности графа (не работать с графом через матрицы смежности, если в некоторых алгоритмах удобнее использовать список ребер - реализовать метод, создающий список рёбер на основе списка смежности);

2. Конструкторы (не менее 3-х):
  - конструктор по умолчанию, создающий пустой граф
  - конструктор, заполняющий данные графа из файла
  - конструктор-копию (аккуратно, не все сразу делают именно копию)
  - специфические конструкторы для удобства тестирования

3. Методы:

  - добавляющие вершину,
  - добавляющие ребро (дугу),
  - удаляющие вершину,
  - удаляющие ребро (дугу),
  - выводящие список смежности в файл (в том числе в пригодном для чтения конструктором формате).
  - Не выполняйте некорректные операции, сообщайте об ошибках.

4. Должны поддерживаться как ориентированные, так и неориентированные графы. Заранее предусмотрите возможность добавления меток и\или весов для дуг. Поддержка мультиграфа не требуется.

5. Добавьте минималистичный консольный интерфейс пользователя (не смешивая его с реализацией!), позволяющий добавлять и удалять вершины и рёбра (дуги) и просматривать текущий список смежности графа.

6. Сгенерируйте не менее 4 входных файлов с разными типами графов (балансируйте на комбинации ориентированность-взвешенность) для тестирования класса в этом и последующих заданиях. Графы должны содержать не менее 7-10 вершин, в том числе петли и изолированные вершины.



Замечание:

В зависимости от выбранного способа хранения графа могут появиться дополнительные трудности при удалении-добавлении, например, необходимость переименования вершин, если граф хранится списком $($например, vector C++, List C#$)$. Этого можно избежать, если хранить в списке пару (имя вершины, список смежных вершин), или хранить в другой структуре (например, Dictionary C#$,$ map в С++, при этом список смежности вершины может также храниться в виде словаря с ключами - смежными вершинами и значениями - весами соответствующих ребер). Идеально, если в качестве вершины реализуется обобщенный тип (generic), но достаточно использовать строковый тип или свой класс.

== Код

Описание графа многофайловое. Полный код проекта приложен в приложении. 

Содержимое файла `graph/types.go`

```go
/*
 * This is a graph package, which contains graoh definition and basic operations
 * on it. As you go through the file, you will see some comments, that are
 * explaining this or that choice, etc.
 *
 * Author: github.com/tolstovrob
 */

package graph

type Option[T any] func(*T) // Type representing functional options pattern

type TKey uint64   // Key type. Can be replaced with any UNIQUE type
type TWeight int64 // Weight type. Can be replaced with any COMPARABLE type
```

Содержимое файла `graph/node.go`

```go
/*
 * This is a graph package, which contains graoh definition and basic operations
 * on it. As you go through the file, you will see some comments, that are
 * explaining this or that choice, etc.
 *
 * Author: github.com/tolstovrob
 */

package graph

/*
 * Node struct represents graph node. It has a unique key and optional label.
 * You can properly construct Node via this:
 *
 * node := MakeNode(1)
 *
 * or this:
 *
 * labeledNode := MakeNode(1, WithNodeLabel("Aboba"))
 */

type Node struct {
	Key   TKey   `json:"key"`
	Label string `json:"label"`
}

func MakeNode(key TKey, options ...Option[Node]) *Node {
	node := &Node{}
	node.Key = key
	for _, opt := range options {
		opt(node)
	}
	return node
}

func (node *Node) UpdateNode(options ...Option[Node]) {
	for _, opt := range options {
		opt(node)
	}
}

func WithNodeLabel(label string) Option[Node] {
	return func(node *Node) {
		node.Label = label
	}
}
```

Содержимое файла `graph/edge.go`

```go
/*
 * This is a graph package, which contains graoh definition and basic operations
 * on it. As you go through the file, you will see some comments, that are
 * explaining this or that choice, etc.
 *
 * Author: github.com/tolstovrob
 */

package graph

/*
 * Edge struct represents graph edge -- a connection between 2 nodes.
 * This edge implementation supposed to be direct. If you need to make
 * undirected edge, you should make 2 edges and work with this.
 *
 * Edge represents connection from Edge.Source to Edge.Destination, optionally
 * labelled and weighted.
 *
 * I.e., there are 2 nodes given:
 *
 * src := MakeNode(1)
 * dst := MakeNode(2)
 *
 * You can properly construct this via this:
 *
 * edge := MakeEdge(1, src.Key, dst.Key)
 *
 * or with optional fields:
 *
 * fullyConstructedEdge := MakeEdge(1, src.Key, dst.Key, WithEdgeLabel("Path"), WithEdgeWeight(69))
 */

type Edge struct {
	Key         TKey    `json:"key"`
	Source      TKey    `json:"source"`
	Destination TKey    `json:"destination"`
	Weight      TWeight `json:"weight"`
	Label       string  `json:"label"`
}

func MakeEdge(key, src, dst TKey, options ...Option[Edge]) *Edge {
	edge := &Edge{}
	edge.Key, edge.Source, edge.Destination = key, src, dst
	for _, opt := range options {
		opt(edge)
	}
	return edge
}

func (edge *Edge) UpdateEdge(options ...Option[Edge]) {
	for _, opt := range options {
		opt(edge)
	}
}

func WithEdgeWeight(weight TWeight) Option[Edge] {
	return func(edge *Edge) {
		edge.Weight = weight
	}
}

func WithEdgeLabel(label string) Option[Edge] {
	return func(edge *Edge) {
		edge.Label = label
	}
}
```

Содержимое файла `graph/graph.go`

```go
/*
 * This is a graph package, which contains graoh definition and basic operations
 * on it. As you go through the file, you will see some comments, that are
 * explaining this or that choice, etc.
 *
 * Author: github.com/tolstovrob
 */

package graph

import (
	"encoding/json"
	"fmt"
	"slices"
)

/*
 * Graph struct.
 *
 * First of all, we got TOptions struct, which represents all possible graph
 * configuration. Now, it only has IsMulti and IsDirected for multigraphs and
 * Directed graphs respectively, but it is easily scalable for other options
 * if neccessary.
 *
 * Graph struct sa it is contains Nodes and Edges lists of Node and Edge
 * pointers respectively, and Options configuration of TOptions.
 *
 * Graph represented via adjacency list of edges. But it also possible to have
 * islands with no connections. You cannot find them in Edges, but in the Nodes.
 *
 * You actually can use default constructor with this one. It will build
 * non-multi undirected graph:
 *
 * gr := Graph{}
 *
 * But to make graph properly, use constructor with options:
 *
 * gr := MakeGraph(WithGraphMulti(true), WithGraphDirected(false))
 *
 * I.e., code above will create undirected multigraph.
 */

type TOptions struct {
	IsMulti    bool `json:"isMulti"`
	IsDirected bool `json:"IsDirected"`
}

type Graph struct {
	Nodes        map[TKey]*Node  `json:"nodes"`
	Edges        map[TKey]*Edge  `json:"edges"`
	AdjacencyMap map[TKey][]TKey `json:"adjacencyMap"`
	Options      TOptions        `json:"options"`
}

func MakeGraph(options ...Option[Graph]) *Graph {
	gr := &Graph{}
	gr.Nodes = make(map[TKey]*Node)
	gr.Edges = make(map[TKey]*Edge)
	gr.AdjacencyMap = make(map[TKey][]TKey)
	for _, opt := range options {
		opt(gr)
	}
	return gr
}

func (gr *Graph) Copy() *Graph {
	newGraph := MakeGraph(
		WithGraphDirected(gr.Options.IsDirected),
		WithGraphMulti(gr.Options.IsMulti),
	)

	for key, node := range gr.Nodes {
		newGraph.Nodes[key] = &Node{
			Key:   node.Key,
			Label: node.Label,
		}
	}

	for key, edge := range gr.Edges {
		newGraph.Edges[key] = &Edge{
			Key:         edge.Key,
			Source:      edge.Source,
			Destination: edge.Destination,
			Weight:      edge.Weight,
			Label:       edge.Label,
		}
	}

	newGraph.RebuildAdjacencyMap()
	return newGraph
}

func (gr *Graph) RebuildEdges() {
	newEdges := make(map[TKey]*Edge)
	edgeKeysUsed := make(map[TKey]bool)
	edgeKeyCounter := TKey(1)

	nextEdgeKey := func() TKey {
		for {
			if !edgeKeysUsed[edgeKeyCounter] {
				edgeKeysUsed[edgeKeyCounter] = true
				key := edgeKeyCounter
				edgeKeyCounter++
				return key
			}
			edgeKeyCounter++
		}
	}

	edgeID := func(src, dst TKey) string {
		if !gr.Options.IsDirected && src > dst {
			src, dst = dst, src
		}
		return fmt.Sprintf("%d-%d", src, dst)
	}

	seenEdges := make(map[string]bool)

	for _, edge := range gr.Edges {
		id := edgeID(edge.Source, edge.Destination)
		if !gr.Options.IsMulti {
			if seenEdges[id] {
				continue
			}
			seenEdges[id] = true
		}

		key := edge.Key
		if key == 0 || edgeKeysUsed[key] {
			key = nextEdgeKey()
		} else {
			edgeKeysUsed[key] = true
		}

		newEdge := &Edge{
			Key:         key,
			Source:      edge.Source,
			Destination: edge.Destination,
			Weight:      edge.Weight,
			Label:       edge.Label,
		}
		newEdges[key] = newEdge
	}

	gr.Edges = newEdges
}

func (gr *Graph) RebuildAdjacencyMap() {
	gr.AdjacencyMap = make(map[TKey][]TKey)
	for _, edge := range gr.Edges {
		gr.AdjacencyMap[edge.Source] = append(gr.AdjacencyMap[edge.Source], edge.Destination)
		if !gr.Options.IsDirected {
			gr.AdjacencyMap[edge.Destination] = append(gr.AdjacencyMap[edge.Destination], edge.Source)
		}
	}
}

/*
 * Later in the code, Graph.RebuildAdjacencyMap will be called many times.
 * It could really affect performance on huge amount of edges, but since
 * it is just academical example, we will pretend it never happens.
 *
 * Anyways, need to fix, so mark this part as WIP!
 */

func (gr *Graph) UpdateGraph(options ...Option[Graph]) {
	oldOptions := gr.Options

	for _, opt := range options {
		opt(gr)
	}

	if oldOptions != gr.Options {
		gr.RebuildEdges()
		gr.RebuildAdjacencyMap()
	}
}

func WithGraphNodes(nodes map[TKey]*Node) Option[Graph] {
	return func(gr *Graph) {
		gr.Nodes = nodes
	}
}

func WithGraphEdges(edges map[TKey]*Edge) Option[Graph] {
	return func(gr *Graph) {
		gr.Edges = edges
	}
}
func WithGraphAdjacencyMap(adj map[TKey][]TKey) Option[Graph] {
	return func(gr *Graph) {
		gr.AdjacencyMap = adj
	}
}

func WithGraphOptions(options TOptions) Option[Graph] {
	return func(gr *Graph) {
		gr.Options = options
	}
}

func WithGraphMulti(isMulti bool) Option[Graph] {
	return func(gr *Graph) {
		gr.Options.IsMulti = isMulti
	}
}

func WithGraphDirected(IsDirected bool) Option[Graph] {
	return func(gr *Graph) {
		gr.Options.IsDirected = IsDirected
	}
}

/*
 * Next coming finding, adding and removing handlers for nodes and edges. I put
 * them apart main Graph struct because they contain both node and edges and
 * graph. All of them will throw an error if the operation is not allowed (I.e.
 * adding existing node or connecting nodes with more then one time in multi).
 */

func (gr *Graph) GetNodeByKey(key TKey) (*Node, error) {
	if gr.Nodes == nil {
		return nil, ThrowNodesListIsNil()
	}

	if _, exists := gr.Nodes[key]; !exists {
		return nil, ThrowNodeWithKeyNotExists(key)
	}

	return gr.Nodes[key], nil
}

func (gr *Graph) AddNode(node *Node) error {
	if node, _ := gr.GetNodeByKey(node.Key); node != nil {
		return ThrowNodeWithKeyExists(node.Key)
	}

	gr.Nodes[node.Key] = node
	return nil
}

func (gr *Graph) RemoveNodeByKey(key TKey) error {
	if _, err := gr.GetNodeByKey(key); err != nil {
		return err
	}

	// Remove all edges connected to this node
	edgesToRemove := make([]TKey, 0)
	for edgeKey, edge := range gr.Edges {
		if edge.Source == key || edge.Destination == key {
			edgesToRemove = append(edgesToRemove, edgeKey)
		}
	}

	for _, edgeKey := range edgesToRemove {
		delete(gr.Edges, edgeKey)
	}

	// Remove the node
	delete(gr.Nodes, key)
	delete(gr.AdjacencyMap, key)

	// Remove references to this node from other nodes' adjacency lists
	for nodeKey, neighbors := range gr.AdjacencyMap {
		filteredNeighbors := make([]TKey, 0)
		for _, neighbor := range neighbors {
			if neighbor != key {
				filteredNeighbors = append(filteredNeighbors, neighbor)
			}
		}
		gr.AdjacencyMap[nodeKey] = filteredNeighbors
	}

	return nil
}

func (gr *Graph) GetEdgeByKey(key TKey) (*Edge, error) {
	if gr.Edges == nil {
		return nil, ThrowEdgesListIsNil()
	}

	if _, exists := gr.Edges[key]; !exists {
		return nil, ThrowEdgeWithKeyNotExists(key)
	}

	return gr.Edges[key], nil
}

func (gr *Graph) AddEdge(edge *Edge) error {
	if edge, _ := gr.GetEdgeByKey(edge.Key); edge != nil {
		return ThrowEdgeWithKeyExists(edge.Key)
	}

	if !gr.Options.IsMulti &&
		(slices.Contains(gr.AdjacencyMap[edge.Source], edge.Destination) ||
			!gr.Options.IsDirected && slices.Contains(gr.AdjacencyMap[edge.Destination], edge.Source)) {
		return ThrowSameEdgeNotAllowed(edge.Source, edge.Destination)
	}

	if src, _ := gr.GetNodeByKey(edge.Source); src == nil {
		return ThrowEdgeEndNotExists(edge.Key, edge.Source)
	}

	if dst, _ := gr.GetNodeByKey(edge.Destination); dst == nil {
		return ThrowEdgeEndNotExists(edge.Key, edge.Destination)
	}

	gr.Edges[edge.Key] = edge
	gr.RebuildAdjacencyMap()
	return nil
}

func (gr *Graph) RemoveEdgeByKey(key TKey) error {
	if edge, _ := gr.GetEdgeByKey(key); edge == nil {
		return ThrowEdgeWithKeyNotExists(key)
	}

	delete(gr.Edges, key)
	gr.RebuildAdjacencyMap()
	return nil
}

/*
 * File handling moved to CLI service -- here will be declared just marshalling
 * and unmarshalling handlers
 */

func (gr *Graph) MarshalJSON() ([]byte, error) {
	type MarshalGraph Graph
	return json.Marshal(&struct {
		*MarshalGraph
	}{
		MarshalGraph: (*MarshalGraph)(gr),
	})
}

func (gr *Graph) UnmarshalJSON(data []byte) error {
	type MarshalGraph Graph
	aux := &struct {
		*MarshalGraph
	}{
		MarshalGraph: (*MarshalGraph)(gr),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return ThrowGraphUnmarshalError()
	}
	gr.RebuildAdjacencyMap()
	return nil
}

func (gr *Graph) ToJSON() (string, error) {
	b, err := json.Marshal(gr)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (gr *Graph) FromJSON(jsonData string) error {
	return json.Unmarshal([]byte(jsonData), gr)
}

/*
 * Following methods are for checking some graph properties.
 * They are can be useful for some tasks
 */

func (gr *Graph) IsTree() bool {
	if len(gr.Nodes) == 0 {
		return true
	}

	// Check edge count: tree must have n-1 edges
	if len(gr.Edges) != len(gr.Nodes)-1 {
		return false
	}

	// Check connectivity and acyclicity
	return gr.IsConnected() && !gr.HasCycle()
}

func (gr *Graph) IsConnected() bool {
	if len(gr.Nodes) == 0 {
		return true
	}

	visited := make(map[TKey]bool)

	// Start from first node
	var startKey TKey
	for key := range gr.Nodes {
		startKey = key
		break
	}

	// BFS traversal
	queue := []TKey{startKey}
	visited[startKey] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, neighbor := range gr.AdjacencyMap[current] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}

	return len(visited) == len(gr.Nodes)
}

func (gr *Graph) HasCycle() bool {
	if len(gr.Nodes) == 0 {
		return false
	}

	visited := make(map[TKey]bool)

	for node := range gr.Nodes {
		if !visited[node] {
			if gr.hasCycleDFS(node, 0, visited) {
				return true
			}
		}
	}

	return false
}

func (gr *Graph) hasCycleDFS(current, parent TKey, visited map[TKey]bool) bool {
	visited[current] = true

	for _, neighbor := range gr.AdjacencyMap[current] {
		if !visited[neighbor] {
			if gr.hasCycleDFS(neighbor, current, visited) {
				return true
			}
		} else if neighbor != parent {
			return true
		}
	}

	return false
}

// GetConnectedComponents returns the number of connected components in the graph
func (gr *Graph) GetConnectedComponents() int {
	if len(gr.Nodes) == 0 {
		return 0
	}

	visited := make(map[TKey]bool)
	componentCount := 0

	for node := range gr.Nodes {
		if !visited[node] {
			componentCount++
			gr.bfsComponent(node, visited)
		}
	}

	return componentCount
}

// bfsComponent performs BFS to mark all nodes in the same connected component
func (gr *Graph) bfsComponent(start TKey, visited map[TKey]bool) {
	queue := []TKey{start}
	visited[start] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, neighbor := range gr.AdjacencyMap[current] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}
}

// GetComponentSizes returns the sizes of all connected components
func (gr *Graph) GetComponentSizes() []int {
	if len(gr.Nodes) == 0 {
		return []int{}
	}

	visited := make(map[TKey]bool)
	var sizes []int

	for node := range gr.Nodes {
		if !visited[node] {
			size := gr.bfsComponentWithSize(node, visited)
			sizes = append(sizes, size)
		}
	}

	return sizes
}

// bfsComponentWithSize performs BFS and returns the size of the component
func (gr *Graph) bfsComponentWithSize(start TKey, visited map[TKey]bool) int {
	queue := []TKey{start}
	visited[start] = true
	size := 1

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, neighbor := range gr.AdjacencyMap[current] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
				size++
			}
		}
	}

	return size
}
```

== Пример интерфейса в консоли

#image("/assets/image.png")

Реализация TUI CLI находится в пакете cli. 
