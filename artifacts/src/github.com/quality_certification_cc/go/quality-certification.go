package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
        "github.com/quality_certification_cc/go/model"
	"math/big"
	"strings"
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
	if function == "userRegistration" {
		return s.userRegistration(APIstub, args)
	} else if function == "userAudit" {
		return s.userAudit(APIstub, args)
	} else if function == "userLogout" {
		return s.userLogout(APIstub)
	} else if function == "certApplication" {
		return s.certApplication(APIstub, args)
	} else if function == "docAudit" {
		return s.docAudit(APIstub, args)
	} else if function == "onsiteAudit" {
		return s.onsiteAudit(APIstub, args)
	} else if function == "certUpload" {
		return s.certUpload(APIstub, args)
	} else if function == "testDataUpload" {
		return s.testDataUpload(APIstub, args)
	} else if function == "trialRunDataUpload" {
		return s.trialRunDataUpload(APIstub, args)
	} else if function == "queryUser" {
		return s.queryUser(APIstub, args)
	} else if function == "queryCert" {
		return s.queryCert(APIstub, args)
	} else if function == "queryAllCerts" {
		return s.queryAllCerts(APIstub)
	} else if function == "queryAllUsers" {
		return s.queryAllUsers(APIstub)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) userRegistration(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}
	creatorByte, _ := APIstub.GetCreator()
	certStart := bytes.IndexAny(creatorByte, "-----BEGIN")
	if certStart == -1 {
		return shim.Error("No certificate found")
	}
	certText := creatorByte[certStart:]
	bl, _ := pem.Decode(certText)
	if bl == nil {
		return shim.Error("Could not decode the PEM structure")
	}

	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		return shim.Error("ParseCertificate failed")
	}
	uname := cert.Subject.CommonName

	var user = model.User{UnitNo: args[0], UnitType: args[1], UserType: args[2], Id: args[3], Name: args[4], State: 0}

	userAsBytes, _ := json.Marshal(user)
	APIstub.PutState(uname, userAsBytes)
	return shim.Success(nil)
}

func (s *SmartContract) userAudit(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	creatorByte, _ := APIstub.GetCreator()
	certStart := bytes.IndexAny(creatorByte, "-----BEGIN")
	if certStart == -1 {
		return shim.Error("No certificate found")
	}
	certText := creatorByte[certStart:]
	bl, _ := pem.Decode(certText)
	if bl == nil {
		return shim.Error("Could not decode the PEM structure")
	}

	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		return shim.Error("ParseCertificate failed")
	}
	uname := cert.Subject.CommonName

	userAsBytes, _ := APIstub.GetState(uname)
	userB := model.User{}
	json.Unmarshal(userAsBytes, &userB)
	if userB.UserType != "2" {
		return shim.Error("You are not an auditor！")
	}

	userAsBytes, _ = APIstub.GetState(args[0])
	userA := model.User{}

	json.Unmarshal(userAsBytes, &userA)
	userA.State = 1
	userA.AuditorID = userB.Id
	userA.AuditorName = userB.Name

	userAsBytes, _ = json.Marshal(userA)
	APIstub.PutState(args[0], userAsBytes)
	return shim.Success(nil)

}

func (s *SmartContract) userLogout(APIstub shim.ChaincodeStubInterface) sc.Response {
	creatorByte, _ := APIstub.GetCreator()
	certStart := bytes.IndexAny(creatorByte, "-----BEGIN")
	if certStart == -1 {
		return shim.Error("No certificate found")
	}
	certText := creatorByte[certStart:]
	bl, _ := pem.Decode(certText)
	if bl == nil {
		return shim.Error("Could not decode the PEM structure")
	}

	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		return shim.Error("ParseCertificate failed")
	}
	uname := cert.Subject.CommonName

	userAsBytes, _ := APIstub.GetState(uname)
	user := model.User{}
	user.State = 2
	userAsBytes, _ = json.Marshal(user)
	APIstub.PutState(uname, userAsBytes)
	return shim.Success(nil)

}

