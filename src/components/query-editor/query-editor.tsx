import { css } from 'emotion';
import React, { ChangeEvent, PureComponent } from 'react';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { Button, InlineFormLabel, LegacyForms, RadioButtonGroup, Select, TextArea } from '@grafana/ui';
import { DataSource } from '../../data-source';
import {
  Aggregations,
  AggregationValue,
  CommandParameters,
  Commands,
  InfoSections,
  InfoSectionValue,
  QueryType,
  QueryTypeValue,
  RedisQuery,
  StreamingDataType,
  StreamingDataTypes,
} from '../../redis';
import { RedisDataSourceOptions } from '../../types';

/**
 * Form Field
 */
const { FormField, Switch } = LegacyForms;

/**
 * Editor Property
 */
type Props = QueryEditorProps<DataSource, RedisQuery, RedisDataSourceOptions>;

/**
 * Query Editor
 */
export class QueryEditor extends PureComponent<Props> {
  /**
   * Change handler for number field
   *
   * @param {ChangeEvent<HTMLInputElement>} event Event
   */
  createNumberFieldHandler = (name: keyof RedisQuery) => (event: ChangeEvent<HTMLInputElement>) => {
    this.props.onChange({ ...this.props.query, [name]: Number(event.target.value) });
  };

  /**
   * Change handler for text field
   *
   * @param {ChangeEvent<HTMLInputElement>} event Event
   */
  createTextFieldHandler = (name: keyof RedisQuery) => (event: ChangeEvent<HTMLInputElement>) => {
    this.props.onChange({ ...this.props.query, [name]: event.target.value });
  };

  /**
   * Change handler for textarea field
   *
   * @param {ChangeEvent<HTMLInputElement>} event Event
   */
  createTextareaFieldHandler = (name: keyof RedisQuery) => (event: ChangeEvent<HTMLTextAreaElement>) => {
    this.props.onChange({ ...this.props.query, [name]: event.target.value });
  };

  /**
   * Change handler for select field
   *
   * @param {ChangeEvent<HTMLInputElement>} event Event
   */
  createSelectFieldHandler<ValueType>(name: keyof RedisQuery) {
    return (val: SelectableValue<ValueType>) => {
      this.props.onChange({ ...this.props.query, [name]: val.value });
    };
  }

  /**
   * Change handler for radio button field
   *
   * @param {value: ValueType}
   */
  createRedioButtonFieldHandler<ValueType>(name: keyof RedisQuery) {
    return (value?: ValueType) => {
      this.props.onChange({ ...this.props.query, [name]: value });
    };
  }

  /**
   * Change handler for switch field
   *
   * @param {ChangeEvent<HTMLInputElement>} event Event
   */
  createSwitchFieldHandler = (name: keyof RedisQuery) => (event: React.SyntheticEvent<HTMLInputElement>) => {
    this.props.onChange({ ...this.props.query, [name]: event.currentTarget.checked });
  };

  /**
   * Key name change
   */
  onKeyNameChange = this.createTextFieldHandler('keyName');

  /**
   * Query change
   */
  onQueryChange = this.createTextareaFieldHandler('query');

  /**
   * Filter change
   */
  onFilterChange = this.createTextFieldHandler('filter');

  /**
   * Field change
   */
  onFieldChange = this.createTextFieldHandler('field');

  /**
   * Legend change
   */
  onLegendChange = this.createTextFieldHandler('legend');

  /**
   * Value change
   */
  onValueChange = this.createTextFieldHandler('value');

  /**
   * Command change
   */
  onCommandChange = this.createSelectFieldHandler<string>('command');

  /**
   * Type change
   *
   * @param val Value
   */
  onTypeChange = (val: SelectableValue<string>) => {
    const { onChange, query } = this.props;
    onChange({
      ...query,
      type: val.value as QueryTypeValue,
      query: '',
      command: '',
    });
  };

  /**
   * Aggregation change
   */
  onAggregationChange = this.createSelectFieldHandler<AggregationValue>('aggregation');

  /**
   * Info section change
   */
  onInfoSectionChange = this.createSelectFieldHandler<InfoSectionValue>('section');

  /**
   * Bucket change
   */
  onBucketChange = this.createNumberFieldHandler('bucket');

  /**
   * Size change
   */
  onSizeChange = this.createNumberFieldHandler('size');

