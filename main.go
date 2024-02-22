package main

import (
	"fmt"
	"os"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/milebdo/wallet/wallet"
)

func main() {

	server := &shim.ChaincodeServer{
		CCID:    os.Getenv("CORE_CHAINCODE_ID"),
		Address: os.Getenv("CORE_CHAINCODE_ADDRESS"),
		CC:      new(wallet.Wallet),
		TLSProps: shim.TLSProperties{
			Disabled: true,
		},
	}
	err := server.Start()
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}

}
