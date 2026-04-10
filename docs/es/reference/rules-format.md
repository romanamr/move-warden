# Formato de `rules.json`

Estructura raíz:

```json
{
  "dry_run": true,
  "delete_empty_directories": false,
  "movements": []
}
```

## Campos raíz

- `movements` (requerido): lista de movimientos.
- `dry_run` (opcional): simula ejecución.
- `delete_empty_directories` (opcional): borra carpetas vacías al final.

## Cada elemento de `movements`

- `source` (requerido): ruta de origen.
- `recursive` (opcional): procesa subcarpetas.
- `change_key_map` (opcional): reemplazos manuales `{key}` -> valor.
- `transformation_rules` (opcional): cambios de ruta/nombre.
- `filter_rules` (opcional): condiciones de inclusión.

## Tipos soportados

- Transformaciones: `path_change`, `extension`, `regex`
- Filtros: `extension`, `regex`, `contains`

Consulta [Ejemplos](./examples) para plantillas completas.
