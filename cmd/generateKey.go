package cmd

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)

var cmdGenerateKey = &cobra.Command{
	Use:   "generate-key",
	Short: "Generate key",
	Long:  `Generate peer key.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ecdsaKey, err := crypto.GenerateKey()
		if err != nil {
			return err
		}

		privKeyStr := hex.EncodeToString(crypto.FromECDSA(ecdsaKey))
		pubKeyStr := hex.EncodeToString(crypto.FromECDSAPub(&ecdsaKey.PublicKey))

		fmt.Printf("Private key: %s\nPublic key: %s\n\n", privKeyStr, pubKeyStr)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(cmdGenerateKey)
}
