package cmd

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/quorumcontrol/tupelo-go-sdk/gossip3/remote"
	"github.com/quorumcontrol/tupelo-go-sdk/p2p"
	"github.com/spf13/cobra"
)

type service struct {
	publicKey *ecdsa.PublicKey
	actor     *actor.PID
}

func setupClientRemote(ctx context.Context, localKey *ecdsa.PrivateKey, services []*service) error {
	fmt.Printf("Setting up remote subsystem\n")
	remote.Start()

	p2pHost, err := p2p.NewLibP2PHost(ctx, localKey, int(servicePort))
	if err != nil {
		return err
	}

	// TODO: Find bootstrap peers
	fmt.Printf("Bootstrapping...\n")
	if _, err = p2pHost.Bootstrap([]string{}); err != nil {
		return err
	}

	if err = p2pHost.WaitForBootstrap(len(services), 15*time.Second); err != nil {
		return err
	}

	remote.NewRouter(p2pHost)

	return nil
}

func setupServices(ctx context.Context, localKey *ecdsa.PrivateKey) ([]*service, error) {
	conf, err := readConf()
	if err != nil {
		return nil, err
	}

	fromID, err := p2p.PeerFromEcdsaKey(&localKey.PublicKey)
	if err != nil {
		return nil, err
	}

	services := make([]*service, 0, len(conf.Services))
	for _, sc := range conf.Services {
		pubKeyB, err := hex.DecodeString(sc.EcdsaHexPublicKey)
		if err != nil {
			return nil, err
		}
		ecdsaPubKey, err := crypto.UnmarshalPubkey(pubKeyB)
		if err != nil {
			return nil, err
		}

		toId, err := p2p.PeerFromEcdsaKey(ecdsaPubKey)
		if err != nil {
			return nil, err
		}
		actorAddr := fmt.Sprintf("%s-%s", fromID.Pretty(), toId.Pretty())
		fmt.Printf("Creating service actor %s\n", actorAddr)
		services = append(services, &service{
			publicKey: ecdsaPubKey,
			actor:     actor.NewPID(actorAddr, fmt.Sprintf("tupelo-%s", toId)),
		})
	}

	return services, nil
}

var cmdClient = &cobra.Command{
	Use:   "client <destination>",
	Short: "Jason's Game land service client",
	Long:  "Jason's Game land service client",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("requires command argument")
		}
		// TODO: Validate key file
		// TODO: Validate command
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		serviceCmd := args[0]

		ctx := context.Background()
		localKey, err := crypto.GenerateKey()
		if err != nil {
			return err
		}

		services, err := setupServices(ctx, localKey)
		if err != nil {
			return err
		}

		if err := setupClientRemote(ctx, localKey, services); err != nil {
			return err
		}

		fmt.Printf("Successfully sent command: %q\n", serviceCmd)
		return nil
	},
}

var clientPort uint

func init() {
	cmdClient.Flags().UintVarP(&clientPort, "port", "p", 0, "Source port")
}
