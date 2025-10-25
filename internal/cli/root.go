package cli

import (
	"os"

	"golang.org/x/term"
)

var Version string = "dev"
var IsTTY = term.IsTerminal(int(os.Stdout.Fd()))
