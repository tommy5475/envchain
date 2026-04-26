# envchain

A tool for chaining and validating environment variable sets across deployment stages.

---

## Installation

```bash
go install github.com/yourorg/envchain@latest
```

Or build from source:

```bash
git clone https://github.com/yourorg/envchain.git && cd envchain && go build ./...
```

---

## Usage

Define your environment variable sets in a `.envchain.yaml` file:

```yaml
stages:
  dev:
    required:
      - DATABASE_URL
      - API_KEY
  production:
    extends: dev
    required:
      - SECRET_TOKEN
      - SENTRY_DSN
```

Then validate your environment before deploying:

```bash
# Validate all variables for a specific stage
envchain validate --stage production

# Chain and export variables across stages
envchain export --stage production --output .env
```

If any required variables are missing, `envchain` will report them clearly:

```
✗ Missing variables for stage "production":
  - SECRET_TOKEN
  - SENTRY_DSN
```

---

## Why envchain?

Managing environment variables across `dev`, `staging`, and `production` is error-prone. `envchain` lets you define inheritance between stages and catch missing variables before they cause runtime failures.

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE)