package p2p

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/quorumcontrol/tupelo-go-sdk/p2p"
)

func Bootstrap(p2pHost *p2p.LibP2PHost) error {
	bootstrapperAddrsRaw := strings.Split(os.Getenv("LAND_SERVICE_BOOTSTRAPPERS"), ",")
	bootstrapperAddrs := []string{}
	for _, addr := range bootstrapperAddrsRaw {
		addr = strings.TrimSpace(addr)
		if addr != "" {
			bootstrapperAddrs = append(bootstrapperAddrs, addr)
		}
	}
	if len(bootstrapperAddrs) == 0 {
		return fmt.Errorf("please define $LAND_SERVICE_BOOTSTRAPPERS")
	}
	fmt.Printf("Bootstrapping...\n")
	if _, err := p2pHost.Bootstrap(bootstrapperAddrs); err != nil {
		return err
	}
	if err := p2pHost.WaitForBootstrap(1, 20*time.Second); err != nil {
		return err
	}
	fmt.Printf("Finished bootstrapping\n")

	return nil
}
