// Command cauldron boots a project and the third-party APIs it depends on.
package main

import (
	"os"

	"github.com/CauldronUp/cauldron/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
