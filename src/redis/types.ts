import { DataQuery, SelectableValue } from '@grafana/data';
import { InfoSectionValue } from './info';
import { QueryTypeValue } from './query';
import { AggregationValue } from './time-series';

export enum StreamingDataType {
  TimeSeries = 'TimeSeries',
  DataFrame = 'DataFrame',
}

export const StreamingDataTypes: Array<SelectableValue<StreamingDataType>> = [
  {
    label: 'Time series',
    value: StreamingDataType.TimeSeries,
  },
  {
    label: 'Data frame',
    value: StreamingDataType.DataFrame,
  }
]

/**
 * Redis Query
 */
export interface RedisQuery extends DataQuery {
    /**
     * Type
     *
     * @type {string}
     */
    type: QueryTypeValue;

    /**
     * Query command
     *
     * @type {string}
     */
    query?: string;

    /**
     * Field
     *
     * @type {string}
     */
    field?: string;

    /**
     * Redis TimeSeries filter
     *
     * @see https://oss.redislabs.com/redistimeseries/commands/#filtering
     * @type {string}
     */
    filter?: string;

    /**
     * Command
     *
     * @type {string}
     */
    command?: string;

    /**
     * Key name
     *
     * @type {string}
     */
    keyName?: string;

    /**
     * Value label
     *
     * @type {string}
     */
    value?: string;

    /**
     * Aggregation
     *
     * @see https://oss.redislabs.com/redistimeseries/commands/#aggregation-compaction-downsampling
     * @type {string}
     */
    aggregation?: AggregationValue;

    /**
     * Bucket
     *
     * @type {number}
     */
    bucket?: number;

    /**
     * Fill
     *
     * @type {boolean}
     */
    fill?: boolean;

    /**
     * Legend label
     *
     * @type {string}
     */
    legend?: string;

    /**
     * Info Section
     *
     * @type {string}
     */
    section?: InfoSectionValue;

    /**
     * Size
     *
     * @type {number}
     */
    size?: number;

    /**
     * Support streaming
     *
     * @type {boolean}
     */
    streaming?: boolean;

    /**
     * Streaming interval in milliseconds
     *
     * @type {number}
     */
    streamingInterval?: number;

    /**
     * Streaming capacity
     *
     * @type {number}
     */
    streamingCapacity?: number;

    /**
     * Streaming data type
     * @type {StreamingDataType}
     */
    streamingDataType?: StreamingDataType;

    /**
     * Cursor for SCAN command
     *
     * @type {string}
     */
    cursor?: string;

    /**
     * Match for SCAN command
     *
     * @type {string}
     */
    match?: string;

    /**
     * Count for SCAN command
     *
     * @type {number}
     */
    count?: number;

    /**
     * Samples for MEMORY USAGE command
     *
     * @type {number}
     */
    samples?: number;

    /**
     * Start for Streams
     *
     * @type {string}
     */
    start?: string;

    /**
     * Stop for Streams
     *
     * @type {string}
     */
    end?: string;
}
