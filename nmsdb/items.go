package nmsdb

import (
	"github.com/Centny/gwf/util"
)

const (
	AT_HTTP = "http"
	AT_DOWN = "down"
	AT_RCC  = "rcc"
	AT_RCS  = "rcs"
)

type Node struct {
	Id    string `bson:"_id"`
	Alias string `bson:"alias"`
	Time  int64  `bson:"time" json:"time"`
}

type Action struct {
	Id    string   `bson:"_id" json:"id"`
	Nid   string   `bson:"nid" json:"nid"`
	Name  string   `bson:"name" json:"name"`
	Sub   string   `bson:"sub" json:"sub"`
	Type  string   `bson:"type" json:"type"`
	Code  int      `bson:"code" json:"code"`
	Used  int64    `bson:"used" json:"used"`
	Len   int      `bson:"len" json:"len"`
	Err   string   `bson:"err" json:"err"`
	Attrs util.Map `bson:"attrs" json:"attrs"`
	Time  int64    `bson:"time" json:"time"`
}
