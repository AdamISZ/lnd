// +build !signrpc

package main

import "github.com/urfave/cli"

// signrpcCommands will return nil for non-signrpc builds.
func signrpcCommands() []cli.Command {
	return nil
}
