/**
 * Supported Commands
 */
export enum RediSearch {
  INFO = 'ft.info',
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
];
