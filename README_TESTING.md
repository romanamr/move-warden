# Testing

Este proyecto usa principalmente pruebas unitarias sobre los paquetes internos.

## Qué se prueba hoy

- `internal/config/config_test.go`
  - reglas de transformación (path, extensión, regex),
  - reglas de filtrado (regex, extensión),
  - parsing/unmarshal de reglas y configuraciones,
  - aplicación de inserciones y secuencia de transformaciones.
- `internal/engine/engine_test.go`
  - mapeo de variables de ruta,
  - creación de directorios destino,
  - ejecución de movimientos de archivo/directorio,
  - recorridos recursivos y propagación de errores,
  - colección de planes de movimiento para escenarios de prueba.

## Comandos de test

### Go directo

```bash
go test ./...
go test -race ./...
go test -coverprofile=cover.out ./...
go tool cover -func=cover.out
```

### Con Makefile

```bash
make test
make test-race
make test-cover
```

## Validación local recomendada

Antes de abrir PR o mergear:

```bash
make check
```

`make check` ejecuta:

1. `make fmt`
2. `make vet`
3. `make test`

## Pruebas manuales del CLI

Además de tests unitarios, se recomienda validar flujo mínimo:

```bash
make run-example-rules
make run-dry RULES_FILE=example_rules.json
```

Esto verifica el parseo básico de flags y la generación de reglas de ejemplo.

## E2E por CLI

Se incluye un runner E2E en Python para simular ejecuciones reales por CLI con 3 escenarios:

1. archivo unico `.txt -> .md`,
2. carpeta con 2 archivos `.txt -> .md`,
3. reorganizacion recursiva por tipo (`JPG/TIFF/PDF`).

Comando:

```bash
python3 scripts/test/e2e_cli_runner.py
make e2e
```

Este script:

- regenera `e2e_generated/` con fixtures de prueba,
- regenera `e2e_test_config/` con reglas JSON por escenario,
- ejecuta por caso `dry-run` y, solo si pasa, `real-run`,
- escribe `e2e_report.md` (detalle humano) y `e2e_report.json` (resumen simple),
- si algo falla sale con `exit code 1`,
- si todo pasa limpia `e2e_generated/` y `e2e_test_config/` automaticamente.

Interpretacion rapida:

- `PASS`: dry-run y real-run del caso pasaron.
- `FAIL_DRY`: fallo en dry-run y se omite real-run de ese caso.
- `FAIL_REAL`: dry-run paso pero real-run o validacion final fallo.
