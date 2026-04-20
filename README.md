# vaultdiff

> CLI tool to diff and audit changes between HashiCorp Vault secret versions across environments.

---

## Installation

```bash
go install github.com/yourusername/vaultdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/vaultdiff.git
cd vaultdiff
go build -o vaultdiff .
```

---

## Usage

Compare two versions of a secret within the same Vault path:

```bash
vaultdiff --path secret/myapp/config --v1 3 --v2 5
```

Diff secrets across environments:

```bash
vaultdiff --src secret/staging/myapp --dst secret/production/myapp
```

Export a diff report to a file:

```bash
vaultdiff --path secret/myapp/config --v1 2 --v2 4 --output report.json
```

### Flags

| Flag | Description |
|------|-------------|
| `--path` | Vault secret path |
| `--v1` | First version to compare |
| `--v2` | Second version to compare |
| `--src` | Source environment path |
| `--dst` | Destination environment path |
| `--output` | Write diff output to a file |
| `--format` | Output format: `text` (default), `json`, or `yaml` |
| `--token` | Vault token (or set `VAULT_TOKEN`) |
| `--addr` | Vault address (or set `VAULT_ADDR`) |
| `--mask` | Mask secret values in output (shows `***` instead of plaintext) |

---

## Requirements

- Go 1.21+
- HashiCorp Vault with KV v2 secrets engine

---

## License

MIT © 2024 [Your Name](https://github.com/yourusername)
