import { css } from 'emotion';
import React, { ChangeEvent, PureComponent } from 'react';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { Button, InlineFormLabel, LegacyForms, Select, TextArea } from '@grafana/ui';
import { DataSource } from '../DataSource';
import { Aggregations, CommandParameters, Commands, InfoSections, QueryType, QueryTypeValue } from '../redis';
import { RedisDataSourceOptions, RedisQuery } from '../types';

/**
 * Form Field
 */
const { FormField } = LegacyForms;

/**
 * Editor Property
 */
type Props = QueryEditorProps<DataSource, RedisQuery, RedisDataSourceOptions>;

/**
 * Query Editor
 */
export class QueryEditor extends PureComponent<Props> {
  /**
   * Key change
   *
   * @param {ChangeEvent<HTMLInputElement>} event Event
   */
  onKeyChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, key: event.target.value });
  };

  /**
   * Query change
   *
   * @param {ChangeEvent<HTMLInputElement>} event Event
   */
  onQueryChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, query: event.target.value });
  };

  /**
   * Filter change
   *
   * @param {ChangeEvent<HTMLInputElement>} event Event
   */
  onFilterChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, filter: event.target.value });
  };

  /**
   * Field change
   *
   * @param {ChangeEvent<HTMLInputElement>} event Event
   */
  onFieldChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, field: event.target.value });
  };

  /**
   * Legend change
   *
   * @param {ChangeEvent<HTMLInputElement>} event Event
   */
  onLegendChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, legend: event.target.value });
  };

  /**
   * Value change
   *
   * @param {ChangeEvent<HTMLInputElement>} event Event
   */
  onValueChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, value: event.target.value });
  };

  /**
   * Command change
   *
   * @param val Value
   */
  onCommandChange = (val: SelectableValue<string>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, command: val.value });
  };

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
   *
   * @param val Value
   */
  onAggregationTextChange = (val: SelectableValue<string>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, aggregation: val.value });
  };

  /**
   * Info section change
   *
   * @param val Value
   */
  onInfoSectionTextChange = (val: SelectableValue<string>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, section: val.value });
  };

  /**
   * Bucket change
   *
   * @param {ChangeEvent<HTMLInputElement>} event Event
   */
  onBucketTextChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, bucket: event.target.value });
  };

  /**
   * Render Editor
   */
  render() {
    const { key, aggregation, bucket, legend, command, field, filter, value, query, type, section } = this.props.query;
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
            {CommandParameters.key.includes(command) && (
              <FormField
                labelWidth={8}
                inputWidth={30}
                value={key}
                onChange={this.onKeyChange}
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
              <FormField
                labelWidth={8}
                inputWidth={30}
                value={field}
                onChange={this.onFieldChange}
                label="Field"
                tooltip="Field"
              />
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

            {CommandParameters.legendLabel.includes(command) && (
              <FormField
                labelWidth={8}
                inputWidth={10}
                value={legend}
                onChange={this.onLegendChange}
                label="Legend Label"
                tooltip="Legend Label"
              />
            )}

            {CommandParameters.valueLabel.includes(command) && (
              <FormField
                labelWidth={8}
                inputWidth={10}
                value={value}
                onChange={this.onValueChange}
                label="Value Label"
                tooltip="Value Label"
              />
            )}
          </div>
        )}

        {type === QueryTypeValue.COMMAND && command && CommandParameters.section.includes(command) && (
          <div className="gf-form">
            <InlineFormLabel width={8}>Section</InlineFormLabel>
            <Select
              options={InfoSections}
              onChange={this.onInfoSectionTextChange}
              value={section}
              menuPlacement="bottom"
            />
          </div>
        )}

        {type === QueryTypeValue.TIMESERIES && command && CommandParameters.aggregation.includes(command) && (
          <div className="gf-form">
            <InlineFormLabel width={8}>Aggregation</InlineFormLabel>
            <Select
              options={Aggregations}
              width={30}
              onChange={this.onAggregationTextChange}
              value={aggregation}
              menuPlacement="bottom"
            />
            {aggregation && (
              <FormField
                labelWidth={8}
                value={bucket}
                type="number"
                onChange={this.onBucketTextChange}
                label="Bucket"
                tooltip="Time bucket for aggregation in milliseconds"
              />
            )}
          </div>
        )}

        <Button onClick={onRunQuery}>Run</Button>
      </div>
    );
  }
}
