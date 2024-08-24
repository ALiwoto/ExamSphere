package dbScripts

import (
	_ "embed"
)

var (
	//go:embed migration1.sql
	Migration1Str string

	//go:embed migration2.sql
	Migration2Str string

	//go:embed migration3.sql
	Migration3Str string

	//go:embed migration4.sql
	Migration4Str string
)
