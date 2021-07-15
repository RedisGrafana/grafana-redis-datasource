import { Observable } from 'rxjs';
import { map as map$, switchMap as switchMap$ } from 'rxjs/operators';
import {
  DataFrame,
  DataQueryRequest,
  DataQueryResponse,
  DataSourceInstanceSettings,
  LoadingState,
  MetricFindValue,
  ScopedVars,
} from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { DefaultStreamingInterval, StreamingDataType } from '../constants';
import { RedisQuery } from '../redis';
import { TimeSeriesStreaming } from '../time-series';
import { RedisDataSourceOptions } from '../types';

/**
 * Redis Data Source
 */
export class DataSource extends DataSourceWithBackend<RedisQuery, RedisDataSourceOptions> {
  /**
   * Constructor
   *
   * @param {DataSourceInstanceSettings<RedisDataSourceOptions>} instanceSettings Instance Settings
   */
  constructor(instanceSettings: DataSourceInstanceSettings<RedisDataSourceOptions>) {
    super(instanceSettings);
  }

  /**
   * Variable query action
   *
   * @param {string} query Query
   * @param {any} options Options
   * @returns {Promise<MetricFindValue[]>} Metric Find Values
   */
  async metricFindQuery?(query: string, options?: any): Promise<MetricFindValue[]> {
    /**
     * If query or datasource not specified
     */
    if (!query || !options.variable.datasource) {
      return Promise.resolve([]);
    }

    /**
     * Run Query
     */
    return this.query({
      targets: [{ refId: 'A', datasource: options.variable.datasource, query: query }],
    } as DataQueryRequest<RedisQuery>)
      .pipe(
        switchMap$((response) => response.data),
        switchMap$((data: DataFrame) => data.fields),
        map$((field) => {
          const values: MetricFindValue[] = [];
          field.values.toArray().forEach((value) => {
            values.push({ text: value });
          });

          return values;
        })
      )
      .toPromise();
  }

  /**
   * Override to apply template variables
   *
   * @param {string} query Query
   * @param {ScopedVars} scopedVars Scoped variables
   */
  applyTemplateVariables(query: RedisQuery, scopedVars: ScopedVars) {
    const templateSrv = getTemplateSrv();

    /**
     * Replace variables
     */
    return {
      ...query,
      keyName: query.keyName ? templateSrv.replace(query.keyName, scopedVars) : '',
      query: query.query ? templateSrv.replace(query.query, scopedVars) : '',
      field: query.field ? templateSrv.replace(query.field, scopedVars) : '',
      filter: query.filter ? templateSrv.replace(query.filter, scopedVars) : '',
      legend: query.legend ? templateSrv.replace(query.legend, scopedVars) : '',
      value: query.value ? templateSrv.replace(query.value, scopedVars) : '',
    };
  }

  /**
   * Override query to support streaming
   */
  query(request: DataQueryRequest<RedisQuery>): Observable<DataQueryResponse> {
    /**
     * No query
     * Need to typescript types narrowing
     */
    if (!request.targets.length) {
      return super.query(request);
    }

    /**
     * No streaming enabled
     */
    const streaming = request.targets.filter((target) => target.streaming);
    if (!streaming.length) {
      return super.query(request);
    }

    /**
     * Streaming enabled
     */
    return new Observable<DataQueryResponse>((subscriber) => {
      const frames: { [id: string]: TimeSeriesStreaming } = {};
      request.targets.forEach((target) => {
        /**
         * Time-series frame
         */
        if (target.streamingDataType !== StreamingDataType.DATAFRAME) {
          frames[target.refId] = new TimeSeriesStreaming(target);
        }
      });

      /**
       * Get minimum Streaming Interval
       */
      const streamingInterval = request.targets.map((target) =>
        target.streamingInterval ? target.streamingInterval : DefaultStreamingInterval
      );

      /**
       * Interval
       */
      const intervalId = setInterval(async () => {
        const response = await super.query(request).toPromise();

        response.data.forEach(async (frame) => {
          if (frames[frame.refId]) {
            frame = await frames[frame.refId].update(frame.fields);
          }

          subscriber.next({
            data: [frame],
            key: frame.refId,
            state: LoadingState.Streaming,
          });
        });
      }, Math.min(...streamingInterval));

      return () => {
        clearInterval(intervalId);
      };
    });
  }
}
