# Transformations

`transformation_rules` are applied in order, one after another.

## `path_change`

Replaces a path fragment with another one.

```json
{
  "type": "path_change",
  "from": "workspace/inbox/photos",
  "to": "workspace/library/{parent_dir}"
}
```

## `extension`

Rewrites file extensions using `from`/`to` pairs.

```json
{
  "type": "extension",
  "extensions": [
    { "from": ".jpeg", "to": ".jpg" },
    { "from": ".heic", "to": ".jpg" }
  ]
}
```

## `regex`

Applies regex replacement over the full path.

```json
{
  "type": "regex",
  "pattern": " ",
  "replacement": "_"
}
```

## Note

The output of one transformation is the input of the next.
