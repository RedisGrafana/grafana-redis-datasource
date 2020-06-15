import React, { ChangeEvent, PureComponent } from 'react';
import { LegacyForms, InlineFormLabel } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from './DataSource';
import { RedisTimeSeriesDataSourceOptions, RedisTimeSeriesQuery } from './types';
import { Select } from '@grafana/ui';

const { FormField } = LegacyForms;

type Props = QueryEditorProps<DataSource, RedisTimeSeriesQuery, RedisTimeSeriesDataSourceOptions>;
const aggreations: Array<SelectableValue<string>> = [
  { label: 'None', description: 'no aggregation', value: '' },
  { label: 'Max', description: 'max', value: 'max' },
  { label: 'Min', description: 'min', value: 'min' },
  { label: 'Rate', description: 'rate', value: 'rate' },
  { label: 'Count', description: 'count number of samples', value: 'count' },
  { label: 'Range', description: 'Diff between max and min in the bucket', value: 'range' },
];

const cmdTypes: Array<SelectableValue<string>> = [
  { label: 'TS.RANGE', description: 'range query', value: 'tsrange' },
  { label: 'HGETALL', description: 'hashtable', value: 'hgetall' },
];
export class QueryEditor extends PureComponent<Props> {
  cmdtype?: string;
  onQueryKeyChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, keyname: event.target.value });
    // executes the query
    onRunQuery();
  };

  onQueryLegendChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, legend: event.target.value });
    // executes the query
    onRunQuery();
  };

  onCmdTypeChange = (val: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, cmd: val.value });
    // executes the query
    this.cmdtype = val.value;
    this.forceUpdate();
    onRunQuery();
  };

  onAggregationTextChange = (val: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, aggregation: val.value });
    // executes the query
    onRunQuery();
  };

  onBucketTextChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, bucket: event.target.value });
    // executes the query
    onRunQuery();
  };

  render() {
    const { keyname, aggregation, bucket, legend, cmd } = this.props.query;
    this.cmdtype = cmd;
    let selected_agg = aggreations[0];
    aggreations.forEach(option => {
      if (option.value === aggregation) {
        selected_agg = option;
      }
    });
    return (
      <div>
        <div className="gf-form">
          <FormField labelWidth={8} value={keyname} onChange={this.onQueryKeyChange} label="Key" tooltip="keyname" />
          <Select options={cmdTypes} menuPlacement="bottom" value={this.cmdtype} onChange={this.onCmdTypeChange} />
        </div>
        {this.cmdtype === 'tsrange' && (
          <div className="gf-form">
            <InlineFormLabel width={8}>Aggregation</InlineFormLabel>
            <Select
              options={aggreations}
              onChange={this.onAggregationTextChange}
              value={selected_agg}
              menuPlacement="bottom"
            />
            <FormField
              labelWidth={8}
              value={bucket}
              onChange={this.onBucketTextChange}
              label="Bucket"
              tooltip="keyname"
            />
            <FormField
              labelWidth={8}
              value={legend}
              onChange={this.onQueryLegendChange}
              label="Legend"
              tooltip="Legend override"
            />
          </div>
        )}
      </div>
    );
  }
}
