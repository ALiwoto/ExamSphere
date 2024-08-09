package dbScripts

import (
	// fs
	"embed"
)

var (
	//go:embed migration1.sql
	Migration1 embed.FS
)

var (
	Migration1Str = func() string {
		fileContent, _ := Migration1.ReadFile("migration1.sql")
		return string(fileContent)
	}()
)
