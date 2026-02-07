package test

import (
	"EthAddressGenerator/generator"
	"fmt"
	"testing"
)

func TestFindETHAddress(t *testing.T) {
	testConfig := &generator.Config{
		LeadChar:   '8',
		LeadCount:  2,
		TrailChar:  '8',
		TrailCount: 2,
		MaxWorkers: 20,
	}

	res := generator.Generate(testConfig)
	fmt.Println("found address    : ", res.Address)
	fmt.Println("found private key: ", res.PrivateKey)
	fmt.Println("found public key : ", res.PublicKey)
	fmt.Println("search cost      : ", res.CostTime)
	fmt.Println("search worker    : ", res.Worker)
}
