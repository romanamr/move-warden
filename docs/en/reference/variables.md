# Variables and placeholders

You can use placeholders in destination paths and rules.

Example:

```json
"to": "target/{filename}.{ext}"
```

## Available variables

- `{filename}`: file name without extension.
- `{ext}`: extension without dot.
- `{parent_dir}`: immediate parent folder.
- `{fragment_0}`, `{fragment_1}`, ...: path fragments.
- `{fragment_init}`: first fragment.
- `{fragment_last}`: last fragment.

## Windows and UNC paths

Windows paths (`C:\\...`) and UNC paths (`\\\\Server\\Share\\...`) are supported.

Fragments are computed using both `/` and `\\` separators.
