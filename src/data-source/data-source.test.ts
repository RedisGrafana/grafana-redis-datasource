import { Observable } from 'rxjs';
import { take } from 'rxjs/operators';
import {
  CircularDataFrame,
  DataQueryRequest,
  DataSourceInstanceSettings,
  FieldType,
  MetricFindValue,
  PluginType,
  toDataFrame,
} from '@grafana/data';
import { DataSourceWithBackend, setTemplateSrv, TemplateSrv } from '@grafana/runtime';
import { ClientTypeValue, StreamingDataType } from '../constants';
import { QueryTypeValue, RedisQuery } from '../redis';
import { getQuery } from '../tests/utils';
import { RedisDataSourceOptions } from '../types';
import { DataSource } from './data-source';

/**
 * Instance Settings
 */
const getInstanceSettings = (overrideSettings: object = {}): DataSourceInstanceSettings<RedisDataSourceOptions> => ({
  uid: '',
  id: 1,
  type: '',
  name: '',
  meta: {
    id: '',
    name: '',
    type: PluginType.datasource,
    info: {} as any,
    module: '',
    baseUrl: '',
  },
  jsonData: {
    poolSize: 0,
    timeout: 0,
    pingInterval: 0,
    pipelineWindow: 0,
    tlsAuth: false,
    tlsSkipVerify: false,
    client: ClientTypeValue.CLUSTER,
    sentinelName: '',
    sentinelUser: '',
    sentinelAcl: false,
    acl: false,
    user: '',
  },
});

/**
 * Override Request
 */
interface OverrideRequest {
  [key: string]: unknown;
  targets?: RedisQuery[];
}

/**
 * Request
 */
const getRequest = (overrideRequest: OverrideRequest = {}): DataQueryRequest<RedisQuery> => ({
  requestId: '',
  interval: '',
  intervalMs: 0,
  range: {} as any,
  scopedVars: {},
  timezone: '',
  app: '',
  startTime: 0,
  ...overrideRequest,
  targets: overrideRequest.targets
    ? overrideRequest.targets
    : [
        {
          datasource: '',
          type: QueryTypeValue.CLI,
          refId: 'A',
          query: '',
          streaming: false,
        },
      ],
});

/**
 * Data Source
 */
