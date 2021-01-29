package model

import (
	"time"

	"github.com/saltbo/gopkg/strutil"
)

const (
	DirTypeSys = iota + 1
	DirTypeUser
	DirFileMaxNum = 65534
)

const (
	AclPublic    = "public"
	AclProtected = "protected"
)

// 文件或者目录
type Matter struct {
	Id         int64      `json:"id"`                                     // id 主键
	Uid        int64      `json:"uid" gorm:"not null"`                    // 用户 id
	Sid        int64      `json:"sid" gorm:"not null"`                    // storage_id  storeage 的主键
	Alias      string     `json:"alias" gorm:"size:16;not null"`          // 别名 slug
	Name       string     `json:"name" gorm:"not null"`                   // 文件名 目录名
	Type       string     `json:"type" gorm:"not null"`                   // application/pdf 等
	Size       int64      `json:"size" gorm:"not null"`                   // 文件大小，单位 bit
	DirType    int8       `json:"dirtype" gorm:"column:dirtype;not null"` // 目录类型
	Parent     string     `json:"parent" gorm:"not null"`                 // 父 Matter 的 name
	Object     string     `json:"object" gorm:"not null"`                 //  fmt.Sprintf("%s/%s%s", prefix, m.Alias, filepath.Ext(p.Name))  prefix = 日期20101010
	ACL        string     `json:"acl" gorm:"not null"`                    // 权限类型（public、protected）
	URL        string     `json:"url" gorm:"-"`                           // 文件 url
	CreatedAt  time.Time  `json:"created" gorm:"not null"`                // 创建时间
	UpdatedAt  time.Time  `json:"updated" gorm:"not null"`                // 更新时间
	UploadedAt *time.Time `json:"uploaded"`                               // 上传时间
	DeletedAt  *time.Time `json:"-"`                                      // 删除时间
	TrashedBy  string     `json:"-" gorm:"size:16;not null"`              // 谁放到回收站的
}

func NewMatter(uid, sid int64, name string) *Matter {
	return &Matter{
		Uid:   uid,
		Sid:   sid,
		Alias: strutil.RandomText(16),
		Name:  name,
		ACL:   AclProtected,
	}
}

func (Matter) TableName() string {
	return "zp_matter"
}

func (m *Matter) Clone() *Matter {
	clone := *m
	clone.Id = 0
	clone.Alias = strutil.RandomText(16)
	return &clone
}

func (m *Matter) FullPath() string {
	fp := m.Parent + m.Name
	if m.IsDir() {
		fp += "/"
	}

	return fp
}

func (m *Matter) IsDir() bool {
	return m.DirType > 0
}

func (m *Matter) Public() bool {
	return m.ACL == AclPublic
}

func (m *Matter) UserAccessible(uid int64) bool {
	return m.Uid == uid
}

func (m *Matter) SetURL(fc func(object string) string) {
	if m.Public() {
		m.URL = fc(m.Object)
	}
}
