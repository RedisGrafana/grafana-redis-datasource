import { Observable } from 'rxjs';
import { FieldType, toDataFrame } from '@grafana/data';
import { QueryTypeValue } from '../../redis';
import { TimeSeriesFormatter } from './time-series';

/**
 * Time Series Formatter
 */
describe('TimeSeriesFormatter', () => {
  it('Should keep previous values', async () => {
    const frame = new TimeSeriesFormatter({ refId: 'A', type: QueryTypeValue.REDIS });
    const data = await frame.update(
      new Observable((subscriber) => {
        subscriber.next({
          data: [
            toDataFrame({
              fields: [
                {
                  name: 'value',
                  type: FieldType.string,
                  values: ['hello'],
                },
              ],
            }),
          ],
        });
        subscriber.complete();
      })
    );
    expect(data.length).toEqual(1);

    const data2 = await frame.update(
      new Observable((subscriber) => {
        subscriber.next({
          data: [
            toDataFrame({
              fields: [
                {
                  name: 'value',
                  type: FieldType.string,
                  values: ['world'],
                },
              ],
            }),
          ],
        });
        subscriber.complete();
      })
    );
    expect(data2.length).toEqual(2);
    expect(data2.fields[0].values.toArray()).toEqual(['hello', 'world']);
  });

  it('If no fields, should work correctly', async () => {
    const frame = new TimeSeriesFormatter({ refId: 'A', type: QueryTypeValue.REDIS });
    const data = await frame.update(
      new Observable((subscriber) => {
        subscriber.next({
          data: [
            toDataFrame({
              fields: [
                {
                  name: 'value',
                  type: FieldType.string,
                  values: ['hello'],
                },
              ],
            }),
          ],
        });
        subscriber.complete();
      })
    );
    expect(data.length).toEqual(1);

    const data2 = await frame.update(
      new Observable((subscriber) => {
        subscriber.next({
          data: [],
        });
        subscriber.complete();
      })
    );
    expect(data2.length).toEqual(2);
  });

  it('Should work correctly if no query', () => {
    const frame = new TimeSeriesFormatter(undefined as any);
    expect(frame).toBeInstanceOf(TimeSeriesFormatter);
  });

  it('Should apply last line if gets more 1 line', async () => {
    const frame = new TimeSeriesFormatter({ refId: 'A', type: QueryTypeValue.REDIS });
    const data = await frame.update(
      new Observable((subscriber) => {
        subscriber.next({
          data: [
            toDataFrame({
              fields: [
                {
                  name: 'value',
                  type: FieldType.string,
                  values: ['hello', 'world', 'bye'],
                },
              ],
            }),
          ],
        });
        subscriber.complete();
      })
    );
    expect(data.length).toEqual(1);
    expect(data.fields.length).toEqual(1);
  });

  it('Should convert string to number if value can be converted', async () => {
    const frame = new TimeSeriesFormatter({ refId: 'A', type: QueryTypeValue.REDIS, streamingCapacity: 10 });
    const data = await frame.update(
      new Observable((subscriber) => {
        subscriber.next({
          data: [
            toDataFrame({
              fields: [
                {
                  name: 'value',
                  type: FieldType.string,
                  values: ['123'],
                },
              ],
            }),
          ],
        });
        subscriber.complete();
      })
    );
    expect(data.fields[0].name === 'value');
    expect(data.fields[0].type === FieldType.number);
    expect(data.fields[0].values.toArray() === [123]);
  });
});
