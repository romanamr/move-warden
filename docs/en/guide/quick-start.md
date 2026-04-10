# Quick Start

## Requirements

- Go `1.25+`

## 1) Generate a sample rules file

```bash
go run . --create-example-rules
```

This creates `example_rules.json` at the project root.

## 2) Run with your rules file

```bash
go run . --rules rules.json
```

## 3) Simulate without moving files

```bash
go run . --dry-run --rules rules.json
```

## Makefile alternative

```bash
make run
make run-dry
make run-example-rules
```

## Next step

- Read [rules.json format](../reference/rules-format)
