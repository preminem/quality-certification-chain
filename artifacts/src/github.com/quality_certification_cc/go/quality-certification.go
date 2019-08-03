package main

import (
	"bytes"
	//"crypto/ecdsa"
	//"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	"github.com/quality_certification_cc/go/model"
	//"github.com/preminem/quality-certification-chain/artifacts/src/github.com/quality_certification_cc/go/model"
	//"math/big"
	//"strings"
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
	} else if function == "conditionalQuery" {
		return s.conditionalQuery(APIstub, args)
	} else if function == "publicQuery" {
		return s.publicQuery(APIstub, args)
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
	json.Unmarshal(userAsBytes, &user)
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
	//var m, n big.Int
	//var rr, ss *big.Int
	//提取证书
	creatorByte, _ := APIstub.GetCreator()
	certStart := bytes.IndexAny(creatorByte, "-----BEGIN")
	if certStart == -1 {
		return shim.Error("No certificate found")
	}
	certText := creatorByte[certStart:]
	//提取用户身份
	bl, _ := pem.Decode(certText)
	if bl == nil {
		return shim.Error("Could not decode the PEM structure")
	}
	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		return shim.Error("ParseCertificate failed")
	}
	uname := cert.Subject.CommonName
	//pub := cert.PublicKey.(*ecdsa.PublicKey)
	//h2 := sha256.New()
	//h2.Write([]byte(args[0]))
	//hashed := h2.Sum(nil)
	//arr := strings.Split(args[5], ",")
	//m.SetString(arr[0], 10) //大于int64的数字要用到SetString函数
	//n.SetString(arr[1], 10)
	//rr = &m
	//ss = &n
	//result := ecdsa.Verify(pub, hashed, rr, ss)
	//if result != true {
	//	return shim.Error("ECDSA Verification failed")
	//}

	//certStart := bytes.IndexAny(creatorByte, "-----BEGIN")
	//if certStart == -1 {
	//	return shim.Error("No certificate found")
	//}
	//certText := creatorByte[certStart:]
	//bl, _ := pem.Decode(certText)
	//if bl == nil {
	//	return shim.Error("Could not decode the PEM structure")
	//}

	userAsBytes, _ := APIstub.GetState(uname)
	user := model.User{}
	json.Unmarshal(userAsBytes, &user)
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
	//var m, n big.Int
	//var rr, ss *big.Int
	//提取证书
	creatorByte, _ := APIstub.GetCreator()
	certStart := bytes.IndexAny(creatorByte, "-----BEGIN")
	if certStart == -1 {
		return shim.Error("No certificate found")
	}
	certText := creatorByte[certStart:]
	//提取用户身份
	bl, _ := pem.Decode(certText)
	if bl == nil {
		return shim.Error("Could not decode the PEM structure")
	}
	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		return shim.Error("ParseCertificate failed")
	}
	uname := cert.Subject.CommonName
	//pub := cert.PublicKey.(*ecdsa.PublicKey)
	//h2 := sha256.New()
	//h2.Write([]byte(args[0]))
	//hashed := h2.Sum(nil)
	//arr := strings.Split(args[5], ",")
	//m.SetString(arr[0], 10) //大于int64的数字要用到SetString函数
	//n.SetString(arr[1], 10)
	//rr = &m
	//ss = &n
	//result := ecdsa.Verify(pub, hashed, rr, ss)
	//if result != true {
	//	return shim.Error("ECDSA Verification failed")
	//}

	//certStart := bytes.IndexAny(creatorByte, "-----BEGIN")
	//if certStart == -1 {
	//	return shim.Error("No certificate found")
	//}
	//certText := creatorByte[certStart:]
	//bl, _ := pem.Decode(certText)
	//if bl == nil {
	//	return shim.Error("Could not decode the PEM structure")
	//}

	userAsBytes, _ := APIstub.GetState(uname)
	user := model.User{}
	json.Unmarshal(userAsBytes, &user)

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
	//var m, n big.Int
	//var rr, ss *big.Int
	//提取证书
	creatorByte, _ := APIstub.GetCreator()
	certStart := bytes.IndexAny(creatorByte, "-----BEGIN")
	if certStart == -1 {
		return shim.Error("No certificate found")
	}
	certText := creatorByte[certStart:]
	//提取用户身份
	bl, _ := pem.Decode(certText)
	if bl == nil {
		return shim.Error("Could not decode the PEM structure")
	}
	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		return shim.Error("ParseCertificate failed")
	}
	uname := cert.Subject.CommonName
	//pub := cert.PublicKey.(*ecdsa.PublicKey)
	//h2 := sha256.New()
	//h2.Write([]byte(args[0]))
	//hashed := h2.Sum(nil)
	//arr := strings.Split(args[5], ",")
	//m.SetString(arr[0], 10) //大于int64的数字要用到SetString函数
	//n.SetString(arr[1], 10)
	//rr = &m
	//ss = &n
	//result := ecdsa.Verify(pub, hashed, rr, ss)
	//if result != true {
	//	return shim.Error("ECDSA Verification failed")
	//}
	//
	//certStart := bytes.IndexAny(creatorByte, "-----BEGIN")
	//if certStart == -1 {
	//	return shim.Error("No certificate found")
	//}
	//certText := creatorByte[certStart:]
	//bl, _ := pem.Decode(certText)
	//if bl == nil {
	//	return shim.Error("Could not decode the PEM structure")
	//}

	userAsBytes, _ := APIstub.GetState(uname)
	user := model.User{}
	json.Unmarshal(userAsBytes, &user)

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
	//var m, n big.Int
	//var rr, ss *big.Int
	creatorByte, _ := APIstub.GetCreator()
	certStart := bytes.IndexAny(creatorByte, "-----BEGIN")
	if certStart == -1 {
		return shim.Error("No certificate found")
	}
	certText := creatorByte[certStart:]
	//提取用户身份
	bl, _ := pem.Decode(certText)
	if bl == nil {
		return shim.Error("Could not decode the PEM structure")
	}
	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		return shim.Error("ParseCertificate failed")
	}
	uname := cert.Subject.CommonName
	//pub := cert.PublicKey.(*ecdsa.PublicKey)
	//h2 := sha256.New()
	//h2.Write([]byte(args[2]))
	//hashed := h2.Sum(nil)
	//arr := strings.Split(args[3], ",")
	//m.SetString(arr[0], 10) //大于int64的数字要用到SetString函数
	//n.SetString(arr[1], 10)
	//rr = &m
	//ss = &n
	//result := ecdsa.Verify(pub, hashed, rr, ss)
	//if result != true {
	//	return shim.Error("ECDSA Verification failed")
	//}
	userAsBytes, _ := APIstub.GetState(uname)
	user := model.User{}
	json.Unmarshal(userAsBytes, &user)

	certUpload := model.CertUpload{BaseData: args[2], EncryptedSummary: args[3], PostPersonID: user.Id, PostPersonName: user.Name}
	var certData = model.CertificationData{CertificateID: args[0], UploadedUnitNo: user.UnitNo, UnitID: args[1], CertUpload: &certUpload}
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
	//var m, n big.Int
	//var rr, ss *big.Int
	creatorByte, _ := APIstub.GetCreator()
	certStart := bytes.IndexAny(creatorByte, "-----BEGIN")
	if certStart == -1 {
		return shim.Error("No certificate found")
	}
	certText := creatorByte[certStart:]
	//提取用户身份
	bl, _ := pem.Decode(certText)
	if bl == nil {
		return shim.Error("Could not decode the PEM structure")
	}
	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		return shim.Error("ParseCertificate failed")
	}
	uname := cert.Subject.CommonName
	//pub := cert.PublicKey.(*ecdsa.PublicKey)
	//h2 := sha256.New()
	//h2.Write([]byte(args[2]))
	//hashed := h2.Sum(nil)
	//arr := strings.Split(args[3], ",")
	//m.SetString(arr[0], 10) //大于int64的数字要用到SetString函数
	//n.SetString(arr[1], 10)
	//rr = &m
	//ss = &n
	//result := ecdsa.Verify(pub, hashed, rr, ss)
	//if result != true {
	//	return shim.Error("ECDSA Verification failed")
	//}

	userAsBytes, _ := APIstub.GetState(uname)
	user := model.User{}
	json.Unmarshal(userAsBytes, &user)

	key := fmt.Sprintf("%s,%s", args[0], args[1])
	certAsBytes, err := APIstub.GetState(key)
	if err != nil {
		return shim.Error(err.Error())
	}
	var certData = model.CertificationData{}
	json.Unmarshal(certAsBytes, &certData)

	testDataUpload := model.TestDataUpload{BaseData: args[2], EncryptedSummary: args[3], PostPersonID: user.Id, PostPersonName: user.Name}
	certData.TestDataUpload = &testDataUpload
	certAsBytes, _ = json.Marshal(certData)

	APIstub.PutState(key, certAsBytes)
	return shim.Success(nil)

}

