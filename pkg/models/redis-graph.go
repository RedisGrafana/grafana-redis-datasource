package models

/**
 * RedisGraph Commands
 */
const (
	GraphConfig  = "graph.config"
	GraphExplain = "graph.explain"
	GraphProfile = "graph.profile"
	GraphQuery   = "graph.query"
	GraphSlowlog = "graph.slowlog"
)

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
