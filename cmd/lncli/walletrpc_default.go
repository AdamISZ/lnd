// +build !walletrpc

package main

import "github.com/urfave/cli"

// walletrpcCommands will return nil for non-walletrpc builds.
func walletrpcCommands() []cli.Command {
	return nil
}
