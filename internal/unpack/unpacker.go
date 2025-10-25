package unpack

import "github.com/evaneliasyoung/phex/internal/phaser"

type Unpacker struct {
	phaser.Atlas
	PackName  string
	InputDir  string
	OutputDir string
	Workers   int
}
