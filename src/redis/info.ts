import { SelectableValue } from '@grafana/data';

/**
 * Info sections
 */
export const InfoSections: Array<SelectableValue<string>> = [
  { label: 'Server', description: 'General information about the Redis server', value: 'server' },
  { label: 'Clients', description: 'Client connections section', value: 'clients' },
  { label: 'Memory', description: 'Memory consumption related information', value: 'memory' },
  { label: 'Persistence', description: 'RDB and AOF related information', value: 'persistence' },
  { label: 'Stats', description: 'General statistics', value: 'stats' },
  { label: 'Replication', description: 'Master/replica replication information', value: 'replication' },
  { label: 'CPU', description: 'CPU consumption statistics', value: 'cpu' },
  { label: 'Command Stats', description: 'Redis command statistics', value: 'commandstats' },
  { label: 'Cluster', description: 'Redis Cluster section', value: 'cluster' },
  { label: 'Keyspace', description: 'Database related statistics', value: 'keyspace' },
];