func (s *SmartContract) certApplication(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}
	//验证签名
	var m, n big.Int
	var rr, ss *big.Int
	creatorByte, _ := APIstub.GetCreator()
	block, _ := pem.Decode(creatorByte)
	if block == nil {
		return shim.Error("block nil!")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return shim.Error("x509 parse err!")
	}
	pub := cert.PublicKey.(*ecdsa.PublicKey)
	h2 := sha256.New()
	h2.Write([]byte(args[0]))
	hashed := h2.Sum(nil)
	arr := strings.Split(args[5], ",")
	m.SetString(arr[0], 10) //大于int64的数字要用到SetString函数
	n.SetString(arr[1], 10)
	rr = &m
	ss = &n
	result := ecdsa.Verify(pub, hashed, rr, ss)
	if result != true {
		return shim.Error("ECDSA Verification failed")
	}

	certStart := bytes.IndexAny(creatorByte, "-----BEGIN")
	if certStart == -1 {
		return shim.Error("No certificate found")
	}
	certText := creatorByte[certStart:]
	bl, _ := pem.Decode(certText)
	if bl == nil {
		return shim.Error("Could not decode the PEM structure")
	}

	uname := cert.Subject.CommonName

	userAsBytes, _ := APIstub.GetState(uname)
	user := model.User{}
	certApp := model.CertApplication{BaseData: args[0], ApplyScanHASH: args[1], LegalPersonScanHASH: args[2],
		Summary: args[3], EncryptedSummary: args[4]}
	user.CertApplications = append(user.CertApplications, &certApp)
	userAsBytes, _ = json.Marshal(user)
	APIstub.PutState(uname, userAsBytes)
	return shim.Success(nil)

}

func (s *SmartContract) docAudit(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	//验证签名
	var m, n big.Int
	var rr, ss *big.Int
	creatorByte, _ := APIstub.GetCreator()
	block, _ := pem.Decode(creatorByte)
	if block == nil {
		return shim.Error("block nil!")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return shim.Error("x509 parse err!")
	}
	pub := cert.PublicKey.(*ecdsa.PublicKey)
	h2 := sha256.New()
	h2.Write([]byte(args[0]))
	hashed := h2.Sum(nil)
	arr := strings.Split(args[5], ",")
	m.SetString(arr[0], 10) //大于int64的数字要用到SetString函数
	n.SetString(arr[1], 10)
	rr = &m
	ss = &n
	result := ecdsa.Verify(pub, hashed, rr, ss)
	if result != true {
		return shim.Error("ECDSA Verification failed")
	}

	certStart := bytes.IndexAny(creatorByte, "-----BEGIN")
	if certStart == -1 {
		return shim.Error("No certificate found")
	}
	certText := creatorByte[certStart:]
	bl, _ := pem.Decode(certText)
	if bl == nil {
		return shim.Error("Could not decode the PEM structure")
	}

	uname := cert.Subject.CommonName

	userAsBytes, _ := APIstub.GetState(uname)
	user := model.User{}
	certApp := model.DocAudit{BaseData: args[0], EncryptedSummary: args[1]}
	user.DocAudits = append(user.DocAudits, &certApp)
	userAsBytes, _ = json.Marshal(user)
	APIstub.PutState(uname, userAsBytes)
	return shim.Success(nil)

}

func (s *SmartContract) onsiteAudit(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	//验证签名
	var m, n big.Int
	var rr, ss *big.Int
	creatorByte, _ := APIstub.GetCreator()
	block, _ := pem.Decode(creatorByte)
	if block == nil {
		return shim.Error("block nil!")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return shim.Error("x509 parse err!")
	}
	pub := cert.PublicKey.(*ecdsa.PublicKey)
	h2 := sha256.New()
	h2.Write([]byte(args[0]))
	hashed := h2.Sum(nil)
	arr := strings.Split(args[5], ",")
	m.SetString(arr[0], 10) //大于int64的数字要用到SetString函数
	n.SetString(arr[1], 10)
	rr = &m
	ss = &n
	result := ecdsa.Verify(pub, hashed, rr, ss)
	if result != true {
		return shim.Error("ECDSA Verification failed")
	}

	certStart := bytes.IndexAny(creatorByte, "-----BEGIN")
	if certStart == -1 {
		return shim.Error("No certificate found")
	}
	certText := creatorByte[certStart:]
	bl, _ := pem.Decode(certText)
	if bl == nil {
		return shim.Error("Could not decode the PEM structure")
	}

	uname := cert.Subject.CommonName

	userAsBytes, _ := APIstub.GetState(uname)
	user := model.User{}
	certApp := model.OnsiteAudit{BaseData: args[0], EncryptedSummary: args[1]}
	user.OnsiteAudits = append(user.OnsiteAudits, &certApp)
	userAsBytes, _ = json.Marshal(user)
	APIstub.PutState(uname, userAsBytes)
	return shim.Success(nil)

}

