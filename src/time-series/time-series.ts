import { CircularDataFrame, Field, FieldType } from '@grafana/data';
import { DefaultStreamingCapacity } from '../constants';
import { RedisQuery } from '../redis';

/**
 * Time Series Streaming
 */
export class TimeSeriesStreaming {
  /**
   * Frame with all values
   */
  frame: CircularDataFrame;

  /**
   * Constructor
   *
   * @param ref
   */
  constructor(ref: RedisQuery) {
    /**
     * This dataframe can have values constantly added, and will never exceed the given capacity
     */
    this.frame = new CircularDataFrame({
      append: 'tail',
      capacity: ref?.streamingCapacity || DefaultStreamingCapacity,
    });

    /**
     * Set refId
     */
    this.frame.refId = ref?.refId;
  }

  /**
   * Add new values for the frame
   *
   * @param request
   */
  async update(fields: any): Promise<CircularDataFrame> {
    let values: { [index: string]: number } = {};

    /**
     * Add fields to frame fields and return values
     */
    fields.forEach((field: Field) => {
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

    /**
     * Add values and return
     */
    this.frame.add(values);
    return Promise.resolve(this.frame);
  }
}
