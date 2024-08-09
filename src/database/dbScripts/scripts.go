package dbScripts

import (
	// fs
	"embed"
)

var (
	//go:embed migration1.sql
	Migration1 embed.FS

	//go:embed migration2.sql
	Migration2 embed.FS

	//go:embed migration3.sql
	Migration3 embed.FS
)

var (
	Migration1Str = func() string {
		fileContent, _ := Migration1.ReadFile("migration1.sql")
		return string(fileContent)
	}()

	Migration2Str = func() string {
		fileContent, _ := Migration2.ReadFile("migration2.sql")
		return string(fileContent)
	}()

	Migration3Str = func() string {
		fileContent, _ := Migration3.ReadFile("migration3.sql")
		return string(fileContent)
	}()
)
