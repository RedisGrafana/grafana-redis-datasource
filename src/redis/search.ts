import { SelectableValue } from '@grafana/data';

/**
 * Supported Commands
 */
export enum RediSearch {
  INFO = 'ft.info',
  SEARCH = 'ft.search',
}

/**
 * Sort Directions
 */
export enum SortDirectionValue {
  NONE = 'None',
  ASC = 'ASC',
  DESC = 'DESC',
}
/**
 * Commands List
 */
export const RediSearchCommands = [
  {
    label: RediSearch.INFO.toUpperCase(),
    description: 'Returns information and statistics on the index',
    value: RediSearch.INFO,
  },
  {
    label: RediSearch.SEARCH.toUpperCase(),
    description: 'Search the index with a textual query, returning either documents or just ids',
    value: RediSearch.SEARCH,
  },
];

export const SortDirection: Array<SelectableValue<SortDirectionValue>> = [
  {
    label: 'None',
    description: "Don't sort anything.",
    value: SortDirectionValue.NONE,
  },
  {
    label: 'Ascending',
    description: 'Sort the field in Ascending order.',
    value: SortDirectionValue.ASC,
  },
  {
    label: 'Descending',
    description: 'Sort the values in descending order.',
    value: SortDirectionValue.DESC,
  },
];
