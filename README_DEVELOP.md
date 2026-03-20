# Desarrollo

Guía técnica para entender cómo está planteada la aplicación y cómo extenderla.

## Stack técnico

- Lenguaje: Go.
- Módulo: `movewarden`.
- Enfoque: CLI para procesamiento de rutas/archivos basado en reglas JSON.

## Arquitectura actual

La aplicación está organizada por paquetes internos con responsabilidades claras:

- `main.go`
  - Parsea flags (`--dry-run`, `--rules`, `--create-example-rules`).
  - Gestiona flujo inicial del CLI.
- `internal/config`
  - Define contratos de configuración (`MovementConfiguration`, `MovementRun`).
  - Implementa parsing flexible de reglas con interfaces:
    - `TransformationRule` para transformar rutas.
    - `FilterRule` para permitir/bloquear archivos.
  - Encapsula lógica de aplicación de transformaciones e inserciones.
- `internal/engine`
  - Orquesta ejecución de movimientos (`Run`, `runMovement`, `runFile`, `runDirectory`).
  - Aísla estrategia de movimiento con `MoveFunc` (real, dry-run, collector para pruebas).
  - Maneja utilidades de filesystem (creación de destino, limpieza de directorios vacíos, recorrido recursivo).
- `internal/tui`
  - Reserva para futura interfaz textual (TUI).

## Flujo de procesamiento (conceptual)

1. Se carga una configuración de movimientos.
2. Para cada movimiento:
   - se identifica si `source` es archivo o directorio,
   - se construye un mapa de variables (`filename`, `ext`, `fragment_*`, etc.),
   - se aplican transformaciones en orden,
   - se aplican inserciones `{key}`,
   - se ejecuta el movimiento mediante una función inyectable.

Este diseño permite probar fácilmente el comportamiento sin tocar el filesystem real.

## Decisiones de diseño

- **Reglas polimórficas**: interfaces para transformación/filtro facilitan agregar nuevos tipos sin romper el esquema.
- **Parsing desacoplado por tipo**: uso de mapas de unmarshallers por `type` para soportar extensión incremental.
- **Inyección de comportamiento de movimiento** (`MoveFunc`): separa "calcular destino" de "efectuar IO", útil para dry-run y testing.
- **Paquetes internos**: evita exposición accidental de APIs no estables.

## Cómo extender el sistema

### Agregar nueva transformación

1. Crear struct que implemente `Apply(source string) string`.
2. Añadir un unmarshaller en `internal/config/config.go`.
3. Registrar la nueva entrada en `transformationRuleUnmarshallers`.
4. Añadir tests unitarios (parsing + comportamiento).

### Agregar nuevo filtro

1. Crear struct que implemente `Allowed(fullSourcePath string) bool`.
2. Registrar unmarshaller en `filterRuleUnmarshallers`.
3. Añadir tests unitarios.

## Convenciones de trabajo recomendadas

- Mantener imports al inicio del archivo.
- Añadir pruebas por cada nueva regla o comportamiento no trivial.
- Ejecutar validación local antes de commit:

```bash
make check
```
