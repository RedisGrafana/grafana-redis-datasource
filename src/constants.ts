/**
 * Default Streaming Interval
 */
export const DefaultStreamingInterval = 1000;

/**
 * Default Streaming Capacity
 */
export const DefaultStreamingCapacity = 1000;

/**
 * Client Type Values
 */
export enum ClientTypeValue {
  CLUSTER = 'cluster',
  SENTINEL = 'sentinel',
  SOCKET = 'socket',
  STANDALONE = 'standalone',
}

/**
 * Client Types
 */
export const ClientType = [
  { label: 'Standalone', value: ClientTypeValue.STANDALONE },
  { label: 'Cluster', value: ClientTypeValue.CLUSTER },
  { label: 'Sentinel', value: ClientTypeValue.SENTINEL },
  { label: 'Socket', value: ClientTypeValue.SOCKET },
];
