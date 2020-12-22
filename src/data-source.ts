import { head } from 'lodash';
import { Observable } from 'rxjs';
import { map as map$, switchMap as switchMap$ } from 'rxjs/operators';
import {
  CircularDataFrame,
  DataFrame,
  DataQueryRequest,
  DataQueryResponse,
  DataSourceInstanceSettings,
  FieldType,
  MetricFindValue,
  ScopedVars,
} from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { RedisQuery } from './redis';
import { RedisDataSourceOptions } from './types';

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
     * No streaming enabled
     */
    if (!refA?.streaming) {
      return super.query(request);
    }

    /**
     * Streaming enabled
     */
    return new Observable<DataQueryResponse>((subscriber) => {
      /**
       * This dataframe can have values constantly added, and will never exceed the given capacity
       */
      const frame = new CircularDataFrame({
        append: 'tail',
        capacity: refA?.streamingCapacity || 1000,
      });

      /**
       * Set refId and Time field
       */
      frame.refId = refA.refId;
      frame.addField({ name: 'time', type: FieldType.time });

      /**
       * Interval
       */
      const intervalId = setInterval(async () => {
        let values: { [index: string]: number } = { time: Date.now() };

        /**
         * Run Query and filter time field out
         */
        const fields = await super
          .query(request)
          .pipe(
            switchMap$((response) => response.data),
            map$((data: DataFrame) => data.fields.filter((field) => (field.name === 'time' ? false : true)))
          )
          .toPromise();

        if (fields) {
          /**
           * Add fields to frame fields and return values
           */
          fields.map((field) =>
            field.values.toArray().map((value) => {
              if (frame.fields.length < fields.length + 1) {
                frame.addField({
                  name: field.name,
                  type: field.type === FieldType.string && !isNaN(value) ? FieldType.number : field.type,
                });
              }
              return (values[field.name] = value);
            })
          );
        }

        /**
         * Add frame
         */
        frame.add(values);
        subscriber.next({
          data: [frame],
          key: refA.refId,
        });
      }, refA.streamingInterval || 1000);

      return () => {
        clearInterval(intervalId);
      };
    });
  }
}
