.PHONY: help run run-dry run-example-rules build test test-race test-cover fmt vet lint check e2e clean

GO ?= go
APP_NAME ?= movewarden
RULES_FILE ?= rules.json

help:
	@echo "Targets disponibles:"
	@echo "  make run                 - Ejecuta la app con rules.json"
	@echo "  make run-dry             - Ejecuta en modo dry-run"
	@echo "  make run-example-rules   - Genera example_rules.json"
	@echo "  make build               - Compila el binario en ./bin/movewarden"
	@echo "  make test                - Ejecuta tests unitarios"
	@echo "  make test-race           - Ejecuta tests con detector de race"
	@echo "  make test-cover          - Genera reporte de cobertura en cover.out"
	@echo "  make fmt                 - Formatea el codigo con gofmt"
	@echo "  make vet                 - Ejecuta analisis estatico con go vet"
	@echo "  make lint                - Ejecuta golangci-lint (si esta instalado)"
	@echo "  make check               - Ejecuta fmt + vet + test"
	@echo "  make e2e                 - Ejecuta test unitarios y luego E2E CLI"
	@echo "  make clean               - Limpia artefactos de build/testing"

run:
	$(GO) run . --rules $(RULES_FILE)

run-dry:
	$(GO) run . --dry-run --rules $(RULES_FILE)

run-example-rules:
	$(GO) run . --create-example-rules

build:
	@mkdir -p bin
	$(GO) build -o bin/$(APP_NAME) .

test:
	$(GO) test ./...

test-race:
	$(GO) test -race ./...

test-cover:
	$(GO) test -coverprofile=cover.out ./...
	$(GO) tool cover -func=cover.out

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

lint:
	golangci-lint run ./...

check: fmt vet test

e2e: test
	python3 scripts/test/e2e_cli_runner.py

clean:
	rm -rf bin cover.out
