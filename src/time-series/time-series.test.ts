import { FieldType } from '@grafana/data';
import { QueryTypeValue } from '../redis';
import { TimeSeriesStreaming, TimeFieldName } from './time-series';

/**
 * Returns a field by name, so assertions do not depend on the position of the time field
 */
const fieldByName = (frame: any, name: string) => frame.fields.find((field: any) => field.name === name);

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
    expect(fieldByName(data2, 'value').values.toArray()).toEqual(['hello', 'world']);
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
    expect(data.fields.length).toEqual(2);
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
    const value = fieldByName(data, 'value');
    expect(value.type).toEqual(FieldType.number);
    expect(value.values.toArray()).toEqual([123]);
  });

  it('Should add a time field so the frame can be plotted', async () => {
    const frame = new TimeSeriesStreaming({ refId: 'A', type: QueryTypeValue.REDIS });
    const before = Date.now();
    const data = await frame.update([
      {
        name: 'used_memory',
        type: FieldType.number,
        values: { toArray: jest.fn().mockImplementation(() => [4686704]) },
      },
    ]);

    const time = fieldByName(data, TimeFieldName);
    expect(time).toBeDefined();
    expect(time.type).toEqual(FieldType.time);

    const timestamps = time.values.toArray();
    expect(timestamps.length).toEqual(1);
    expect(timestamps[0]).toBeGreaterThanOrEqual(before);
    expect(timestamps[0]).toBeLessThanOrEqual(Date.now());
  });

  it('Should keep a time field returned by the command', async () => {
    const frame = new TimeSeriesStreaming({ refId: 'A', type: QueryTypeValue.REDIS });
    const data = await frame.update([
      {
        name: 'Timestamp',
        type: FieldType.time,
        values: { toArray: jest.fn().mockImplementation(() => [1789052949000]) },
      },
      {
        name: 'Duration',
        type: FieldType.number,
        values: { toArray: jest.fn().mockImplementation(() => [13]) },
      },
    ]);

    expect(fieldByName(data, TimeFieldName)).toBeUndefined();
    expect(fieldByName(data, 'Timestamp').values.toArray()).toEqual([1789052949000]);
  });
});
