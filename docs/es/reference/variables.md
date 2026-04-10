# Variables y placeholders

Puedes usar placeholders entre llaves en rutas de destino y reglas.

Ejemplo:

```json
"to": "destino/{filename}.{ext}"
```

## Variables disponibles

- `{filename}`: nombre sin extensión.
- `{ext}`: extensión sin el punto.
- `{parent_dir}`: carpeta padre inmediata.
- `{fragment_0}`, `{fragment_1}`, ...: fragmentos de ruta.
- `{fragment_init}`: primer fragmento.
- `{fragment_last}`: último fragmento.

## Windows y rutas UNC

También se soportan rutas Windows (`C:\\...`) y rutas de red (`\\\\Servidor\\Share\\...`).

Los fragmentos se calculan usando separadores `/` y `\\`.
