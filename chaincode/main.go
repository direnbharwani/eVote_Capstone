package main

import (
	chaincode "github.com/direnbharwani/eVote_Capstone/src"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	eVoteSmartContract := new(chaincode.SmartContract)
	chaincode, err := contractapi.NewChaincode(eVoteSmartContract)
	if err != nil {
		panic(err.Error())
	}

	if err = chaincode.Start(); err != nil {
		panic(err.Error())
	}
}
