// +build walletrpc

package main

import (
	"context"
	"encoding/hex"

	"github.com/lightningnetwork/lnd/lnrpc/walletrpc"
	"github.com/urfave/cli"
)

// walletrpcCommands will return nil for non-walletrpc builds.
func walletrpcCommands() []cli.Command {
	return []cli.Command{
		keyForAddressCommand,
	}
}

func getWalletKitClient(ctx *cli.Context) (walletrpc.WalletKitClient, func()) {
	conn := getClientConn(ctx, false)

	cleanUp := func() {
		conn.Close()
	}

	return walletrpc.NewWalletKitClient(conn), cleanUp
}

var keyForAddressCommand = cli.Command{
	Name:      "keyforaddress",
	Category:  "On-chain",
	Usage:     "Get the public key for a given in-wallet address",
	ArgsUsage: "bitcoin-address",
	Description: `
	On receipt of a (segwit) address that is owned by the wallet,
	this call returns the corresponding pubkey in serialized compressed
	form. If the bitcoin address is not recognized as owned by the
	wallet, an error is returned.
	`,
	Flags:  []cli.Flag{},
	Action: actionDecorator(keyForAddress),
}

// keyForAddress is the function called by the keyforaddress command
// of the walletkit rpc client.
func keyForAddress(ctx *cli.Context) error {
	var (
		address string
		err     error
	)
	args := ctx.Args()
	address = args.First()

	ctxb := context.Background()
	client, cleanUp := getWalletKitClient(ctx)
	defer cleanUp()

	req := &walletrpc.KeyForAddressRequest{
		AddrIn: address,
	}

	resp, err := client.KeyForAddress(ctxb, req)
	if err != nil {
		return err
	}

	// The response is a signrpc.KeyDescriptor which holds
	// raw bytes for the key, convert to hex for reading:
	printJSON(hex.EncodeToString(resp.RawKeyBytes))

	return nil
}
