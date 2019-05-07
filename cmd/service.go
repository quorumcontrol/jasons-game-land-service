package cmd

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-crypto"
	net "github.com/libp2p/go-libp2p-net"
	"github.com/multiformats/go-multiaddr"
	"github.com/spf13/cobra"
)

var servicePort uint

func readData(rw *bufio.ReadWriter) {
	str, _ := rw.ReadString('\n')
	if str == "" || str == "\n" {
		return
	}

	// Green console colour: 	\x1b[32m
	// Reset console colour: 	\x1b[0m
	fmt.Printf("\x1b[32m%s\x1b[0m", str)
}

func handleStream(s net.Stream) {
	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go readData(rw)
}

var cmdService = &cobra.Command{
	Use:   "service",
	Short: "Service",
	Long:  `Service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		var privKey crypto.PrivKey
		privKeyStr := os.Getenv("LAND_SERVICE_PRIVATE_KEY")
		privKeyStr = strings.TrimSpace(privKeyStr)
		if privKeyStr != "" {
			privKeyB, err := base64.StdEncoding.DecodeString(privKeyStr)
			if err != nil {
				return err
			}
			privKey, err = crypto.UnmarshalPrivateKey(privKeyB)
			if err != nil {
				return err
			}
		} else {
			r := rand.Reader
			var err error
			privKey, _, err = crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
			if err != nil {
				return err
			}
		}

		sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", servicePort))
		host, err := libp2p.New(ctx, libp2p.ListenAddrs(sourceMultiAddr),
			libp2p.Identity(privKey))
		if err != nil {
			return err
		}

		host.SetStreamHandler("/jasons-game-land-service/1.0.0", handleStream)

		// Let's get the actual TCP port from our listen multiaddr, in case we're using 0 (default; random available port).
		var port string
		for _, la := range host.Network().ListenAddresses() {
			if p, err := la.ValueForProtocol(multiaddr.P_TCP); err == nil {
				port = p
				break
			}
		}

		if port == "" {
			panic("was not able to find actual local port")
		}

		fmt.Printf("Run './jasons-game-land-service client /ip4/127.0.0.1/tcp/%v/p2p/%s <command>' in another console.\n", port, host.ID().Pretty())
		fmt.Println("You can replace 127.0.0.1 with public IP as well.")
		fmt.Printf("\nWaiting for incoming connection\n\n")

		// Hang forever
		<-make(chan struct{})

		return nil
	},
}

func init() {
	cmdService.Flags().UintVarP(&servicePort, "port", "p", 0, "Source port")
}
