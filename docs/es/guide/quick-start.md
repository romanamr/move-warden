# Inicio rápido

## Requisitos

- Go `1.25+`

## 1) Generar un ejemplo de reglas

```bash
go run . --create-example-rules
```

Genera `example_rules.json` en la raíz del proyecto.

## 2) Ejecutar con tu archivo de reglas

```bash
go run . --rules rules.json
```

## 3) Simular sin mover archivos

```bash
go run . --dry-run --rules rules.json
```

## Bloque corto: filtro `contains`

```json
"filter_rules": [
	{ "type": "contains", "text": ["/docs/"] }
]
```

Permite solo archivos cuya ruta contiene `/docs/`.

## Alternativa con Makefile

```bash
make run
make run-dry
make run-example-rules
```

## Siguiente paso

- Revisa [Formato de rules.json](../reference/rules-format)
