package list

import (
	"fmt"
	set "github.com/deckarep/golang-set"
)

// Grapher defines the commonly used interfaces of a graph structure
type Grapher interface {
	// AddVertex inserts a node to the graph
	AddVertex(vertex Vertex) Vertex
	// DeleteVertex remove a node from the graph by vertex id. The edges that are related to the node will be removed in the meanwhile
	DeleteVertex(id string) error
	// AddEdge connects two nodes (the source and the target node) in the graph
	AddEdge(edge Edge)
	// DeleteEdge will delete the edge which connects a source vertex to the target. The related inbound and outbound edges will be deleted in the meanwhile
	DeleteEdge(edge Edge) error
	// ContainsVertex check if a vertex exists in the graph by id
	ContainsVertex(id string) bool
	// ContainsEdge checks if there is a directed connection between two nodes in the graph
	ContainsEdge(sourceId string, targetId string) bool
	// Vertices returns the list of all the vertices in the graph.
	Vertices() []Vertex
	// Edges returns the list of all the edges in the graph.
	Edges() []Edge
	// Size returns the number of the vertices in the graph
	Size() int
	// Draw generates a visual representation of the graph
	Draw() string
}

// Graph is a structure that defines a directed graph that contains vertices and directed edges
type Graph struct {
	// vertices defines all the vertices in the graph
	vertices set.Set
	// edges defines all the edges in the graph
	edges set.Set
	// inboundEdges defines a mapping of the vertices and their inbound edges in the graph.
	// The key is the vertex id and the value is a set of all the inbound edges to the vertex.
	inboundEdges map[string]set.Set
	// outboundEdges defines a mapping of the vertices and their outbound edges in the graph.
	// The key is the vertex id and the value is a set of all the outbound edges from the vertex.
	outboundEdges map[string]set.Set
}

// DirectedAcyclicGraph defines a directed graph with no directed cycles
type DirectedAcyclicGraph struct {
	Graph
}

// NewDirectedAcyclicGraph creates an empty DAG without any vertices and edges
func NewDirectedAcyclicGraph() *DirectedAcyclicGraph {
	return &DirectedAcyclicGraph{Graph: *NewGraph()}
}

type Identifier interface {
	Id() string
}

// Vertex defines the node with in a graph with a unique id for index
type Vertex interface {
	Identifier
}

// Edge defines the directed connection between two nodes in a graph with a unique id for index
type Edge interface {
	Identifier
	Source() Vertex
	Target() Vertex
}

// NewGraph creates an empty graph without any vertices and edges
func NewGraph() *Graph {
	return &Graph{
		vertices:      set.NewSet(),
		edges:         set.NewSet(),
		inboundEdges:  make(map[string]set.Set),
		outboundEdges: make(map[string]set.Set),
	}
}

// Size returns the number of the vertices in the graph
func (g *Graph) Size() int {
	return g.vertices.Cardinality()
}

// Draw generates a visual representation of the graph
func (g *Graph) Draw() string {
	panic("implement me")
}

// AddVertex inserts a node to the graph
func (g *Graph) AddVertex(vertex Vertex) Vertex {
	if vertex == nil {
		return nil
	}
	g.vertices.Add(vertex)
	return vertex
}

// DeleteVertex remove a node from the graph by vertex id. The edges that are related to the node will be removed in the meanwhile
func (g *Graph) DeleteVertex(id string) error {

	if !g.vertices.Contains(id) {
		return VertexUnknownError{id}
	}

	// delete the vertex from its upstream vertices' outbound edges
	if _, ok := g.inboundEdges[id]; ok {
		for upstream := range g.inboundEdges[id].Iter() {
			g.outboundEdges[upstream.(Vertex).Id()].Remove(id)
		}
	}

	// delete the vertex from its downstream vertices' inbound edges
	if _, ok := g.outboundEdges[id]; ok {
		for downstream := range g.outboundEdges[id].Iter() {
			g.inboundEdges[downstream.(Vertex).Id()].Remove(id)
		}
	}

	// delete the vertex's inbound edges
	delete(g.inboundEdges, id)
	// delete the vertex's outbound edges
	delete(g.outboundEdges, id)

	// delete the vertex from the vertices map
	g.vertices.Remove(id)
	return nil
}

