package models

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