func (s *SmartContract) trialRunDataUpload(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	//验证签名
	//var m, n big.Int
	//var rr, ss *big.Int
	creatorByte, _ := APIstub.GetCreator()
	certStart := bytes.IndexAny(creatorByte, "-----BEGIN")
	if certStart == -1 {
		return shim.Error("No certificate found")
	}
	certText := creatorByte[certStart:]
	//提取用户身份
	bl, _ := pem.Decode(certText)
	if bl == nil {
		return shim.Error("Could not decode the PEM structure")
	}
	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		return shim.Error("ParseCertificate failed")
	}
	uname := cert.Subject.CommonName
	//pub := cert.PublicKey.(*ecdsa.PublicKey)
	//h2 := sha256.New()
	//h2.Write([]byte(args[2]))
	//hashed := h2.Sum(nil)
	//arr := strings.Split(args[3], ",")
	//m.SetString(arr[0], 10) //大于int64的数字要用到SetString函数
	//n.SetString(arr[1], 10)
	//rr = &m
	//ss = &n
	//result := ecdsa.Verify(pub, hashed, rr, ss)
	//if result != true {
	//	return shim.Error("ECDSA Verification failed")
	//}
	userAsBytes, _ := APIstub.GetState(uname)
	user := model.User{}
	json.Unmarshal(userAsBytes, &user)

	key := fmt.Sprintf("%s,%s", args[0], args[1])
	certAsBytes, err := APIstub.GetState(key)
	if err != nil {
		return shim.Error(err.Error())
	}
	var certData = model.CertificationData{}
	json.Unmarshal(certAsBytes, &certData)

	trialRunDataUpload := model.TrialRunDataUpload{BaseData: args[2], EncryptedSummary: args[3], PostPersonID: user.Id, PostPersonName: user.Name}
	certData.TrialRunDataUpload = &trialRunDataUpload
	certAsBytes, _ = json.Marshal(certData)

	APIstub.PutState(key, certAsBytes)
	return shim.Success(nil)

}