  /**
   * Count change
   */
  onCountChange = this.createNumberFieldHandler('count');

  /**
   * Samples change
   */
  onSamplesChange = this.createNumberFieldHandler('samples');

  /**
   * Cursor change
   */
  onCursorChange = this.createTextFieldHandler('cursor');

  /**
   * Match change
   */
  onMatchChange = this.createTextFieldHandler('match');

  /**
   * Start change
   */
  onStartChange = this.createTextFieldHandler('start');

  /**
   * End change
   */
  onEndChange = this.createTextFieldHandler('end');

  /**
   * Fill change
   */
  onFillChange = this.createSwitchFieldHandler('fill');

  /**
   * Streaming change
   */
  onStreamingChange = this.createSwitchFieldHandler('streaming');

  /**
   * Streaming interval change
   */
  onStreamingIntervalChange = this.createNumberFieldHandler('streamingInterval');

  /**
   * Streaming capacity change
   */
  onStreamingCapacityChange = this.createNumberFieldHandler('streamingCapacity');

  /**
   * Streaming data type change
   */
  onStreamingDataTypeChange = this.createRedioButtonFieldHandler<StreamingDataType>('streamingDataType');

  /**
   * Render Editor
   */
  render() {
    const {
      keyName,
      aggregation,
      bucket,
      legend,
      command,
      field,
      filter,
      value,
      query,
      type,
      section,
      size,
      fill,
      cursor,
      count,
      match,
      samples,
      start,
      end,
      streaming,
      streamingInterval,
      streamingCapacity,
      streamingDataType,
      refId,
    } = this.props.query;
    const { onRunQuery } = this.props;

    /**
     * Return
     */
    return (
      <div>
        <div className="gf-form">
          <InlineFormLabel width={8}>Type</InlineFormLabel>
          <Select
            className={css`
              margin-right: 5px;
            `}
            width={40}
            options={QueryType}
            menuPlacement="bottom"
            value={type}
            onChange={this.onTypeChange}
          />

          {type === QueryTypeValue.CLI && (
            <>
              <InlineFormLabel width={8}>Command</InlineFormLabel>
              <TextArea css="" value={query} className="gf-form-input" onChange={this.onQueryChange} />
            </>
          )}
          {type && type !== QueryTypeValue.CLI && (
            <>
              <InlineFormLabel width={8}>Command</InlineFormLabel>
              <Select options={Commands[type]} menuPlacement="bottom" value={command} onChange={this.onCommandChange} />
            </>
          )}
        </div>

        {type !== QueryTypeValue.CLI && command && (
          <div className="gf-form">
            {CommandParameters.keyName.includes(command) && (
              <FormField
                labelWidth={8}
                inputWidth={30}
                value={keyName}
                onChange={this.onKeyNameChange}
                label="Key"
                tooltip="Key name"
              />
            )}

            {CommandParameters.filter.includes(command) && (
              <FormField
                labelWidth={8}
                inputWidth={30}
                value={filter}
                onChange={this.onFilterChange}
                label="Label Filter"
                tooltip="Label Filter"
              />
            )}

            {CommandParameters.field.includes(command) && (
              <FormField labelWidth={8} inputWidth={30} value={field} onChange={this.onFieldChange} label="Field" />
            )}

            {CommandParameters.legend.includes(command) && (
              <FormField
                labelWidth={8}
                inputWidth={20}
                value={legend}
                onChange={this.onLegendChange}
                label="Legend"
                tooltip="Legend override"
              />
            )}

            {CommandParameters.value.includes(command) && (
              <FormField
                labelWidth={8}
                inputWidth={10}
                value={value}
                onChange={this.onValueChange}
                label="Value"
                tooltip="Value override"
              />
            )}

            {CommandParameters.legendLabel.includes(command) && (
              <FormField
                labelWidth={8}
                inputWidth={10}
                value={legend}
                onChange={this.onLegendChange}
                label="Legend Label"
              />
            )}

            {CommandParameters.valueLabel.includes(command) && (
              <FormField
                labelWidth={8}
                inputWidth={10}
                value={value}
                onChange={this.onValueChange}
                label="Value Label"
              />
            )}

            {CommandParameters.size.includes(command) && (
              <FormField
                labelWidth={8}
                inputWidth={10}
                value={size}
                type="number"
                onChange={this.onSizeChange}
                label="Size"
              />
            )}

            {CommandParameters.cursor.includes(command) && (
              <FormField labelWidth={8} inputWidth={10} value={cursor} onChange={this.onCursorChange} label="Cursor" />
            )}

            {CommandParameters.match.includes(command) && (
              <FormField
                labelWidth={8}
                inputWidth={10}
                value={match}
                onChange={this.onMatchChange}
                placeholder="*"
                label="Match pattern"
              />
            )}

            {CommandParameters.start.includes(command) && (
              <FormField
                labelWidth={8}
                inputWidth={10}
                value={start}
                onChange={this.onStartChange}
                placeholder="-"
                label="Start"
              />
            )}

            {CommandParameters.end.includes(command) && (
              <FormField
                labelWidth={8}
                inputWidth={10}
                value={end}
                onChange={this.onEndChange}
                placeholder="+"
                label="End"
              />
            )}
          </div>
        )}

        {type !== QueryTypeValue.CLI && command && (
          <div className="gf-form">
            {CommandParameters.count.includes(command) && (
              <FormField
                labelWidth={8}
                inputWidth={10}
                value={count}
                type="number"
                onChange={this.onCountChange}
                label="Count"
                tooltip="Can cause latency and is not recommended to use in Production."
              />
            )}

            {CommandParameters.samples.includes(command) && (
              <FormField
                labelWidth={8}
                inputWidth={10}
                value={samples}
                type="number"
                onChange={this.onSamplesChange}
                label="Samples"
                placeholder="5"
                tooltip="Number of sampled nested values. 0 (all values) is not supported."
              />
            )}
          </div>
        )}

        {type === QueryTypeValue.REDIS && command && CommandParameters.section.includes(command) && (
          <div className="gf-form">
            <InlineFormLabel width={8}>Section</InlineFormLabel>
            <Select options={InfoSections} onChange={this.onInfoSectionChange} value={section} menuPlacement="bottom" />
          </div>
        )}

        {type === QueryTypeValue.TIMESERIES && command && CommandParameters.aggregation.includes(command) && (
          <div className="gf-form">
            <InlineFormLabel width={8}>Aggregation</InlineFormLabel>
            <Select
              className={css`
                margin-right: 5px;
              `}
              options={Aggregations}
              width={30}
              onChange={this.onAggregationChange}
              value={aggregation}
              menuPlacement="bottom"
            />
            {aggregation && (
              <FormField
                labelWidth={8}
                value={bucket}
                type="number"
                onChange={this.onBucketChange}
                label="Time Bucket"
                tooltip="Time bucket for aggregation in milliseconds"
              />
            )}
            {aggregation && bucket && CommandParameters.fill.includes(command) && (
              <Switch
                label="Fill Missing"
                labelClass="width-10"
                tooltip="If checked, the datasource will fill missing intervals."
                checked={fill || false}
                onChange={this.onFillChange}
              />
            )}
          </div>
        )}

        {refId === 'A' && (
          <div className="gf-form">
            <Switch
              label="Streaming"
              labelClass="width-8"
              tooltip="If checked, the datasource will stream data. Only Ref A is supported."
              checked={streaming || false}
              onChange={this.onStreamingChange}
            />
            {streaming && (
              <>
                <FormField
                  labelWidth={8}
                  value={streamingInterval}
                  type="number"
                  onChange={this.onStreamingIntervalChange}
                  label="Interval"
                  tooltip="Streaming interval in milliseconds. Default is 1000ms."
                  placeholder="1000"
                />
                <FormField
                  labelWidth={8}
                  value={streamingCapacity}
                  type="number"
                  onChange={this.onStreamingCapacityChange}
                  label="Capacity"
                  tooltip="Values will be constantly added and will never exceed the given capacity. Default is 1000."
                  placeholder="1000"
                />
                <div className="gf-form">
                  <InlineFormLabel width={8} tooltip="If checked Time series, the last line of data will be applied.">
                    Data type
                  </InlineFormLabel>
                  <RadioButtonGroup
                    options={StreamingDataTypes}
                    value={streamingDataType || StreamingDataType.TimeSeries}
                    onChange={this.onStreamingDataTypeChange}
                  />
                </div>
              </>
            )}
          </div>
        )}

        <Button onClick={onRunQuery}>Run</Button>
      </div>
    );
  }
}
