/**
 * Supported Commands
 */
export enum RedisGraph {
  CONFIG = 'graph.config',
  PROFILE = 'graph.profile',
  EXPLAIN = 'graph.explain',
  QUERY = 'graph.query',
  SLOWLOG = 'graph.slowlog',
}

/**
 * Commands List
 */
export const RedisGraphCommands = [
  {
    label: RedisGraph.CONFIG.toUpperCase(),
    description: 'Retrieves a RedisGraph configuration',
    value: RedisGraph.CONFIG,
  },
  {
    label: RedisGraph.EXPLAIN.toUpperCase(),
    description: 'Constructs a query execution plan but does not run it',
    value: RedisGraph.EXPLAIN,
  },
  {
    label: RedisGraph.PROFILE.toUpperCase(),
    description:
      "Executes a query and produces an execution plan augmented with metrics for each operation's execution",
    value: RedisGraph.PROFILE,
  },
  {
    label: RedisGraph.QUERY.toUpperCase(),
    description: 'Executes the given query against a specified graph',
    value: RedisGraph.QUERY,
  },
  {
    label: RedisGraph.SLOWLOG.toUpperCase(),
    description: 'Returns a list containing up to 10 of the slowest queries issued against the given graph ID',
    value: RedisGraph.SLOWLOG,
  },
];
