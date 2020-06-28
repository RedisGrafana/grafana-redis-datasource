import { map as map$, switchMap as switchMap$ } from 'rxjs/operators';
import { DataFrame, DataQueryRequest, DataSourceInstanceSettings, MetricFindValue } from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { RedisDataSourceOptions, RedisQuery } from './types';

/**
 * Redis Data Source
 */
export class DataSource extends DataSourceWithBackend<RedisQuery, RedisDataSourceOptions> {
  /**
   * Constructor
   *
   * @param instanceSettings Instance Settings
   */
  constructor(instanceSettings: DataSourceInstanceSettings<RedisDataSourceOptions>) {
    super(instanceSettings);
  }

  /**
   * Variable query action
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
        switchMap$(response => response.data),
        switchMap$((data: DataFrame) => data.fields),
        map$(field =>
          field.values.toArray().map(value => {
            return { text: value };
          })
        )
      )
      .toPromise();
  }

  /**
   * Override to apply template variables
   */
  applyTemplateVariables(query: RedisQuery) {
    const templateSrv = getTemplateSrv();

    /**
     * Replace variables
     */
    return {
      ...query,
      key: query.key ? templateSrv.replace(query.key) : '',
      query: query.query ? templateSrv.replace(query.query) : '',
      filter: query.filter ? templateSrv.replace(query.filter) : '',
    };
  }
}
