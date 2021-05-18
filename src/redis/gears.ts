/**
 * Supported Commands
 */
export enum RedisGears {
  DUMPREGISTRATIONS = 'rg.dumpregistrations',
  PYSTATS = 'rg.pystats',
  PYDUMPREQS = 'rg.pydumpreqs',
  PYEXECUTE = 'rg.pyexecute',
}

/**
 * Commands List
 */
export const RedisGearsCommands = [
  {
    label: RedisGears.DUMPREGISTRATIONS.toUpperCase(),
    description: 'Outputs the list of function registrations',
    value: RedisGears.DUMPREGISTRATIONS,
  },
  {
    label: RedisGears.PYSTATS.toUpperCase(),
    description: 'Returns memory usage statistics from the Python interpreter',
    value: RedisGears.PYSTATS,
  },
  {
    label: RedisGears.PYDUMPREQS.toUpperCase(),
    description: 'Returns a list of all the python requirements available',
    value: RedisGears.PYDUMPREQS,
  },
];
