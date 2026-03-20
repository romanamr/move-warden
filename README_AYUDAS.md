# Ayudas y scripts

Este documento describe utilidades para trabajar en el proyecto de forma consistente.

## Makefile (principal)

Se añadió `Makefile` en la raíz para centralizar comandos de desarrollo.

### Targets disponibles

- `make help`: lista todos los comandos.
- `make run`: ejecuta la app usando `rules.json`.
- `make run-dry`: ejecuta la app con `--dry-run`.
- `make run-example-rules`: genera `example_rules.json`.
- `make build`: compila el binario en `bin/movewarden`.
- `make test`: ejecuta tests unitarios (`go test ./...`).
- `make test-race`: ejecuta tests con detector de race.
- `make test-cover`: genera `cover.out` y resumen de cobertura.
- `make fmt`: formatea código (`go fmt ./...`).
- `make vet`: ejecuta chequeos de `go vet`.
- `make lint`: ejecuta `golangci-lint` (si está instalado).
- `make check`: pipeline local recomendado (`fmt + vet + test`).
- `make clean`: limpia artefactos (`bin/`, `cover.out`).

### Variables sobrescribibles

- `GO`: binario de Go (por defecto `go`).
- `APP_NAME`: nombre del binario (`movewarden`).
- `RULES_FILE`: archivo de reglas para `run`/`run-dry` (`rules.json`).

Ejemplo:

```bash
make run RULES_FILE=example_rules.json
```

## Scripts existentes

Carpeta `scripts/`:

- `scripts/crear-labels.ps1`: script PowerShell orientado a tareas de soporte en repositorio.

## Flujo recomendado diario

```bash
make fmt
make vet
make test
make run-dry RULES_FILE=rules.json
```

Si todo está correcto, se puede usar `make check` para validación rápida previa a commit.
