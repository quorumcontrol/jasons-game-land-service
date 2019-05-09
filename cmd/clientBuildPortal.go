package cmd

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/quorumcontrol/jasons-game-land-service/messages"
	"github.com/quorumcontrol/jasons-game-land-service/p2p"
	"github.com/quorumcontrol/tupelo-go-sdk/gossip3/remote"
	sdkp2p "github.com/quorumcontrol/tupelo-go-sdk/p2p"
	"github.com/spf13/cobra"
)

type service struct {
	publicKey *ecdsa.PublicKey
	actor     *actor.PID
}

func setupClientRemote(ctx context.Context, localKey *ecdsa.PrivateKey) error {
	remote.Start()

	p2pHost, err := sdkp2p.NewLibP2PHost(ctx, localKey, int(clientPort))
	if err != nil {
		return err
	}
	if err = p2p.Bootstrap(p2pHost); err != nil {
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

	fromID, err := sdkp2p.PeerFromEcdsaKey(&localKey.PublicKey)
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

		toAddr := crypto.PubkeyToAddress(*ecdsaPubKey).String()
		toId, err := sdkp2p.PeerFromEcdsaKey(ecdsaPubKey)
		if err != nil {
			return nil, err
		}
		actorAddr := fmt.Sprintf("%s-%s", fromID.Pretty(), toId.Pretty())
		act := actor.NewPID(actorAddr, fmt.Sprintf("service-%s", toAddr))
		fmt.Printf("Creating service actor with ID %q and address %q\n", act.Id, act.Address)
		services = append(services, &service{
			publicKey: ecdsaPubKey,
			actor:     act,
		})
	}

	return services, nil
}

var cmdBuildPortal = &cobra.Command{
	Use:   "build-portal",
	Short: "Build portal",
	Long:  "Send request to service to build a portal.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		localKey, err := crypto.GenerateKey()
		if err != nil {
			return err
		}

		services, err := setupServices(ctx, localKey)
		if err != nil {
			return err
		}

		if err := setupClientRemote(ctx, localKey); err != nil {
			return err
		}

		fut := actor.EmptyRootContext.RequestFuture(services[0].actor, &messages.BuildPortal{},
			1*time.Second)
		if err = fut.Wait(); err != nil {
			return errors.Wrapf(err, "request to build portal failed")
		}

		fmt.Printf("Successfully requested building of portal\n")
		return nil
	},
}

func init() {
	cmdClient.AddCommand(cmdBuildPortal)
}
