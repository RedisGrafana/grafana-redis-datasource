import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { RedisDataSourceOptions, RedisQuery } from './types';

export interface DataQueryFilter {
  /**
   * Reference Id
   */
  refId: string;

  /**
   * Redis key
   *
   * @type {string}
   */
  key: string;

  /**
   * Redis TimeSeries filter
   *
   * @see https://oss.redislabs.com/redistimeseries/commands/#filtering
   * @type {string}
   */
  filter: string;
}

export class DataSource extends DataSourceWithBackend<RedisQuery, RedisDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<RedisDataSourceOptions>) {
    super(instanceSettings);
  }

  /**
   * Override to apply template variables
   */
  applyTemplateVariables(query: DataQueryFilter) {
    /**
     * Replace variables in Key
     */
    if (query.key) {
      query.key = getTemplateSrv().replace(query.key);
    }

    /**
     * Replace veriables in filter
     */
    if (query.filter) {
      query.filter = getTemplateSrv().replace(query.filter);
    }

    /**
     * Return
     */
    return query;
  }
}
