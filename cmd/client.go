package cmd

import (
	"bufio"
	"context"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-crypto"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/multiformats/go-multiaddr"
	"github.com/spf13/cobra"
)

var cmdClient = &cobra.Command{
	Use:   "client <destination>",
	Short: "Jason's Game land service client",
	Long:  "Jason's Game land service client",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("requires destination and command arguments")
		}
		// TODO: Validate address
		// TODO: Validate command
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		dest := args[0]
		serviceCmd := args[1]

		ctx := context.Background()

		r := rand.Reader
		privKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
		if err != nil {
			return err
		}

		sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", clientPort))
		host, err := libp2p.New(ctx, libp2p.ListenAddrs(sourceMultiAddr),
			libp2p.Identity(privKey))
		if err != nil {
			return err
		}

		fmt.Println("This node's multiaddresses:")
		for _, la := range host.Addrs() {
			fmt.Printf(" - %v\n", la)
		}
		fmt.Println()

		maddr, err := multiaddr.NewMultiaddr(dest)
		if err != nil {
			return err
		}
		info, err := peerstore.InfoFromP2pAddr(maddr)
		if err != nil {
			return err
		}

		host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

		s, err := host.NewStream(context.Background(), info.ID, "/jasons-game-land-service/1.0.0")
		if err != nil {
			return err
		}

		w := bufio.NewWriter(s)
		if _, err = w.WriteString(fmt.Sprintf("%s\n", serviceCmd)); err != nil {
			return err
		}
		w.Flush()

		fmt.Printf("Successfully sent command: %q\n", serviceCmd)
		return nil
	},
}

var clientPort uint

func init() {
	cmdClient.Flags().UintVarP(&clientPort, "port", "p", 0, "Source port")
}
