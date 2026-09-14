import { SelectableValue } from '@grafana/data';

/**
 * Info Section Values
 *
 * The backend accepts any section name, so this list only controls what the
 * query editor offers. It follows the order Redis 8.10 uses in INFO EVERYTHING.
 */
export enum InfoSectionValue {
  SERVER = 'server',
  CLIENTS = 'clients',
  MEMORY = 'memory',
  PERSISTENCE = 'persistence',
  THREADS = 'threads',
  STATS = 'stats',
  REPLICATION = 'replication',
  CPU = 'cpu',
  HOTKEYS = 'hotkeys',
  MODULES = 'modules',
  COMMANDSTATS = 'commandstats',
  ERRORSTATS = 'errorstats',
  LATENCYSTATS = 'latencystats',
  CLUSTER = 'cluster',
  KEYSPACE = 'keyspace',
  KEYSIZES = 'keysizes',
  SEARCH = 'search',
  EVERYTHING = 'everything',
}

/**
 * Info sections
 */
export const InfoSections: Array<SelectableValue<InfoSectionValue>> = [
  { label: 'Server', description: 'General information about the Redis server', value: InfoSectionValue.SERVER },
  { label: 'Clients', description: 'Client connections section', value: InfoSectionValue.CLIENTS },
  { label: 'Memory', description: 'Memory consumption related information', value: InfoSectionValue.MEMORY },
  { label: 'Persistence', description: 'RDB and AOF related information', value: InfoSectionValue.PERSISTENCE },
  {
    label: 'Threads',
    description: 'Per I/O thread client, read and write counters (Redis 8.0)',
    value: InfoSectionValue.THREADS,
  },
  { label: 'Stats', description: 'General statistics', value: InfoSectionValue.STATS },
  { label: 'Replication', description: 'Master/replica replication information', value: InfoSectionValue.REPLICATION },
  { label: 'CPU', description: 'CPU consumption statistics', value: InfoSectionValue.CPU },
  {
    label: 'Hotkeys',
    description: 'Hot key tracking status, populated after HOTKEYS START (Redis 8.6)',
    value: InfoSectionValue.HOTKEYS,
  },
  {
    label: 'Modules',
    description: 'Loaded modules with version, API level and dependencies',
    value: InfoSectionValue.MODULES,
  },
  { label: 'Command Stats', description: 'Command statistics', value: InfoSectionValue.COMMANDSTATS },
  { label: 'Error Stats', description: 'Error statistics (Redis 6.2)', value: InfoSectionValue.ERRORSTATS },
  {
    label: 'Latency Stats',
    description: 'Per command latency percentiles in microseconds (Redis 7.0)',
    value: InfoSectionValue.LATENCYSTATS,
  },
  { label: 'Cluster', description: 'Cluster section', value: InfoSectionValue.CLUSTER },
  { label: 'Keyspace', description: 'Database related statistics', value: InfoSectionValue.KEYSPACE },
  {
    label: 'Key Sizes',
    description: 'Key and element size distribution histograms per database (Redis 8.0)',
    value: InfoSectionValue.KEYSIZES,
  },
  {
    label: 'Search',
    description: 'Query engine statistics, requires the search module (Redis 8.0)',
    value: InfoSectionValue.SEARCH,
  },
  {
    label: 'Everything',
    description: 'All sections at once, returned as one frame per section',
    value: InfoSectionValue.EVERYTHING,
  },
];
