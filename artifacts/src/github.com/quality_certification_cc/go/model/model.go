package model

//
type User struct {
	UnitNo           string             `json:"unitNo"`      //"机构编号（营业执照号）",
	UnitType         string             `json:"unitType"`    //机构类型（1-认证机构、2-检测检验机构、3-政府监管机构）",
	UserType         string             `json:"userType"`    //用户类型（1-业务员、2-审核员）",
	Id               string             `json:"id"`          //身份证号码",
	Name             string             `json:"name"`        //姓名",
	State            int                `json:"state"`       //状态（0-未审核、1-已审核、2-注销）
	AuditorID        string             `json:"auditorId"`   //审核员身份证号码",
	AuditorName      string             `json:"auditorName"` //审核员姓名",
	CertApplications []*CertApplication `json:"certApplications"`
	DocAudits        []*DocAudit        `json:"docAudits"`
	OnsiteAudits     []*OnsiteAudit     `json:"onsiteAudits"`
	//CertificationDatas []*CertificationData `json:"certificationDatas"`
}

type CertApplication struct {
	BaseData            string `json:"baseData"`
	ApplyScanHASH       string `json:"applyScanHASH"`
	LegalPersonScanHASH string `json:"legalPersonScanHASH"`
	Summary             string `json:"summary"`
	EncryptedSummary    string `json:"encryptedSummary"`
	//PostPersonID        string `json:"postPersonID"`
	//PostPersonName      string `json:"postPersonName"`
}

type DocAudit struct {
	BaseData         string `json:"baseData"`
	EncryptedSummary string `json:"encryptedSummary"`
	//PostPersonID     string `json:"postPersonID"`
	//PostPersonName   string `json:"postPersonName"`
}

type OnsiteAudit struct {
	BaseData         string `json:"baseData"`
	EncryptedSummary string `json:"encryptedSummary"`
	//PostPersonID     string `json:"postPersonID"`
	//PostPersonName   string `json:"postPersonName"`
}

type CertificationData struct {
	CertificateID      string              `json:"certificateID"`
	UnitID             string              `json:"unitID"`
	CertUpload         *CertUpload         `json:"certUpload"`
	TestDataUpload     *TestDataUpload     `json:"testDataUpload"`
	TrialRunDataUpload *TrialRunDataUpload `json:"trialRunDataUpload"`
}

type CertUpload struct {
	BaseData string `json:"baseData"`
	//CertificateID    string `json:"certificateID"`
	//UnitID           string `json:"unitID"`
	UploadedUnitNo   string `json:"uploadedUnitNo"` //上传者机构编号
	EncryptedSummary string `json:"encryptedSummary"`
	PostPersonID     string `json:"postPersonID"`
	PostPersonName   string `json:"postPersonName"`
}

type TestDataUpload struct {
	BaseData string `json:"baseData"`
	//CertificateID    string `json:"certificateID"`
	//UnitID           string `json:"unitID"`
	UploadedUnitNo   string `json:"uploadedUnitNo"` //上传者机构编号
	EncryptedSummary string `json:"encryptedSummary"`
	PostPersonID     string `json:"postPersonID"`
	PostPersonName   string `json:"postPersonName"`
}

type TrialRunDataUpload struct {
	BaseData string `json:"baseData"`
	//CertificateID    string `json:"certificateID"`
	//UnitID           string `json:"unitID"`
	UploadedUnitNo   string `json:"uploadedUnitNo"` //上传者机构编号
	EncryptedSummary string `json:"encryptedSummary"`
	PostPersonID     string `json:"postPersonID"`
	PostPersonName   string `json:"postPersonName"`
}

type PublicData struct {
	CertificateID      string `json:"certificateID"`
	UnitName           string `json:"unitName"`           //机构名称
	PlatformName       string `json:"platformName"`       //交易平台名称
	CertificationClass string `json:"certificationClass"` //认证级别
	ValidityTerm       string `json:"validityTerm"`       //有效期至
}
