package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"EthAddressGenerator/generator"
)

type Output struct {
	PrivateKey   string    `json:"private_key"`
	PublicKey    string    `json:"public_key"`
	Address      string    `json:"address"`
	CostTime     string    `json:"cost_time_ms"`
	WorkerUsed   int       `json:"worker_used"`
	SearchTimeAt time.Time `json:"search_time_at"`
}

func main() {
	// 前缀字符
	leadChar := flag.String("lead-char", "", "leading hex character (0-9, a-f)")
	// 前缀字符数
	leadCount := flag.Int("lead-count", 0, "number of leading characters")
	// 后缀字符
	trailChar := flag.String("trail-char", "", "trailing hex character (0-9, a-f)")
	// 后缀字符数
	trailCount := flag.Int("trail-count", 0, "number of trailing characters")
	// 最大工作线程数，未设置或数值大于当前进程可见CPU核数时，默认使用CPU核数
	workers := flag.Int("workers", 0, "max workers (default: NumCPU)")

	flag.Parse()

	if *leadCount > 0 {
		if len(*leadChar) != 1 || !isHexChar((*leadChar)[0]) {
			fmt.Println("❌ invalid lead-char")
			os.Exit(1)
		}
	}

	if *trailCount > 0 {
		if len(*trailChar) != 1 || !isHexChar((*trailChar)[0]) {
			fmt.Println("❌ invalid trail-char")
		}
	}

	var leadCharByte byte
	var trailCharByte byte

	if *leadCount > 0 {
		leadCharByte = strings.ToLower(*leadChar)[0]
	}

	if *trailCount > 0 {
		trailCharByte = strings.ToLower(*trailChar)[0]
	}

	if *leadCount == 0 && *trailCount == 0 {
		fmt.Println("❌ at least one of lead-count or trail-count must be > 0")
		os.Exit(1)
	}

	cfg := &generator.Config{
		LeadChar:   leadCharByte,
		LeadCount:  *leadCount,
		TrailChar:  trailCharByte,
		TrailCount: *trailCount,
		MaxWorkers: *workers,
	}

	result := generator.Generate(cfg)

	output := Output{
		PrivateKey:   result.PrivateKey,
		PublicKey:    result.PublicKey,
		Address:      result.Address,
		CostTime:     fmt.Sprintf("%.4fms", float64(result.CostTime.Microseconds())/1000),
		SearchTimeAt: time.Now(),
		WorkerUsed:   result.Worker,
	}

	filename := "search_results.json"
	if err := appendJSONFile(filename, output); err != nil {
		fmt.Println("❌ failed to append result to file:", err)
		fmt.Println("💣💣💣Private Key: ", output.PrivateKey)
		os.Exit(1)
	}

	fmt.Printf("✅ Result appended to %s\n", filename)
}

func isHexChar(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')
}

func appendJSONFile(filename string, newOutput Output) error {
	var outputs []Output

	if data, err := os.ReadFile(filename); err == nil {
		if len(data) > 0 {
			if err := json.Unmarshal(data, &outputs); err != nil {
				outputs = nil
			}
		}
	}

	outputs = append(outputs, newOutput)

	dataToWrite, err := json.MarshalIndent(outputs, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, dataToWrite, 0600)
}
