package models

/**
 * Redis Commands
 */
const (
	ClientList   = "clientList"
	ClusterInfo  = "clusterInfo"
	ClusterNodes = "clusterNodes"
	Get          = "get"
	HGet         = "hget"
	HGetAll      = "hgetall"
	HKeys        = "hkeys"
	HLen         = "hlen"
	HMGet        = "hmget"
	Info         = "info"
	LLen         = "llen"
	SCard        = "scard"
	SlowlogGet   = "slowlogGet"
	SMembers     = "smembers"
	TTL          = "ttl"
	Type         = "type"
	ZRange       = "zrange"
	XInfoStream  = "xinfoStream"
	XLen         = "xlen"
	XRange       = "xrange"
	XRevRange    = "xrevrange"
)
