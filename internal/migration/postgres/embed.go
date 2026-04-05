// Package postgres provides embedded SQL migration files.
package postgres

import "embed"

// FS contains migration files embedded at compile time.
//
//go:embed *.sql
var FS embed.FS
