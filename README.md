# Move Warden

Herramienta CLI en Go para definir y ejecutar reglas de movimiento/renombrado de archivos y rutas mediante configuración JSON.

## Estado actual

- La entrada principal está en `main.go`.
- Soporta flags de ejecución para:
  - ejecutar con archivo de reglas (`--rules`),
  - simulación (`--dry-run`),
  - generar ejemplo (`--create-example-rules`).
- El motor de movimientos y reglas vive en `internal/engine` y `internal/config`.
- La integración completa del motor desde `main.go` está en evolución; hoy el CLI valida/lee reglas y permite generar `example_rules.json`.

## Requisitos

- Go `1.25+` (según `go.mod`).

## Inicio rápido

1. Generar archivo de ejemplo:

```bash
go run . --create-example-rules
```

Este comando genera `example_rules.json` con una configuración realista de referencia:

- movimiento recursivo de fotos (`workspace/inbox/photos` -> `workspace/library/...`),
- normalización de extensiones (`.jpeg/.heic` a `.jpg`),
- ejemplo de renombrado por regex (espacios a `_`),
- movimiento de documentos (`workspace/inbox/docs` -> `workspace/library/docs`).

2. Ejecutar usando reglas:

```bash
go run . --rules rules.json
```

3. Ejecutar en modo simulación:

```bash
go run . --dry-run --rules rules.json
```

## Estructura del archivo de reglas

Formato raíz:

```json
{
  "dry_run": true,
  "delete_empty_directories": false,
  "movements": []
}
```

Campos requeridos y opcionales:

- `movements` (**requerido**): lista de movimientos a procesar.
- `dry_run` (opcional): si no se informa, se considera `false`.
- `delete_empty_directories` (opcional): si no se informa, se considera `false`.

En cada elemento de `movements`:

- `source` (**requerido**): ruta origen (archivo o directorio).
- `recursive` (opcional): por defecto `false`.
- `change_key_map` (opcional): reemplazos manuales para placeholders `{key}`.
- `transformation_rules` (opcional): transformaciones de ruta/nombre.
- `filter_rules` (opcional): filtros que deben cumplirse para permitir el archivo.

Tipos de regla soportados:

- transformación:
  - `path_change`
  - `extension`
  - `regex`
- filtro:
  - `extension`
  - `regex`

## Uso con Makefile

Este repositorio incluye un `Makefile` para simplificar tareas comunes:

```bash
make help
make run
make run-dry
make run-example-rules
make test
make check
```

Más detalle en `README_AYUDAS.md`.

## Estructura del proyecto

- `main.go`: punto de entrada del CLI.
- `internal/config`: tipos de configuración, parsing JSON y reglas de transformación/filtro.
- `internal/engine`: ejecución de movimientos, utilidades de filesystem y planificación de moves.
- `internal/tui`: espacio para interfaz TUI (actualmente mínimo).
- `scripts`: scripts de apoyo del proyecto.

## Documentación complementaria

- `README_AYUDAS.md`: scripts y automatizaciones (incluye `Makefile`).
- `README_DEVELOP.md`: arquitectura, decisiones de diseño y guía de desarrollo.
- `README_TESTING.md`: estrategia y comandos de test.

### Uso de llaves `{}` y variables de mapeo en rutas destino

En los movimientos configurados en el JSON, puedes usar llaves `{}` para indicar variables dinámicas dentro de las rutas de destino, por ejemplo:

```json
"to": "destino/{filename}.{ext}"
```

Las variables que puedes usar entre `{}` dependen del análisis automático de la ruta de origen (ver implementación en `internal/engine/utils.go`). Estas son las variables disponibles por defecto:

- `{filename}`: nombre base del archivo sin extensión.
  - Ejemplo: `/foo/bar/foto.jpeg` → `foto`
- `{ext}`: la extensión del archivo (sin el punto).
  - Ejemplo: `/foo/bar/foto.jpeg` → `jpeg`
- `{parent_dir}`: nombre del directorio padre inmediato de la ruta origen.
  - Ejemplo: `/foo/bar/foto.jpeg` → `bar`
- `{fragment_0}`, `{fragment_1}`, ...: fragmentos individuales de la ruta (separados por `/` o `\`), empezando desde el inicio de la ruta.
  - Ejemplo: `/tmp/origin/file.jpg`
    - `fragment_0`: ""
    - `fragment_1`: "tmp"
    - `fragment_2`: "origin"
    - `fragment_3`: "file.jpg"
- `{fragment_init}`: primer fragmento de la ruta.
- `{fragment_last}`: último fragmento de la ruta (usualmente el archivo objetivo o la última carpeta).

> Estas variables también pueden usarse como reemplazos en reglas de transformación (`transformation_rules`) y como claves en `change_key_map`. Puedes combinarlas para definir rutas destino flexibles.

#### Notas sobre rutas en Windows y rutas de red

El sistema soporta tanto rutas tipo Unix (`/home/usuario/file.txt`) como rutas de Windows (`C:\Users\usuario\file.txt`) y rutas en red de Windows (`\\SERVIDOR\share\carpeta\archivo.pdf`):

- **Separadores**: Los fragmentos de ruta se extraen usando ambos separadores (`/` y `\`), para soportar ambos estilos de sistema operativo.
- **Rutas en red (`UNC`) de Windows**: Si usas rutas tipo `\\Servidor\Share\Carpeta\Archivo.txt`, los fragmentos se tomarán después del doble backslash inicial:
  - Por ejemplo, la ruta `\\Servidor\Share\Docs\Archivo.txt` se desglosará como:
    - `fragment_0`: "Servidor"
    - `fragment_1`: "Share"
    - `fragment_2`: "Docs"
    - `fragment_3`: "Archivo.txt"
    - Por tanto, `{parent_dir}` sobre `Archivo.txt` sería `"Docs"`.
- **Ejemplo típico en Windows**:
  - Ruta: `C:\Users\Maria\Fotos\viaje.jpg`
    - `fragment_0`: "C:"
    - `fragment_1`: "Users"
    - `fragment_2`: "Maria"
    - `fragment_3`: "Fotos"
    - `fragment_4`: "viaje.jpg"
    - `{filename}`: "viaje"
    - `{ext}`: "jpg"
    - `{parent_dir}`: "Fotos"

Esto permite definir reglas y rutas destino independientemente de si ejecutas el programa en Windows o Linux, y utilizar parámetros dinámicos en tus rutas de salida.
