package main

//---------------------------------------  数据上链结构体

// Kyc 数据结构
type Kyc struct {
	// file
	KycID     string `json:"id" pkey:""` // 主键
	KycType   string `json:"kycType"`    // 数据类型
	KycTime   string `json:"kycTime"`    // 上链时间
	KycString string `json:"kycString"`  // 内容密文
	SignKey   string `json:"signKey"`    // 签名key
	SignPower string `json:"signPower"`  // 访问权限
	SignUser  Use    `json:"signUser"`   //签名用户信息
}

//	Byc 数据结构
type Byc struct {
	BycID       string     `json:"id" pkey:""` // 用户信息
	BycUser     Use        `json:"bycUser"`    // 用户信息
	BycMes      []BycMes   `json:"kycMes"`     // 数据关联信息
	BycSign     Sign       `json:"sign"`       // 权限信息
	BycSignList []SignList `json:"signList"`   // 权限信息列表
}

// Byc 数据关联结构
type BycMes struct {
	BycID  string `json:"bycId"`  //外键
	BycKey string `json:"bycKey"` //秘钥

	ParBycID  string `json:"parBycID"`  //外键
	ParBycKey string `json:"parBycKey"` //秘钥

	BycRelation string `json:"bycRelation"` //外键关系  parent  sun
	BySol       string `json:"bySol"`       //开关
}

// Byc 权限结构
type Sign struct {
	OrgKey  string `json:"orgKey"`  // 组织秘钥
	UserKey string `json:"userKey"` // 用户秘钥
	Public  string `json:"public"`  // 公共秘钥
	Self    string `json:"self"`    // 本数据秘钥
}

// Byc 访问权限列表
type SignList struct {
	RelID   string `json:"relID"`   // 访问ID
	RelType string `json:"relType"` // 用户 || 数据
	RelTime string `json:"relTime"` // 访问最后时间
	RelNum  string `json:"relNum"`  // 访问数量
	RelSol  string `json:"relSol"`  // 是否开启
}

// Kyc||Byc|| 用户信息记录
type Use struct {
	UseName    string `json:"useName"`    // 用户名称
	UseID      string `json:"useId"`      // 用户ID
	UseOrgName string `json:"useOrgName"` // 组织名称
	UseOrgID   string `json:"useOrgId"`   // 组织ID
	UseType    string `json:"useType"`    // 用户类型
	UseCa      string `json:"useCa"`      // 用户ca名称 ||
}

//---------------------------------------  链码方法函数

//	链码||函数方法
const (
	uploadAsset          string = "uploadAsset"
	updateAsset          string = "updateAsset"
	getAssetByID         string = "getAssetByID"
	getAssetByUpdateTime string = "getAssetByUpdateTime"
	getAssetList         string = "getAssetList"
)

//	Key|| Key + Key
const (
	KYC string = "KYC"
	BYC string = "BYC"
)

//  父子 标签
const (
	Title_Parent string = "p"
	Title_Sun    string = "s"
	Title_Tong   string = "t"
)
