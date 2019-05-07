package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	crypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/spf13/cobra"
)

var cmdGenerateKey = &cobra.Command{
	Use:   "generate-key",
	Short: "Generate key",
	Long:  `Generate peer key.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		privKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
		if err != nil {
			return err
		}

		bytes, err := crypto.MarshalPrivateKey(privKey)
		if err != nil {
			return err
		}
		privKeyStr := base64.StdEncoding.EncodeToString(bytes)

		id, err := peer.IDFromPrivateKey(privKey)
		if err != nil {
			return err
		}
		fmt.Printf("Private key: %s\nPeer ID: %s\n", privKeyStr, id.Pretty())

		return nil
	},
}
