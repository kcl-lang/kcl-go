package list

import (
	"fmt"
	ds "kusionstack.io/kclvm-go/pkg/data_structure"
)

// Grapher defines the commonly used interfaces of a graph structure
type Grapher interface {
	// AddVertex inserts a node to the graph
	AddVertex(vertex Vertex)
	// AddEdge connects two nodes (the source and the target node) in the graph
	AddEdge(edge Edge)
	// DeleteVertex remove a node from the graph by vertex id. The edges that are related to the node will be removed in the meanwhile
	DeleteVertex(vertex Vertex) error
	// DeleteEdge will delete the edge which connects a source vertex to the target. The related inbound and outbound edges will be deleted in the meanwhile
	DeleteEdge(edge Edge) error
	// ContainsVertex check if a vertex exists in the graph by id
	ContainsVertex(vertex Vertex) bool
	// ContainsEdge checks if there is a directed connection between two nodes in the graph
	ContainsEdge(source, target Vertex) bool
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
	vertices *ds.IdentifierSet
	// edges defines all the edges in the graph
	edges map[Edge]bool
	// inboundEdges defines a mapping of the vertices and their inbound edges in the graph.
	// The key is the vertex id and the value is a map of all the source vertex ids of the inbound edges to the vertex.
	inboundEdges map[string]*ds.IdentifierSet
	// outboundEdges defines a mapping of the vertices and their outbound edges in the graph.
	// The key is the vertex id and the value is a map of all the target vertex ids of the outbound edges from the vertex.
	outboundEdges map[string]*ds.IdentifierSet
}

// DirectedAcyclicGraph defines a directed graph with no directed cycles
type DirectedAcyclicGraph struct {
	Graph
}

// NewDirectedAcyclicGraph creates an empty DAG without any vertices and edges
func NewDirectedAcyclicGraph() *DirectedAcyclicGraph {
	return &DirectedAcyclicGraph{Graph: *NewGraph()}
}

// Vertex defines the node with in a graph with a unique id for index
type Vertex interface {
	ds.Identifier
}

// Edge defines the directed connection between two nodes in a graph with a unique id for index
type Edge interface {
	ds.Identifier
	Source() Vertex
	Target() Vertex
}

// NewGraph creates an empty graph without any vertices and edges
func NewGraph() *Graph {
	return &Graph{
		vertices:      &ds.IdentifierSet{Ids: make(map[string]ds.Identifier)},
		edges:         make(map[Edge]bool),
		inboundEdges:  make(map[string]*ds.IdentifierSet),
		outboundEdges: make(map[string]*ds.IdentifierSet),
	}
}

// Size returns the number of the vertices in the graph
func (g *Graph) Size() int {
	return len(g.vertices.Ids)
}

// Draw generates a visual representation of the graph
func (g *Graph) Draw() string {
	panic("implement me")
}

// AddVertex inserts a node to the graph. Nil vertices will not be added.
func (g *Graph) AddVertex(vertex Vertex) {
	if vertex != nil {
		g.vertices.Add(vertex)
	}
}

// DeleteVertex remove a node from the graph by vertex id. The edges that are related to the node will be removed in the meanwhile
func (g *Graph) DeleteVertex(vertex Vertex) error {
	if vertex == nil {
		return VertexUnknownError{}
	}
	if !g.vertices.Contains(vertex.Id()) {
		return VertexUnknownError{vertex.Id()}
	}
	id := vertex.Id()
	// delete the vertex from its upstream vertices' outbound edges
	if _, ok := g.inboundEdges[id]; ok {
		for upstreamId := range g.inboundEdges[id].Ids {
			g.outboundEdges[upstreamId].Remove(id)
		}
	}

	// delete the vertex from its downstream vertices' inbound edges
	if _, ok := g.outboundEdges[id]; ok {
		for downstreamId := range g.outboundEdges[id].Ids {
			g.inboundEdges[downstreamId].Remove(id)
		}
	}

	// delete the vertex's inbound edges
	delete(g.inboundEdges, id)
	// delete the vertex's outbound edges
	delete(g.outboundEdges, id)

	// delete the vertex from the vertices map
	g.vertices.Remove(id)

	// delete the edge from the edges map
	for edge := range g.edges {
		if edge.Source().Id() == id || edge.Target().Id() == id {
			delete(g.edges, edge)
		}
	}
	return nil
}

