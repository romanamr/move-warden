# Transformaciones

Las `transformation_rules` se aplican en orden, una detrás de otra.

## `path_change`

Reemplaza un fragmento de ruta por otro.

```json
{
  "type": "path_change",
  "from": "workspace/inbox/photos",
  "to": "workspace/library/{parent_dir}"
}
```

## `extension`

Cambia extensiones según pares `from`/`to`.

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

Aplica reemplazo regex sobre la ruta completa.

```json
{
  "type": "regex",
  "pattern": " ",
  "replacement": "_"
}
```

## Nota

El resultado de una transformación es la entrada de la siguiente.
