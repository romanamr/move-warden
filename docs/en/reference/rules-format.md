# `rules.json` format

Root structure:

```json
{
  "dry_run": true,
  "delete_empty_directories": false,
  "movements": []
}
```

## Root fields

- `movements` (required): list of movements.
- `dry_run` (optional): simulation mode.
- `delete_empty_directories` (optional): remove empty folders at the end.

## Fields inside each `movements` item

- `source` (required): source path.
- `recursive` (optional): process subdirectories.
- `change_key_map` (optional): manual `{key}` -> value replacements.
- `transformation_rules` (optional): path/name rewrite rules.
- `filter_rules` (optional): inclusion constraints.

## Supported types

- Transformations: `path_change`, `extension`, `regex`
- Filters: `extension`, `regex`

See [Examples](./examples) for complete templates.
