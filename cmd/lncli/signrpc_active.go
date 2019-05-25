// +build signrpc

package main

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/lightningnetwork/lnd/lnrpc/signrpc"
	"github.com/urfave/cli"
)

// signrpcCommands will return nil for non-signrpc builds.
func signrpcCommands() []cli.Command {
	return []cli.Command{
		signOutputsCommand,
	}
}

func getSignClient(ctx *cli.Context) (signrpc.SignerClient, func()) {
	conn := getClientConn(ctx, false)

	cleanUp := func() {
		conn.Close()
	}

	return signrpc.NewSignerClient(conn), cleanUp
}

var signOutputsCommand = cli.Command{
	Name:      "signoutputs",
	Category:  "On-chain",
	Usage:     "Sign a spending of one or more utxos in wallet",
	ArgsUsage: "raw-tx-hex pubkey-hex --input-index=n",
	Description: `
	Signs the given bitcoin transaction with the private key corresponding
	to the pubkey provided as second argument, at transaction index
	specified with the flag --input-index (or index 0 if not specified).
	`,
	Flags: []cli.Flag{
		cli.Int64Flag{
			Name:  "key_family",
			Usage: "the HD key family of the key to sign with",
		},
		cli.Int64Flag{
			Name:  "key_index",
			Usage: "the HD key index of the key to sign with",
		},
		cli.Int64Flag{
			Name:  "input_index",
			Usage: "the input index at which to sign the tx, default=0",
		},
		cli.StringFlag{
			Name: "tweak_bytes",
			Usage: "For signing the transaction with a tweaked " +
				"private key; specify a 64 character hex string " +
				"to be added to the private key before signing.",
		},
	},
	Action: actionDecorator(signOutputs),
}

func signOutputs(ctx *cli.Context) error {
	var (
		hextxbytes     string
		hexpubkeybytes string
		//hextweakbytes  string
		txinputindex int64
		//keyfamily      int64
		//keyindex       int64
		err error
	)
	args := ctx.Args()
	if !ctx.IsSet("input_index") {
		txinputindex = 0
	} else {
		txinputindex = ctx.Int64("input_index")
	}
	hextxbytes = args.First()
	args = args.Tail()
	hexpubkeybytes = args.First()
	rawtxbytes, err := hex.DecodeString(hextxbytes)
	if err != nil {
		fmt.Errorf("Failed to decode tx bytes from hex")
	}
	rawpubkeybytes, err := hex.DecodeString(hexpubkeybytes)
	if err != nil {
		fmt.Errorf("Failed to decode pubkey bytes from hex")
	}

	ctxb := context.Background()
	client, cleanUp := getSignClient(ctx)
	defer cleanUp()

	keyDesc := &signrpc.KeyDescriptor{
		RawKeyBytes: rawpubkeybytes,
	}
	var signDescs []*signrpc.SignDescriptor
	signDescs = append(signDescs, &signrpc.SignDescriptor{
		InputIndex: int32(txinputindex),
		KeyDesc:    keyDesc,
	})
	req := &signrpc.SignReq{
		RawTxBytes: rawtxbytes,
		SignDescs:  signDescs,
	}

	resp, err := client.SignOutputRaw(ctxb, req)
	if err != nil {
		return err
	}

	printJSON(resp)

	return nil
}
