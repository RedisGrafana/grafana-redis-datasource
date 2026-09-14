import { lastValueFrom, Observable } from 'rxjs';
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
import { ClientTypeValue, StreamingBufferLimit, StreamingDataType } from '../constants';
import { QueryTypeValue, RedisQuery } from '../redis';
import { getQuery } from '../tests/utils';
import { RedisDataSourceOptions } from '../types';
import { DataSource } from './datasource';

/**
 * Instance Settings
 */
const getInstanceSettings = (overrideSettings: object = {}): DataSourceInstanceSettings<RedisDataSourceOptions> => ({
  uid: '',
  id: 1,
  type: '',
  name: '',
  access: 'direct',
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
    cliDisabled: false,
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
              refId: 'A',
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
            type: QueryTypeValue.CLI,
            refId: 'A',
            query: '',
            streaming: true,
            streamingCapacity: 1,
            streamingInterval: 1,
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
            type: QueryTypeValue.CLI,
            refId: 'A',
            query: '',
            streaming: true,
            streamingCapacity: 1,
            streamingDataType: StreamingDataType.DATAFRAME,
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

    /**
     * Length of one streaming reading, taken and unsubscribed as a refresh would
     *
     * The length is read as the value arrives rather than returned with the
     * frame: a streaming frame keeps growing, so a length read after the next
     * reading would report that one instead.
     */
    const streamOnce = async (overrideRequest: OverrideRequest) => {
      const value = await lastValueFrom(dataSource.query(getRequest(overrideRequest)).pipe(take(1)));
      return value.data[0].length;
    };

    /**
     * Streaming target
     */
    const streamingTarget = (overrideTarget: Partial<RedisQuery> = {}): RedisQuery => ({
      type: QueryTypeValue.CLI,
      refId: 'A',
      query: '',
      streaming: true,
      streamingCapacity: 100,
      streamingInterval: 1,
      ...overrideTarget,
    });

    it('Should continue a series across the re-subscription a dashboard refresh makes', (done) => {
      const target = streamingTarget();
      const lengths: number[] = [];

      /**
       * A refresh tears the subscription down and issues the same query again as
       * a new request, so the series has to survive the gap between the two.
       */
      dataSource
        .query(getRequest({ dashboardId: 1, panelId: 1, targets: [target] }))
        .pipe(take(3))
        .subscribe({
          next: (value) => lengths.push(value.data[0].length),
          complete: () => {
            dataSource
              .query(getRequest({ dashboardId: 1, panelId: 1, targets: [target] }))
              .pipe(take(1))
              .subscribe({
                next: (value) => lengths.push(value.data[0].length),
                complete: () => {
                  expect(lengths).toEqual([1, 2, 3, 4]);
                  done();
                },
              });
          },
        });
    });

    it('Should start a new series when the interpolated query changes', async () => {
      const first = await streamOnce({
        dashboardId: 1,
        panelId: 1,
        targets: [streamingTarget({ query: 'INFO' })],
      });
      const second = await streamOnce({
        dashboardId: 1,
        panelId: 1,
        targets: [streamingTarget({ query: 'INFO' })],
      });
      const edited = await streamOnce({
        dashboardId: 1,
        panelId: 1,
        targets: [streamingTarget({ query: 'DBSIZE' })],
      });

      /**
       * The same query continues, the edited one is a series of its own
       */
      expect(first).toEqual(1);
      expect(second).toEqual(2);
      expect(edited).toEqual(1);
    });

    it('Should drop the least recently used series once the buffer limit is reached', async () => {
      await streamOnce({ dashboardId: 1, panelId: 0, targets: [streamingTarget()] });

      /**
       * Fill the cache past its limit, so that panel 0, the oldest, is evicted
       */
      for (let panelId = 1; panelId <= StreamingBufferLimit; panelId++) {
        await streamOnce({ dashboardId: 1, panelId, targets: [streamingTarget()] });
      }

      const evicted = await streamOnce({ dashboardId: 1, panelId: 0, targets: [streamingTarget()] });
      const kept = await streamOnce({
        dashboardId: 1,
        panelId: StreamingBufferLimit,
        targets: [streamingTarget()],
      });

      expect(evicted).toEqual(1);
      expect(kept).toEqual(2);
    });

    it('Should read on subscription rather than waiting out the first interval', async () => {
      /**
       * A refresh clears the timer and starts it from zero, so an interval
       * longer than the refresh period never reached its deadline and the panel
       * drew nothing at all. Jest's own timeout is what asserts it here: a
       * reading that waited for this interval could not arrive in time.
       */
      const reading = await streamOnce({
        dashboardId: 1,
        panelId: 2,
        targets: [streamingTarget({ streamingInterval: 60000 })],
      });

      expect(reading).toEqual(1);
    });

    it('Should start no timer when the interval is zero, so the dashboard refresh drives it', async () => {
      const setIntervalSpy = jest.spyOn(global, 'setInterval');

      /**
       * Grafana resubscribes on every refresh, so a subscription that reads once
       * and starts no timer reads exactly once per refresh. Two subscriptions
       * stand in for two refreshes here, and the series continues across them.
       */
      const first = await streamOnce({
        dashboardId: 1,
        panelId: 3,
        targets: [streamingTarget({ streamingInterval: 0 })],
      });
      const second = await streamOnce({
        dashboardId: 1,
        panelId: 3,
        targets: [streamingTarget({ streamingInterval: 0 })],
      });

      expect(setIntervalSpy).not.toHaveBeenCalled();
      expect(first).toEqual(1);
      expect(second).toEqual(2);

      setIntervalSpy.mockRestore();
    });

    it('Should keep the timer when only one of several streaming targets is zero', () => {
      const setIntervalSpy = jest.spyOn(global, 'setInterval');

      const subscription = dataSource
        .query(
          getRequest({
            targets: [
              streamingTarget({ streamingInterval: 0 }),
              streamingTarget({ refId: 'B', streamingInterval: 4000 }),
            ],
          })
        )
        .subscribe();

      expect(setIntervalSpy).toHaveBeenCalledWith(expect.any(Function), 4000);

      subscription.unsubscribe();
      setIntervalSpy.mockRestore();
    });

    it('Should take the streaming interval from the streaming targets alone', () => {
      const setIntervalSpy = jest.spyOn(global, 'setInterval');

      /**
       * A panel whose second target does not stream. That target carries no
       * interval, because the query editor never offered it one, and taking it
       * into account would hold the panel at the 1000ms default.
       */
      const subscription = dataSource
        .query(
          getRequest({
            targets: [
              streamingTarget({ streamingInterval: 60000 }),
              { type: QueryTypeValue.CLI, refId: 'B', query: '' },
            ],
          })
        )
        .subscribe();

      expect(setIntervalSpy).toHaveBeenCalledWith(expect.any(Function), 60000);

      subscription.unsubscribe();
      setIntervalSpy.mockRestore();
    });
  });

  /**
   * ApplyTemplateVariables Method
   */
  describe('applyTemplateVariables', () => {
    type KeyType = keyof RedisQuery;
    const testedFieldKeys: KeyType[] = [
      'keyName',
      'query',
      'searchQuery',
      'field',
      'filter',
      'legend',
      'value',
      'path',
      'cypher',
    ];

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

    it('If options.variables.datasource is null (default) should return empty array', (done) => {
      dataSource.metricFindQuery &&
        dataSource.metricFindQuery('123', { variable: { datasource: null } }).then((result: MetricFindValue[]) => {
          expect(result.length).toEqual(3);
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
