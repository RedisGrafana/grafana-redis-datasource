import { SelectableValue } from '@grafana/data';

/**
 * Info Section Values
 */
export enum InfoSectionValue {
  SERVER = 'server',
  CLIENTS = 'clients',
  MEMORY = 'memory',
  PERSISTENCE = 'persistence',
  STATS = 'stats',
  REPLICATION = 'replication',
  CPU = 'cpu',
  COMMANDSTATS = 'commandstats',
  CLUSTER = 'cluster',
  KEYSPACE = 'keyspace',
}

/**
 * Info sections
 */
export const InfoSections: Array<SelectableValue<InfoSectionValue>> = [
  { label: 'Server', description: 'General information about the Redis server', value: InfoSectionValue.SERVER },
  { label: 'Clients', description: 'Client connections section', value: InfoSectionValue.CLIENTS },
  { label: 'Memory', description: 'Memory consumption related information', value: InfoSectionValue.MEMORY },
  { label: 'Persistence', description: 'RDB and AOF related information', value: InfoSectionValue.PERSISTENCE },
  { label: 'Stats', description: 'General statistics', value: InfoSectionValue.STATS },
  { label: 'Replication', description: 'Master/replica replication information', value: InfoSectionValue.REPLICATION },
  { label: 'CPU', description: 'CPU consumption statistics', value: InfoSectionValue.CPU },
  { label: 'Command Stats', description: 'Redis command statistics', value: InfoSectionValue.COMMANDSTATS },
  { label: 'Cluster', description: 'Redis Cluster section', value: InfoSectionValue.CLUSTER },
  { label: 'Keyspace', description: 'Database related statistics', value: InfoSectionValue.KEYSPACE },
];
