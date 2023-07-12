import { FieldType } from '@grafana/data';
import { QueryTypeValue } from '../redis';
import { TimeSeriesStreaming } from './time-series';

/**
 * Time Series Streaming
 */
describe('TimeSeriesStreaming', () => {
  it('Should keep previous values', async () => {
    const frame = new TimeSeriesStreaming({ refId: 'A', type: QueryTypeValue.REDIS });
    const data = await frame.update([
      {
        name: 'value',
        type: FieldType.string,
        values: { toArray: jest.fn().mockImplementation(() => ['hello']) },
      },
    ]);
    expect(data.length).toEqual(1);

    const data2 = await frame.update([
      {
        name: 'value',
        type: FieldType.string,
        values: { toArray: jest.fn().mockImplementation(() => ['world']) },
      },
    ]);
    expect(data2.length).toEqual(2);
    expect(data2.fields[0].values.toArray()).toEqual(['hello', 'world']);
  });

  it('If no fields, should work correctly', async () => {
    const frame = new TimeSeriesStreaming({ refId: 'A', type: QueryTypeValue.REDIS });
    const data = await frame.update([
      {
        name: 'value',
        type: FieldType.string,
        values: { toArray: jest.fn().mockImplementation(() => ['hello']) },
      },
    ]);
    expect(data.length).toEqual(1);

    const data2 = await frame.update([]);
    expect(data2.length).toEqual(2);
  });

  it('Should work correctly if no query', () => {
    const frame = new TimeSeriesStreaming(undefined as any);
    expect(frame).toBeInstanceOf(TimeSeriesStreaming);
  });

  it('Should apply last line if gets more 1 line', async () => {
    const frame = new TimeSeriesStreaming({ refId: 'A', type: QueryTypeValue.REDIS });
    const data = await frame.update([
      {
        name: 'value',
        type: FieldType.string,
        values: { toArray: jest.fn().mockImplementation(() => ['hello', 'world', 'bye']) },
      },
    ]);
    expect(data.length).toEqual(1);
    expect(data.fields.length).toEqual(1);
  });

  it('Should convert string to number if value can be converted', async () => {
    const frame = new TimeSeriesStreaming({ refId: 'A', type: QueryTypeValue.REDIS, streamingCapacity: 10 });
    const data = await frame.update([
      {
        name: 'value',
        type: FieldType.string,
        values: { toArray: jest.fn().mockImplementation(() => ['123']) },
      },
    ]);
    expect(data.fields[0].name === 'value');
    expect(data.fields[0].type === FieldType.number);
    const fieldsArr = data.fields[0].values.toArray();
    expect(fieldsArr.length === 1 && fieldsArr[0] === 123);
  });
});