// ContainsVertex check if a vertex is contained in the graph by id
func (g *Graph) ContainsVertex(id string) bool {
	return g.vertices.Contains(id)
}

// ContainsEdge checks if the edge is contained in the graph
// The graph contains the edge if there's already an edge which has the same source and target vertex with the given edge
func (g *Graph) ContainsEdge(source, target string) bool {
	if _, ok := g.inboundEdges[target]; !ok {
		return false
	}
	if !g.inboundEdges[target].Contains(source) {
		return false
	}
	if _, ok := g.outboundEdges[source]; !ok {
		return false
	}
	if !g.outboundEdges[source].Contains(target) {
		return false
	}
	return true
}

// AddEdge adds an edge with the given source and target and records it in the inbound edges and the outbound edges
func (g *Graph) AddEdge(edge Edge) {
	if g.ContainsEdge(edge.Source().Id(), edge.Target().Id()) {
		// the graph already contains an edge which connects the edge.Source to the edge.Target
		return
	}
	g.edges.Add(edge)

	if _, ok := g.outboundEdges[edge.Source().Id()]; !ok {
		g.outboundEdges[edge.Source().Id()] = set.NewSet()
	}
	g.outboundEdges[edge.Source().Id()].Add(edge.Target())

	if _, ok := g.inboundEdges[edge.Target().Id()]; !ok {
		g.inboundEdges[edge.Target().Id()] = set.NewSet()
	}
	g.inboundEdges[edge.Target().Id()].Add(edge.Source())
}

// DeleteEdge will delete the edge which connects a source vertex to the target. The related inbound and outbound edges will be deleted in the meanwhile
func (g *Graph) DeleteEdge(edge Edge) error {
	if !g.ContainsEdge(edge.Source().Id(), edge.Target().Id()) {
		return EdgeUnknownError{source: edge.Source().Id(), target: edge.Target().Id()}
	}
	g.edges.Remove(edge.Id())
	g.inboundEdges[edge.Target().Id()].Remove(edge.Source().Id())
	g.outboundEdges[edge.Source().Id()].Remove(edge.Target().Id())
	return nil
}

// visit traverses a DAG from a start node in the specified direction and returns a map of the visited nodes' id
func visit(g *DirectedAcyclicGraph, id string, visited map[string]bool, down bool) {
	var nexts set.Set
	if down {
		nexts = g.outboundEdges[id]
	} else {
		nexts = g.inboundEdges[id]
	}

	if nexts != nil {
		for next := range nexts.Iter() {
			if _, ok := visited[next.(Vertex).Id()]; ok {
				continue
			}
			visit(g, next.(Vertex).Id(), visited, down)
		}
	}

	visited[id] = true
}

// Vertices returns the list of all the vertices in the graph
func (g *Graph) Vertices() []Vertex {
	list := make([]Vertex, 0, g.vertices.Cardinality())
	for elem := range g.vertices.Iter() {
		list = append(list, elem.(Vertex))
	}
	return list
}

// Edges returns the list of all the edges in the graph
func (g *Graph) Edges() []Edge {
	list := make([]Edge, 0, g.edges.Cardinality())
	for elem := range g.edges.Iter() {
		list = append(list, elem.(Edge))
	}
	return list
}

// EdgeUnknownError is the error type when the given edge does not exist in the graph
type EdgeUnknownError struct {
	source string
	target string
}

func (e EdgeUnknownError) Error() string {
	return fmt.Sprintf("edge between '%s' and '%s' is unknown", e.source, e.target)
}

// VertexUnknownError is the error type when the given vertex does not exist in the graph
type VertexUnknownError struct {
	vertex string
}

func (e VertexUnknownError) Error() string {
	return fmt.Sprintf("vertex(name:%s) is unknown", e.vertex)
}