// ContainsVertex check if a vertex is contained in the graph by id
func (g *Graph) ContainsVertex(vertex Vertex) bool {
	if vertex == nil {
		return false
	}
	return g.vertices.Contains(vertex.Id())
}

// ContainsEdge checks if the edge is contained in the graph
// The edge is contained if there's already an edge which has the same source and target vertex with the given edge
func (g *Graph) ContainsEdge(source, target Vertex) bool {
	if source == nil || target == nil {
		return false
	}
	if _, ok := g.inboundEdges[target.Id()]; !ok {
		return false
	}
	if !g.inboundEdges[target.Id()].Contains(source.Id()) {
		return false
	}
	if _, ok := g.outboundEdges[source.Id()]; !ok {
		return false
	}
	if !g.outboundEdges[source.Id()].Contains(target.Id()) {
		return false
	}
	for edge := range g.edges {
		if edge.Source().Id() == source.Id() && edge.Target().Id() == target.Id() {
			return true
		}
	}
	return false
}

// AddEdge adds an edge with the given source and target and records it in the inbound edges and the outbound edges. Nil edges will not be added
func (g *Graph) AddEdge(edge Edge) {
	if edge == nil || g.ContainsEdge(edge.Source(), edge.Target()) {
		// the edge is nil or the graph already contains an edge which connects the (edge.Source) to the (edge.Target)
		return
	}
	g.edges[edge] = true

	if _, ok := g.outboundEdges[edge.Source().Id()]; !ok {
		g.outboundEdges[edge.Source().Id()] = &ds.IdentifierSet{Ids: make(map[string]ds.Identifier)}
	}
	g.outboundEdges[edge.Source().Id()].Add(edge.Target())

	if _, ok := g.inboundEdges[edge.Target().Id()]; !ok {
		g.inboundEdges[edge.Target().Id()] = &ds.IdentifierSet{Ids: make(map[string]ds.Identifier)}
	}
	g.inboundEdges[edge.Target().Id()].Add(edge.Source())
}

// DeleteEdge will delete the edge which connects a source vertex to the target. The related inbound and outbound edges will be deleted in the meanwhile
func (g *Graph) DeleteEdge(edge Edge) error {
	if edge == nil {
		return EdgeUnknownError{}
	}
	if !g.ContainsEdge(edge.Source(), edge.Target()) {
		return EdgeUnknownError{source: edge.Source().Id(), target: edge.Target().Id()}
	}
	delete(g.edges, edge)
	g.inboundEdges[edge.Target().Id()].Remove(edge.Source().Id())
	g.outboundEdges[edge.Source().Id()].Remove(edge.Target().Id())
	return nil
}

// visit traverses a DAG from a start node in the specified direction and returns a map of the visited nodes' id
func visit(g *DirectedAcyclicGraph, id string, visited map[string]bool, down bool) {
	var nexts *ds.IdentifierSet
	if down {
		nexts = g.outboundEdges[id]
	} else {
		nexts = g.inboundEdges[id]
	}

	if nexts != nil {
		for next := range nexts.Ids {
			if _, ok := visited[next]; ok {
				continue
			}
			visit(g, next, visited, down)
		}
	}

	visited[id] = true
}

// Vertices returns the list of all the vertices in the graph
func (g *Graph) Vertices() []Vertex {
	list := make([]Vertex, 0, g.vertices.Size())
	for _, elem := range g.vertices.Ids {
		list = append(list, elem.(Vertex))
	}
	return list
}

// Edges returns the list of all the edges in the graph
func (g *Graph) Edges() []Edge {
	list := make([]Edge, 0, len(g.edges))
	for edge := range g.edges {
		list = append(list, edge)
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
