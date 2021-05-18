package models

/**
 * Commands
 */
const GraphConfig = "graph.config"
const GraphExplain = "graph.explain"
const GraphProfile = "graph.profile"
const GraphQuery = "graph.query"
const GraphSlowlog = "graph.slowlog"

/**
 *  Represents node
 */
type NodeEntry struct {
	Id       string
	Title    string
	SubTitle string
	MainStat string
	Arc      int64
}

/**
 *  Represents edge
 */
type EdgeEntry struct {
	Id       string
	Source   string
	Target   string
	MainStat string
}
