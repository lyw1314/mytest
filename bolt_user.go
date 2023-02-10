package models

import (
	"encoding/json"
	"errors"
)

type Identities struct {
	ID           string    `json:"id"`
	Token        string    `json:"token"`         // 加密后的助记词或私钥
	TokenType    int       `json:"token_type"`    // 验证方式 1:密码，2:人脸或指纹
	Name         string    `json:"name"`          // 钱包名称
	SortIndex    int       `json:"sort_index"`    // 排序字段
	AvatarImage  string    `json:"avatar_image"`  // 钱包头像
	SupportChain string    `json:"support_chain"` // 支持的币种，如果是助记词为all，其余为主链名，多个逗号分隔
	UserType     int       `json:"user_type"`     // 钱包类型：0：助记词；1：为私钥；2：观察钱包；3：硬件钱包；
	AddSource    int       `json:"add_source"`    // 钱包添加来源：0：新创建；1：导入的钱包
	IsBackup     int       `json:"is_backup"`     // 助记词是否已备份：0：未备份；1：已备份；
	Salt         string    `json:"-"`             // 盐
	HardwareInfo *Hardware `json:"hardware_info"` // 硬件钱包信息
}
type Hardware struct {
	Type           int    `json:"type"`            // 钱包类型：1：web3硬件钱包（EVM）；2：厂商钱包（BTC+EVM等）
	Channel        string `json:"channel"`         // 硬件钱包厂商，如keystone等
	KeyData        string `json:"key_data"`        // 硬件钱包key信息
	Path           string `json:"path"`            // 地址路径
	ChildrenPath   string `json:"children_path"`   // 地址派生路径
	FullPath       string `json:"full_path"`       // 派生出当前地址的完整路径
	FingerPrinter  string `json:"finger_printer"`  // 硬件指纹
	Depth          int    `json:"depth"`           // 地址路径深度
	XPub           string `json:"x_pub"`           // 公钥publicKey
	ChainCode      string `json:"chain_code"`      // 所属链编码
	SequenceNumber string `json:"sequence_number"` // 派生出当前钱包地址的序号
	Extra1         string `json:"extra_1"`         // 扩展字段1
	Extra2         string `json:"extra_2"`         // 扩展字段2
	Extra3         string `json:"extra_3"`         // 扩展字段3
}

type UserBolt struct{}

const (
	path   = "./"
	dbName = "my.db"
)

var (
	bucketName = []byte("user")
)

func (ub *UserBolt) GetList(conf map[string]interface{}) ([]Identities, error) {
	db, err := NewBoltDb(path, dbName)
	defer db.Close()
	if err != nil {
		return nil, err
	}
	list, err := db.Range(bucketName)
	if err != nil {
		return nil, err
	}
	ret := make([]Identities, 0)
	for _, v := range list {
		var tmp Identities
		err := json.Unmarshal(v["value"], &tmp)
		if err != nil {
			continue
		}
		// 钱包类型筛选
		if userType, ok := conf["user_type"]; ok {
			if userType != tmp.UserType {
				continue
			}
		}
		// 钱包添加来源筛选
		if addSource, ok := conf["add_source"]; ok {
			if addSource != tmp.AddSource {
				continue
			}
		}
		ret = append(ret, tmp)
	}
	return ret, nil
}

func (ub *UserBolt) GetInfoById(id string) (*Identities, error) {
	db, err := NewBoltDb(path, dbName)
	defer db.Close()
	if err != nil {
		return nil, err
	}
	get, err := db.Get(bucketName, []byte(id))
	if err != nil || get == nil {
		return nil, err
	}
	var ret Identities
	err = json.Unmarshal(get, &ret)
	return &ret, err
}

func (ub *UserBolt) Add(param Identities) error {
	db, err := NewBoltDb(path, dbName)
	defer db.Close()
	if err != nil {
		return err
	}
	marshal, err := json.Marshal(param)
	if err != nil {
		return err
	}
	return db.Put(bucketName, []byte(param.ID), marshal)
}

func (ub *UserBolt) Delete(id string) error {
	db, err := NewBoltDb(path, dbName)
	defer db.Close()
	if err != nil {
		return err
	}
	return db.Delete(bucketName, []byte(id))
}

func (ub *UserBolt) Update(id string, updateData map[string]interface{}) error {
	getInfo, err := ub.GetInfoById(id)
	if err != nil {
		return err
	}
	if getInfo == nil {
		return errors.New("record not found")
	}
	//token_type,name,sort_index,avatar_image,is_backup
	if tokenType, ok := updateData["token_type"]; ok {
		getInfo.TokenType = tokenType.(int)
	}

	if name, ok := updateData["name"]; ok {
		getInfo.Name = name.(string)
	}

	if sortIndex, ok := updateData["sort_index"]; ok {
		getInfo.SortIndex = sortIndex.(int)
	}

	if avatarImage, ok := updateData["avatar_image"]; ok {
		getInfo.AvatarImage = avatarImage.(string)
	}

	if isBackup, ok := updateData["is_backup"]; ok {
		getInfo.IsBackup = isBackup.(int)
	}
	return ub.Add(*getInfo)
}