func (s *SmartContract) certUpload(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	//验证签名
	var m, n big.Int
	var rr, ss *big.Int
	creatorByte, _ := APIstub.GetCreator()
	block, _ := pem.Decode(creatorByte)
	if block == nil {
		return shim.Error("block nil!")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return shim.Error("x509 parse err!")
	}
	pub := cert.PublicKey.(*ecdsa.PublicKey)
	h2 := sha256.New()
	h2.Write([]byte(args[2]))
	hashed := h2.Sum(nil)
	arr := strings.Split(args[3], ",")
	m.SetString(arr[0], 10) //大于int64的数字要用到SetString函数
	n.SetString(arr[1], 10)
	rr = &m
	ss = &n
	result := ecdsa.Verify(pub, hashed, rr, ss)
	if result != true {
		return shim.Error("ECDSA Verification failed")
	}
	uname := cert.Subject.CommonName
	userAsBytes, _ := APIstub.GetState(uname)
	user := model.User{}
	json.Unmarshal(userAsBytes, &user)

	certUpload := model.CertUpload{BaseData: args[2], EncryptedSummary: args[3], PostPersonID: user.Id, PostPersonName: user.Name}
	var certData = model.CertificationData{CertificateID: args[0], UnitID: args[1], CertUpload: &certUpload}
	certAsBytes, _ := json.Marshal(certData)
	key := fmt.Sprintf("%s,%s", args[0], args[1])
	APIstub.PutState(key, certAsBytes)
	return shim.Success(nil)

}

func (s *SmartContract) testDataUpload(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	//验证签名
	var m, n big.Int
	var rr, ss *big.Int
	creatorByte, _ := APIstub.GetCreator()
	block, _ := pem.Decode(creatorByte)
	if block == nil {
		return shim.Error("block nil!")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return shim.Error("x509 parse err!")
	}
	pub := cert.PublicKey.(*ecdsa.PublicKey)
	h2 := sha256.New()
	h2.Write([]byte(args[2]))
	hashed := h2.Sum(nil)
	arr := strings.Split(args[3], ",")
	m.SetString(arr[0], 10) //大于int64的数字要用到SetString函数
	n.SetString(arr[1], 10)
	rr = &m
	ss = &n
	result := ecdsa.Verify(pub, hashed, rr, ss)
	if result != true {
		return shim.Error("ECDSA Verification failed")
	}
	uname := cert.Subject.CommonName
	userAsBytes, _ := APIstub.GetState(uname)
	user := model.User{}
	json.Unmarshal(userAsBytes, &user)

	certUpload := model.TestDataUpload{BaseData: args[2], EncryptedSummary: args[3], PostPersonID: user.Id, PostPersonName: user.Name}
	var certData = model.CertificationData{CertificateID: args[0], UnitID: args[1], TestDataUpload: &certUpload}
	certAsBytes, _ := json.Marshal(certData)
	key := fmt.Sprintf("%s,%s", args[0], args[1])
	APIstub.PutState(key, certAsBytes)
	return shim.Success(nil)

}

func (s *SmartContract) trialRunDataUpload(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	//验证签名
	var m, n big.Int
	var rr, ss *big.Int
	creatorByte, _ := APIstub.GetCreator()
	block, _ := pem.Decode(creatorByte)
	if block == nil {
		return shim.Error("block nil!")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return shim.Error("x509 parse err!")
	}
	pub := cert.PublicKey.(*ecdsa.PublicKey)
	h2 := sha256.New()
	h2.Write([]byte(args[2]))
	hashed := h2.Sum(nil)
	arr := strings.Split(args[3], ",")
	m.SetString(arr[0], 10) //大于int64的数字要用到SetString函数
	n.SetString(arr[1], 10)
	rr = &m
	ss = &n
	result := ecdsa.Verify(pub, hashed, rr, ss)
	if result != true {
		return shim.Error("ECDSA Verification failed")
	}

	uname := cert.Subject.CommonName
	userAsBytes, _ := APIstub.GetState(uname)
	user := model.User{}
	json.Unmarshal(userAsBytes, &user)

	certUpload := model.TrialRunDataUpload{BaseData: args[2], EncryptedSummary: args[3], PostPersonID: user.Id, PostPersonName: user.Name}
	var certData = model.CertificationData{CertificateID: args[0], UnitID: args[1], TrialRunDataUpload: &certUpload}
	certAsBytes, _ := json.Marshal(certData)
	key := fmt.Sprintf("%s,%s", args[0], args[1])
	APIstub.PutState(key, certAsBytes)
	return shim.Success(nil)

}

func (s *SmartContract) queryUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	userAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(userAsBytes)
}

func (s *SmartContract) queryCert(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	key := fmt.Sprintf("%s,%s", args[0], args[1])
	carAsBytes, _ := APIstub.GetState(key)
	return shim.Success(carAsBytes)
}

func (s *SmartContract) queryAllUsers(APIstub shim.ChaincodeStubInterface) sc.Response {
	var queryString = "{\"selector\":{\"id\":{\"$regex\":\"(?i)\"}}}"

	resultsIterator, err := APIstub.GetQueryResult(queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")

		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- allUsers:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())

}

func (s *SmartContract) queryAllCerts(APIstub shim.ChaincodeStubInterface) sc.Response {
	var queryString = "{\"selector\":{\"certificateID\":{\"$regex\":\"(?i)\"}}}"

	resultsIterator, err := APIstub.GetQueryResult(queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")

		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- allUsers:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())

}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
