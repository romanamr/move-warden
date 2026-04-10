# Filtros

Las `filter_rules` deben cumplirse todas para permitir un archivo.

## `extension`

Permite solo extensiones incluidas en la lista.

```json
{
  "type": "extension",
  "extensions": [".jpg", ".jpeg", ".png", ".heic"]
}
```

## `regex`

Permite solo rutas que cumplen el patrón.

```json
{
  "type": "regex",
  "pattern": ".*"
}
```

## Nota

Si una regla falla, el archivo queda fuera del movimiento.
