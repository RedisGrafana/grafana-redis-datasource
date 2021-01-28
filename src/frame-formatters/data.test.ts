import { Observable } from 'rxjs';
import { FieldType, toDataFrame } from '@grafana/data';
import { DataFrameFormatter } from './data';

/**
 * Data Frame Formatter
 */
describe('DataFrameFormatter', () => {
  it('Should return dataFrame', async () => {
    const frame = new DataFrameFormatter();
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
    expect(data.fields[0].name === 'value');
    expect(data.fields[0].values.toArray() === ['hello']);
  });
});