describe('DataSource', () => {
  const superQueryMock = jest.spyOn(DataSourceWithBackend.prototype, 'query').mockImplementation(
    () =>
      new Observable((subscriber) => {
        subscriber.next({
          data: [
            toDataFrame({
              fields: [
                {
                  name: 'get',
                  type: FieldType.number,
                  values: [1, 2, 3],
                },
              ],
            }),
          ],
        });
        subscriber.complete();
      })
  );
  const instanceSettings = getInstanceSettings();
  let dataSource: DataSource;
  let templateSrv: TemplateSrv;

  beforeEach(() => {
    superQueryMock.mockClear();
    dataSource = new DataSource(instanceSettings);
    templateSrv = {
      replace: jest.fn().mockImplementation((value) => `replaced:${value}`),
      getVariables: jest.fn(),
    };
    setTemplateSrv(templateSrv);
  });

  /**
   * Query Method
   */
  describe('query', () => {
    it('If no streaming should use super.query', (done) => {
      const request = getRequest();
      dataSource.query(request).subscribe(() => {
        expect(superQueryMock).toHaveBeenCalledWith(request);
        done();
      });
    });

    it('If no query should use super.query', (done) => {
      const request = getRequest({ targets: [] });
      dataSource.query(request).subscribe(() => {
        expect(superQueryMock).toHaveBeenCalledWith(request);
        done();
      });
    });

    it('Should use TimeSeries as default streamingDataType', (done) => {
      const request = getRequest({
        targets: [
          {
            datasource: '',
            type: QueryTypeValue.CLI,
            refId: 'A',
            query: '',
            streaming: true,
            streamingCapacity: 1,
          },
        ],
      });

      /**
       * Query
       */
      dataSource
        .query(request)
        .pipe(take(3))
        .subscribe(
          (value) => {
            value.data.forEach((item) => {
              expect(item).toBeInstanceOf(CircularDataFrame);
            });
          },
          null,
          () => {
            done();
          }
        );
    });

    it('If streaming exists should get frames in interval', (done) => {
      const request = getRequest({
        targets: [
          {
            datasource: '',
            type: QueryTypeValue.CLI,
            refId: 'A',
            query: '',
            streaming: true,
            streamingCapacity: 1,
            streamingDataType: StreamingDataType.DataFrame,
          },
        ],
      });

      /**
       * Query
       */
      dataSource
        .query(request)
        .pipe(take(3))
        .subscribe(
          (value) => {
            value.data.forEach((item) => {
              expect(item).not.toBeInstanceOf(CircularDataFrame);
              expect(item.fields).toBeDefined();
            });
          },
          null,
          () => {
            done();
          }
        );
    });
  });

  /**
   * ApplyTemplateVariables Method
   */
  describe('applyTemplateVariables', () => {
    type KeyType = keyof RedisQuery;
    const testedFieldKeys: KeyType[] = ['keyName', 'query', 'field', 'filter', 'legend', 'value'];

    testedFieldKeys.forEach((fieldKey) => {
      describe(fieldKey, () => {
        it('If value is exist should replace via templateSrv', () => {
          const value = 'filter';
          const query = getQuery({ [fieldKey]: value });
          const scopedVars = { keyName: { value: 'key', key: '', text: '' } };
          const resultQuery = dataSource.applyTemplateVariables(query, scopedVars);
          expect(templateSrv.replace).toHaveBeenCalledWith(query[fieldKey], scopedVars);
          expect(resultQuery).toEqual({
            ...query,
            [fieldKey]: `replaced:${value}`,
          });
        });

        it('If value is empty should not replace via templateSrv', () => {
          const value = '';
          const query = getQuery({ [fieldKey]: value });
          const scopedVars = { keyName: { value: 'key', key: '', text: '' } };
          const resultQuery = dataSource.applyTemplateVariables(query, scopedVars);
          expect(templateSrv.replace).not.toHaveBeenCalled();
          expect(resultQuery).toEqual({
            ...query,
            [fieldKey]: value,
          });
        });
      });
    });
  });

  /**
   * MetricFindQuery Method
   */
  describe('metricFindQuery', () => {
    it('If query is empty should return empty array', (done) => {
      dataSource.metricFindQuery &&
        dataSource.metricFindQuery('', { variable: { datasource: '123' } }).then((result: MetricFindValue[]) => {
          expect(result).toEqual([]);
          done();
        });
    });

    it('If options.variables.datasource is empty should return empty array', (done) => {
      dataSource.metricFindQuery &&
        dataSource.metricFindQuery('123', { variable: { datasource: '' } }).then((result: MetricFindValue[]) => {
          expect(result).toEqual([]);
          done();
        });
    });

    it('Should call query method', (done) => {
      const querySpyMethod = jest.spyOn(dataSource, 'query').mockImplementation(
        () =>
          new Observable((subscriber) => {
            subscriber.next({
              data: [
                {
                  fields: [
                    {
                      name: 'get',
                      values: {
                        toArray() {
                          return ['1', '2', '3'];
                        },
                      },
                    },
                  ],
                  length: 1,
                },
              ],
            });
            subscriber.complete();
          })
      );

      dataSource.metricFindQuery &&
        dataSource.metricFindQuery('123', { variable: { datasource: '123' } }).then((result: MetricFindValue[]) => {
          expect(querySpyMethod).toHaveBeenCalled();
          expect(result).toEqual([{ text: '1' }, { text: '2' }, { text: '3' }]);
          done();
        });
    });
  });

  it('Should call query method with numbers', (done) => {
    const querySpyMethod = jest.spyOn(dataSource, 'query').mockImplementation(
      () =>
        new Observable((subscriber) => {
          subscriber.next({
            data: [
              {
                fields: [
                  {
                    name: 'get',
                    values: {
                      toArray() {
                        return new Float64Array([21, 31]);
                      },
                    },
                  },
                ],
                length: 1,
              },
            ],
          });
          subscriber.complete();
        })
    );

    dataSource.metricFindQuery &&
      dataSource.metricFindQuery('123', { variable: { datasource: '123' } }).then((result: MetricFindValue[]) => {
        expect(querySpyMethod).toHaveBeenCalled();
        expect(result).toEqual([{ text: 21 }, { text: 31 }]);
        done();
      });
  });

  afterAll(() => {
    superQueryMock.mockReset();
    setTemplateSrv(null as any);
  });
});
