import { head } from 'lodash';
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
import { DataFrameFormatter, TimeSeriesFormatter } from '../frame-formatters';
import { RedisQuery } from '../redis';
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
      targets: [{ datasource: options.variable.datasource, query: query }],
    } as DataQueryRequest<RedisQuery>)
      .pipe(
        switchMap$((response) => response.data),
        switchMap$((data: DataFrame) => data.fields),
        map$((field) =>
          field.values.toArray().map((value) => {
            return { text: value };
          })
        )
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
    const refA = head(request.targets);

    /**
     * No query
     * Need to typescript types narrowing
     */
    if (!refA) {
      return super.query(request);
    }

    /**
     * No streaming enabled
     */
    if (!refA?.streaming) {
      return super.query(request);
    }

    /**
     * Streaming enabled
     */
    return new Observable<DataQueryResponse>((subscriber) => {
      const { streamingDataType = StreamingDataType.TimeSeries } = refA;

      /**
       * Apply frame formatted by streamingDataType
       */
      let frame: TimeSeriesFormatter | DataFrameFormatter = new TimeSeriesFormatter(refA);
      if (streamingDataType === StreamingDataType.DataFrame) {
        frame = new DataFrameFormatter();
      }

      /**
       * Interval
       */
      const intervalId = setInterval(async () => {
        /**
         * Run Query
         */
        const data = await frame.update(super.query(request));
        if (!data) {
          return;
        }

        subscriber.next({
          data: [data],
          key: refA.refId,
          state: LoadingState.Streaming,
        });
      }, refA.streamingInterval || DefaultStreamingInterval);

      return () => {
        clearInterval(intervalId);
      };
    });
  }
}
