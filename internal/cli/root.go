package cli

import (
	"os"

	"golang.org/x/term"
)

var IsTTY = term.IsTerminal(int(os.Stdout.Fd()))
