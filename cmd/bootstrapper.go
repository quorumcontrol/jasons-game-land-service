package cmd

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/quorumcontrol/tupelo-go-sdk/p2p"
	"github.com/spf13/cobra"
)

var bootstrapperPort uint

var cmdBootstrapper = &cobra.Command{
	Use:   "bootstrapper",
	Short: "Run a bootstrap node",
	RunE: func(cmd *cobra.Command, args []string) error {
		conf, err := readConf()
		if err != nil {
			return err
		}
		// TODO: Detect node identity
		nodeIdx := 0

		ecdsaKeyB, err := hex.DecodeString(conf.Bootstrappers[nodeIdx].EcdsaHexPrivateKey)
		if err != nil {
			return err
		}
		ecdsaKey, err := crypto.ToECDSA(ecdsaKeyB)
		if err != nil {
			return err
		}

		ctx := context.Background()
		host, err := p2p.NewRelayLibP2PHost(ctx, ecdsaKey, int(bootstrapperPort))
		if err != nil {
			return err
		}

		otherBootstrappers := []string{}
		// for i, addr := range p2p.BootstrapNodes() {
		// 	if i == nodeIdx {
		// 		continue
		// 	}

		// 	otherBootstrappers = append(otherBootstrappers, addr)
		// }
		if len(otherBootstrappers) > 0 {
			if _, err = host.Bootstrap(otherBootstrappers); err != nil {
				return err
			}
		}

		fmt.Println("Bootstrap node running at:")
		for _, addr := range host.Addresses() {
			fmt.Println(addr)
		}
		select {}
	},
}

func init() {
	rootCmd.AddCommand(cmdBootstrapper)
	cmdBootstrapper.Flags().UintVarP(&bootstrapperPort, "port", "p", 0,
		"What port to use (default random)")
}
