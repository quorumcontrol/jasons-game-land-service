package p2p

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/quorumcontrol/jasons-game-land-service/config"
	"github.com/quorumcontrol/tupelo-go-sdk/p2p"
)

func BootstrapperAddrs(conf config.Configuration) ([]string, error) {
	bootstrapperAddrs := []string{}
	// For now assume that bootstrappers run locally
	for _, bc := range conf.Bootstrappers {
		ecdsaKeyB, err := hex.DecodeString(bc.EcdsaHexPrivateKey)
		if err != nil {
			return nil, err
		}
		ecdsaKey, err := crypto.ToECDSA(ecdsaKeyB)
		if err != nil {
			return nil, err
		}

		peerId, err := p2p.PeerFromEcdsaKey(&ecdsaKey.PublicKey)
		if err != nil {
			return nil, err
		}
		addr := fmt.Sprintf("/ip4/127.0.0.1/tcp/34001/ipfs/%s", peerId)
		bootstrapperAddrs = append(bootstrapperAddrs, addr)
	}

	return bootstrapperAddrs, nil
}

func Bootstrap(p2pHost *p2p.LibP2PHost, conf config.Configuration) error {
	bootstrapperAddrs, err := BootstrapperAddrs(conf)
	if err != nil {
		return err
	}
	if _, err := p2pHost.Bootstrap(bootstrapperAddrs); err != nil {
		return err
	}
	if err := p2pHost.WaitForBootstrap(1, 20*time.Second); err != nil {
		return err
	}

	return nil
}
