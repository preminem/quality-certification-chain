package model

//
type User struct {
	UnitNo             string               `json:"make"` //"机构编号（营业执照号）",
	UnitType           string               `json:"make"` //机构类型（1-认证机构、2-检测检验机构、3-政府监管机构）",
	UserType           string               `json:"make"` //用户类型（1-业务员、2-审核员）",
	Id                 string               `json:"make"` //身份证号码",
	Name               string               `json:"make"` //姓名",
	State              string               `json:"make"` //状态（0-未审核、1-已审核、2-注销）
	AuditorID          string               `json:"make"` //审核员身份证号码",
	AuditorName        string               `json:"make"` //审核员姓名",
	CertApplications   []*CertApplication   `json:"certApplications"`
	DocAudits          []*DocAudit          `json:"docAudits"`
	OnsiteAudits       []*OnsiteAudit       `json:"onsiteAudits"`
	CertificationDatas []*CertificationData `json:"certificationDatas"`
}

type CertApplication struct {
	BaseData            string `json:"baseData"`
	ApplyScanHASH       string `json:"applyScanHASH"`
	LegalPersonScanHASH string `json:"legalPersonScanHASH"`
	Summary             string `json:"summary"`
	EncryptedSummary    string `json:"encryptedSummary"`
	PostPersonID        string `json:"postPersonID"`
	PostPersonName      string `json:"postPersonName"`
}

type DocAudit struct {
	BaseData         string `json:"baseData"`
	EncryptedSummary string `json:"encryptedSummary"`
	PostPersonID     string `json:"postPersonID"`
	PostPersonName   string `json:"postPersonName"`
}

type OnsiteAudit struct {
	BaseData         string `json:"baseData"`
	EncryptedSummary string `json:"encryptedSummary"`
	PostPersonID     string `json:"postPersonID"`
	PostPersonName   string `json:"postPersonName"`
}

type CertificationData struct {
	CertUpload         *CertUpload         `json:"certUpload"`
	TestDataUpload     *TestDataUpload     `json:"testDataUpload"`
	TrialRunDataUpload *TrialRunDataUpload `json:"trialRunDataUpload"`
}

type CertUpload struct {
	BaseData         string `json:"baseData"`
	CertificateID    string `json:"certificateID"`
	UnitID           string `json:"unitID"`
	EncryptedSummary string `json:"encryptedSummary"`
	PostPersonID     string `json:"postPersonID"`
	PostPersonName   string `json:"postPersonName"`
}

type TestDataUpload struct {
	BaseData         string `json:"baseData"`
	CertificateID    string `json:"certificateID"`
	UnitID           string `json:"unitID"`
	EncryptedSummary string `json:"encryptedSummary"`
	PostPersonID     string `json:"postPersonID"`
	PostPersonName   string `json:"postPersonName"`
}

type TrialRunDataUpload struct {
	BaseData         string `json:"baseData"`
	CertificateID    string `json:"certificateID"`
	UnitID           string `json:"unitID"`
	EncryptedSummary string `json:"encryptedSummary"`
	PostPersonID     string `json:"postPersonID"`
	PostPersonName   string `json:"postPersonName"`
}
