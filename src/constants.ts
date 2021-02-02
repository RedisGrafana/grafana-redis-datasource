import { SelectableValue } from '@grafana/data';

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

/**
 * Streaming Data Type
 */
export enum StreamingDataType {
  TimeSeries = 'TimeSeries',
  DataFrame = 'DataFrame',
}

/**
 * Streaming
 */
export const StreamingDataTypes: Array<SelectableValue<StreamingDataType>> = [
  {
    label: 'Time series',
    value: StreamingDataType.TimeSeries,
  },
  {
    label: 'Data frame',
    value: StreamingDataType.DataFrame,
  },
];
