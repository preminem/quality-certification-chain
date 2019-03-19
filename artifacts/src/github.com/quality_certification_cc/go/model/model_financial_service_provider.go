package model

//
type User struct {
	UnitNo      string `json:"make"` //"机构编号（营业执照号）",
	UnitType    string `json:"make"` //机构类型（1-认证机构、2-检测检验机构、3-政府监管机构）",
	userType    string `json:"make"` //用户类型（1-业务员、2-审核员）",
	Id          string `json:"make"` //身份证号码",
	Name        string `json:"make"` //姓名",
	Passwd      string `json:"make"` //密码",
	State       string `json:"make"` //状态（0-未审核、1-已审核、2-注销）
	AuditorID   string `json:"make"` //审核员身份证号码",
	AuditorName string `json:"make"` //审核员姓名",
}
