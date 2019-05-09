package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/quorumcontrol/jasons-game-land-service/p2p"
	srv "github.com/quorumcontrol/jasons-game-land-service/service"
	"github.com/quorumcontrol/tupelo-go-sdk/gossip3/remote"
	sdkp2p "github.com/quorumcontrol/tupelo-go-sdk/p2p"
	"github.com/spf13/cobra"
)

var servicePort uint

func stopOnSignal(actors ...*actor.PID) {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Printf("Caught signal %d\n", sig)
		for _, act := range actors {
			fmt.Printf("Gracefully stopping actor\n")
			err := actor.EmptyRootContext.PoisonFuture(act).Wait()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Actor failed to stop gracefully: %s\n", err)
			}
		}
		done <- true
	}()
	<-done
}

func setupServiceRemote(ctx context.Context, nodeIdx int) (*actor.PID, error) {
	fmt.Printf("Setting up remote subsystem\n")
	remote.Start()

	conf, err := readConf()
	if err != nil {
		return nil, err
	}
	ecdsaKeyB, err := hex.DecodeString(conf.Services[nodeIdx].EcdsaHexPrivateKey)
	if err != nil {
		return nil, err
	}
	ecdsaKey, err := crypto.ToECDSA(ecdsaKeyB)
	if err != nil {
		return nil, err
	}
	peerId := crypto.PubkeyToAddress(ecdsaKey.PublicKey).String()

	p2pHost, err := sdkp2p.NewLibP2PHost(ctx, ecdsaKey, int(servicePort))
	if err != nil {
		return nil, err
	}

	if err := p2p.Bootstrap(p2pHost); err != nil {
		return nil, err
	}

	remote.NewRouter(p2pHost)

	cfg := srv.ServiceConfig{
		// Self:              localSigner,
		// NotaryGroup:       group,
		// CurrentStateStore: badgerCurrent,
		// PubSubSystem:      remote.NewNetworkPubSub(p2pHost),
	}
	name := "service-" + peerId
	props := actor.PropsFromProducer(func() actor.Actor {
		return srv.NewServiceActor(cfg)
	})
	act, err := actor.EmptyRootContext.SpawnNamed(props, name)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Started actor with ID %q and address %q\n", act.Id, act.Address)

	return act, nil
}

var cmdService = &cobra.Command{
	Use:   "service",
	Short: "Service",
	Long:  `Service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		// TODO: Detect node identity
		nodeIdx := 0
		fmt.Printf("Service #%d starting\n", nodeIdx+1)
		act, err := setupServiceRemote(ctx, nodeIdx)
		if err != nil {
			return err
		}

		fmt.Printf("All set up!\n")

		stopOnSignal(act)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cmdService)
	cmdService.Flags().UintVarP(&servicePort, "port", "p", 0, "Source port")
}
