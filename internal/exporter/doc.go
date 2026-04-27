// Package exporter provides output formatting for validated environment
// variable sets produced by envchain.
//
// Supported formats:
//
//   - shell   — emits `export KEY=VALUE` lines suitable for sourcing in bash/zsh
//   - dotenv  — emits `KEY=VALUE` lines compatible with .env file loaders
//   - json    — emits a JSON object mapping variable names to their values
//
// Example usage:
//
//	e := exporter.New(os.Stdout, exporter.FormatShell)
//	if err := e.Export(map[string]string{"APP_ENV": "production"}); err != nil {
//		log.Fatal(err)
//	}
package exporter
