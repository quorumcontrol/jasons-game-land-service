package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

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

		fmt.Printf("Bootstrapper #%d starting\n", nodeIdx+1)

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
		anAddr := host.Addresses()[0].String()
		keySlice := strings.Split(anAddr, "/")
		key := keySlice[len(keySlice)-1]

		fmt.Printf("anAddr: %q, key: %q\n", anAddr, key)
		otherBootstrappers := []string{}
		// for i, addr := range p2p.BootstrapNodes() {
		// 	if i == nodeIdx {
		// 		continue
		// 	}

		// 	otherBootstrappers = append(otherBootstrappers, addr)
		// }
		if len(otherBootstrappers) > 0 {
			fmt.Printf("Bootstrapping self\n")
			if _, err = host.Bootstrap(otherBootstrappers); err != nil {
				return err
			}
		} else {
			fmt.Printf("Not bootstrapping self\n")
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
