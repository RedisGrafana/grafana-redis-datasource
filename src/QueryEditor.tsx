import { Aggregations, CommandParameters, Commands } from 'command';
import React, { ChangeEvent, PureComponent } from 'react';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { InlineFormLabel, LegacyForms, Select } from '@grafana/ui';
import { DataSource } from './DataSource';
import { RedisDataSourceOptions, RedisQuery } from './types';

const { FormField } = LegacyForms;

type Props = QueryEditorProps<DataSource, RedisQuery, RedisDataSourceOptions>;

/**
 * Query Editor
 */
export class QueryEditor extends PureComponent<Props> {
  /**
   * Key change
   *
   * @param event Event
   */
  onKeyChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, key: event.target.value });
    onRunQuery();
  };

  /**
   * Filter change
   *
   * @param event Event
   */
  onFilterChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, filter: event.target.value });
    onRunQuery();
  };

  /**
   * Field change
   *
   * @param event Event
   */
  onFieldChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, field: event.target.value });
    onRunQuery();
  };

  /**
   * Legend change
   *
   * @param event Event
   */
  onLegendChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, legend: event.target.value });
    onRunQuery();
  };

  /**
   * Value change
   *
   * @param event Event
   */
  onValueChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, value: event.target.value });
    onRunQuery();
  };

  /**
   * Command change
   *
   * @param val Value
   */
  onCommandChange = (val: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, command: val.value });

    /**
     * Update form
     */
    this.forceUpdate();

    /**
     * Executes the query
     */
    onRunQuery();
  };

  /**
   * Aggregation change
   *
   * @param val Value
   */
  onAggregationTextChange = (val: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, aggregation: val.value });
    // executes the query
    onRunQuery();
  };

  /**
   * Bucket change
   *
   * @param val Value
   */
  onBucketTextChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, bucket: event.target.value });
    // executes the query
    onRunQuery();
  };

  /**
   * Render
   */
  render() {
    const { key, aggregation, bucket, legend, command, field, filter, value } = this.props.query;

    /**
     * Return
     */
    return (
      <div>
        <div className="gf-form">
          <InlineFormLabel width={8}>Command</InlineFormLabel>
          <Select options={Commands} menuPlacement="bottom" value={command} onChange={this.onCommandChange} />
        </div>

        <div className="gf-form">
          {command && CommandParameters.key.includes(command) && (
            <FormField
              labelWidth={8}
              inputWidth={30}
              value={key}
              onChange={this.onKeyChange}
              label="Key"
              tooltip="Key name"
            />
          )}

          {command && CommandParameters.filter.includes(command) && (
            <FormField
              labelWidth={8}
              inputWidth={30}
              value={filter}
              onChange={this.onFilterChange}
              label="Label Filter"
              tooltip="Label Filter"
            />
          )}

          {command && CommandParameters.field.includes(command) && (
            <FormField
              labelWidth={8}
              inputWidth={30}
              value={field}
              onChange={this.onFieldChange}
              label="Field"
              tooltip="Field"
            />
          )}

          {command && CommandParameters.legend.includes(command) && (
            <FormField
              labelWidth={8}
              inputWidth={20}
              value={legend}
              onChange={this.onLegendChange}
              label="Legend"
              tooltip="Legend override"
            />
          )}

          {command && CommandParameters.legendLabel.includes(command) && (
            <FormField
              labelWidth={8}
              inputWidth={20}
              value={legend}
              onChange={this.onLegendChange}
              label="Legend Label"
              tooltip="Legend Label"
            />
          )}

          {command && CommandParameters.valueLabel.includes(command) && (
            <FormField
              labelWidth={8}
              inputWidth={20}
              value={value}
              onChange={this.onValueChange}
              label="Value Label"
              tooltip="Value Label"
            />
          )}
        </div>

        {command && CommandParameters.aggregation.includes(command) && (
          <div className="gf-form">
            <InlineFormLabel width={8}>Aggregation</InlineFormLabel>
            <Select
              options={Aggregations}
              width={30}
              onChange={this.onAggregationTextChange}
              value={aggregation}
              menuPlacement="bottom"
            />
            <FormField
              labelWidth={8}
              value={bucket}
              type="number"
              onChange={this.onBucketTextChange}
              label="Bucket"
              tooltip="Time bucket for aggregation in milliseconds"
            />
          </div>
        )}
      </div>
    );
  }
}