func (s *SmartContract) queryUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	userAsBytes, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error("Certificate doesn't exist!")
	}
	return shim.Success(userAsBytes)
}

func (s *SmartContract) queryCert(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
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

	userAsBytes, err := APIstub.GetState(uname)
	if err != nil {
		return shim.Error("Certificate doesn't exist!")
	}
	user := model.User{}
	json.Unmarshal(userAsBytes, &user)

	key := fmt.Sprintf("%s,%s", args[0], args[1])
	cerAsBytes, _ := APIstub.GetState(key)
	cer := model.CertificationData{}
	json.Unmarshal(cerAsBytes, &cer)

	if cer.UploadedUnitNo != user.UnitNo {
		return shim.Error("Certificate not belonging to your unit")
	}
	return shim.Success(cerAsBytes)
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
	json.Unmarshal(userAsBytes, &user)

	var queryString string
	if uname == "Admin" {
		queryString = "{\"selector\":{\"CertificateID\":{\"$regex\":\"(?i)\"}}}"
	} else {
		queryString = fmt.Sprintf("{\"selector\":{\"uploadedUnitNo\":\"%s\"}}", user.UnitNo)
	}

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

func (s *SmartContract) publicQuery(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	var queryString string
	queryString = "{\"selector\":{\"CertificateID\":{\"$regex\":\"(?i)\"}}}"

	resultsIterator, err := APIstub.GetQueryResult(queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	cer := model.CertificationData{}
	pubData := model.PublicData{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		json.Unmarshal(queryResponse.Value, &cer)

		json.Unmarshal([]byte(cer.CertUpload.BaseData), &pubData)

		if pubData.CertificateID == args[0] && pubData.UnitName == args[1] && pubData.PlatformName == args[2] {
			basedataAsBytes, _ := json.Marshal(pubData)
			return shim.Success(basedataAsBytes)
		}
	}
	return shim.Error("Certificate doesn't exist!")
}

func (s *SmartContract) conditionalQuery(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
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
	user := model.User{}
	json.Unmarshal(userAsBytes, &user)

	var queryString string
	if uname == "Admin" {
		queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"id\":{\"$regex\":\"(?i)\"}},{\"%s\":\"%s\"}}}", args[0], args[1])
	} else {
		queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"uploadedUnitNo\":\"%s\"},{\"%s\":\"%s\"}}}", user.UnitNo, args[0], args[1])
	}

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
