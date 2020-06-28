import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { Button, InlineFormLabel, LegacyForms, Select, TextArea } from '@grafana/ui';
import { Aggregations, CommandParameters, Commands, QueryType, QueryTypeValue, InfoSections } from 'command';
import React, { ChangeEvent, PureComponent } from 'react';
import { DataSource } from './DataSource';
import { RedisDataSourceOptions, RedisQuery } from './types';

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
   * @param event Event
   */
  onKeyChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, key: event.target.value });
  };

  /**
   * Query change
   *
   * @param event Event
   */
  onQueryChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, query: event.target.value });
  };

  /**
   * Execute the Query
   */
  executeQuery = () => {
    const { onRunQuery } = this.props;
    onRunQuery();
  };

  /**
   * Filter change
   *
   * @param event Event
   */
  onFilterChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, filter: event.target.value });
  };

  /**
   * Field change
   *
   * @param event Event
   */
  onFieldChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, field: event.target.value });
  };

  /**
   * Legend change
   *
   * @param event Event
   */
  onLegendChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, legend: event.target.value });
  };

  /**
   * Value change
   *
   * @param event Event
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
    onChange({ ...query, type: val.value, query: val.value === QueryTypeValue.COMMAND ? '' : query.query });
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
   * @param val Value
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

    /**
     * Return
     */
    return (
      <div>
        <div className="gf-form">
          <InlineFormLabel width={8}>Type</InlineFormLabel>
          <Select options={QueryType} menuPlacement="bottom" value={type} onChange={this.onTypeChange} />
        </div>

        {type === QueryTypeValue.CLI && (
          <div className="gf-form">
            <InlineFormLabel width={8}>Command</InlineFormLabel>
            <TextArea value={query} onChange={this.onQueryChange} />
          </div>
        )}

        {type === QueryTypeValue.COMMAND && (
          <div className="gf-form">
            <InlineFormLabel width={8}>Command</InlineFormLabel>
            <Select options={Commands} menuPlacement="bottom" value={command} onChange={this.onCommandChange} />
          </div>
        )}

        {type === QueryTypeValue.COMMAND && command && (
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

        {type === QueryTypeValue.COMMAND && command && CommandParameters.aggregation.includes(command) && (
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

        <Button onClick={this.executeQuery}>Run</Button>
      </div>
    );
  }
}
