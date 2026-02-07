package generator

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

type Config struct {
	LeadChar   byte
	LeadCount  int
	TrailChar  byte
	TrailCount int
	MaxWorkers int
}

type Result struct {
	PrivateKey string
	PublicKey  string
	Address    string
	CostTime   time.Duration
	Worker     int
}

func Generate(cfg *Config) Result {
	numCPU := runtime.NumCPU()
	workers := numCPU

	if cfg.MaxWorkers > 0 && cfg.MaxWorkers < numCPU {
		workers = cfg.MaxWorkers
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resultChan := make(chan Result, 1)

	var wg sync.WaitGroup
	wg.Add(workers)

	fmt.Printf(
		"🔨 Starting find address with %d workers | %d leading '%c' | %d trail '%c' \n",
		workers,
		cfg.LeadCount,
		cfg.LeadChar,
		cfg.TrailCount,
		cfg.TrailChar,
	)

	start := time.Now()

	for i := 0; i < workers; i++ {
		go worker(
			ctx,
			cancel,
			&wg,
			resultChan,
			cfg.LeadChar,
			cfg.LeadCount,
			cfg.TrailChar,
			cfg.TrailCount,
		)

	}

	result := <-resultChan
	result.CostTime = time.Since(start)
	result.Worker = workers
	fmt.Println("✅Successfully search address!!!")

	wg.Wait()

	return result
}

func worker(
	ctx context.Context,
	cancel context.CancelFunc,
	wg *sync.WaitGroup,
	resultChan chan<- Result,
	leadChar byte,
	leadCount int,
	trailChar byte,
	trailCount int,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			priv, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
			if err != nil {
				continue
			}

			addr := strings.ToLower(
				crypto.PubkeyToAddress(priv.PublicKey).Hex()[2:],
			)

			if matchRule(
				addr,
				leadChar,
				leadCount,
				trailChar,
				trailCount,
			) {
				res := Result{
					PrivateKey: hex.EncodeToString(crypto.FromECDSA(priv)),
					PublicKey:  hex.EncodeToString(crypto.FromECDSAPub(&priv.PublicKey)),
					Address:    "0x" + addr,
				}

				select {
				case resultChan <- res:
					cancel()
				default:
				}
				return
			}
		}
	}
}

func matchRule(
	addr string,
	leadChar byte,
	leadCount int,
	trailChar byte,
	trailCount int,
) bool {
	l := len(addr)

	if leadCount+trailCount > l {
		return false
	}

	for i := 0; i < leadCount; i++ {
		if addr[i] != leadChar {
			return false
		}
	}

	for i := 0; i < trailCount; i++ {
		if addr[l-1-i] != trailChar {
			return false
		}
	}

	return true
}
