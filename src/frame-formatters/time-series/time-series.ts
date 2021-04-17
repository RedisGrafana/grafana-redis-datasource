import { Observable } from 'rxjs';
import { map as map$, switchMap as switchMap$ } from 'rxjs/operators';
import { CircularDataFrame, DataFrame, DataQueryResponse, FieldType } from '@grafana/data';
import { DefaultStreamingCapacity } from '../../constants';
import { RedisQuery } from '../../redis';

/**
 * Time Series Formatter
 */
export class TimeSeriesFormatter {
  /**
   * Frame with all values
   */
  frame: CircularDataFrame;

  /**
   * Constructor
   *
   * @param refA
   */
  constructor(refA: RedisQuery) {
    /**
     * This dataframe can have values constantly added, and will never exceed the given capacity
     */
    this.frame = new CircularDataFrame({
      append: 'tail',
      capacity: refA?.streamingCapacity || DefaultStreamingCapacity,
    });

    /**
     * Set refId and Time field
     */
    this.frame.refId = refA?.refId;
    this.frame.addField({ name: 'time', type: FieldType.time });
  }

  /**
   * Add new values for the frame
   *
   * @param request
   */
  async update(request: Observable<DataQueryResponse>): Promise<CircularDataFrame> {
    let values: { [index: string]: number } = { time: Date.now() };

    /**
     * Fields
     */
    const fields = await request
      .pipe(
        switchMap$((response) => response.data),
        map$((data: DataFrame) => data.fields.filter((field) => (field.name === 'time' ? false : true)))
      )
      .toPromise();

    if (fields) {
      /**
       * Add fields to frame fields and return values
       */
      fields.forEach((field) => {
        /**
         * Add new fields if frame does not have the field
         */
        const fieldValues = field.values.toArray();
        const value = fieldValues[fieldValues.length - 1];

        if (!this.frame.fields.some((addedField) => addedField.name === field.name)) {
          this.frame.addField({
            name: field.name,
            type: field.type === FieldType.string && !isNaN(value) ? FieldType.number : field.type,
          });
        }

        /**
         * Set values. If values.length > 1, should be set the last line
         */
        values[field.name] = value;
      });
    }

    /**
     * Add values
     */
    this.frame.add(values);

    return Promise.resolve(this.frame);
  }
}
