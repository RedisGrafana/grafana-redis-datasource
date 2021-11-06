/**
 * Supported Commands
 */
export enum RedisJson {
  OBJKEYS = 'json.objkeys',
  OBJLEN = 'json.objlen',
  GET = 'json.get',
  TYPE = 'json.type',
  ARRLEN = 'json.arrlen',
}

/**
 * Commands List
 */
export const RedisJsonCommands = [
  {
    label: RedisJson.ARRLEN.toUpperCase(),
    description: 'Report the length of the JSON Array at path in key',
    value: RedisJson.ARRLEN,
  },
  {
    label: RedisJson.GET.toUpperCase(),
    description: 'Return the value at path',
    value: RedisJson.GET,
  },
  {
    label: RedisJson.OBJKEYS.toUpperCase(),
    description: "Return the keys in the object that's referenced by path",
    value: RedisJson.OBJKEYS,
  },
  {
    label: RedisJson.OBJLEN.toUpperCase(),
    description: 'Report the number of keys in the JSON Object at path in key',
    value: RedisJson.OBJLEN,
  },
  {
    label: RedisJson.TYPE.toUpperCase(),
    description: 'Report the type of JSON value at path',
    value: RedisJson.TYPE,
  },
];
