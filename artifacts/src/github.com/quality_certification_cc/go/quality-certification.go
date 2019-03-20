package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "initUserInfo" {
		return s.initUserInfo(APIstub, args)
	}else if function == "userReview"{
		return s.userReview(APIstub,args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

//Initialize user information
func (s *SmartContract) initUserInfo(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	return shim.Success(nil)
}

//User review
func (s *SmartContract) userReview(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	return shim.Success(nil)

}
