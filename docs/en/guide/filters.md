# Filters

All `filter_rules` must pass for a file to be allowed.

## `extension`

Allows only extensions listed in `extensions`.

```json
{
  "type": "extension",
  "extensions": [".jpg", ".jpeg", ".png", ".heic"]
}
```

## `regex`

Allows only paths matching the regex pattern.

```json
{
  "type": "regex",
  "pattern": ".*"
}
```

## `contains`

Allows only paths that contain at least one value from `text`.

```json
{
  "type": "contains",
  "text": ["/docs/", "manual"]
}
```

## Note

If one rule fails, the file is excluded from movement.
