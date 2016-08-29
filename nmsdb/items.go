package nmsdb

import (
	"github.com/Centny/gwf/util"
	"gopkg.in/mgo.v2"
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
	Uri   string   `bson:"uri" json:"uri"`
	Sub   string   `bson:"sub" json:"sub"`
	Type  string   `bson:"type" json:"type"`
	Code  int      `bson:"code" json:"code"`
	Used  int64    `bson:"used" json:"used"`
	Len   int      `bson:"len" json:"len"`
	Err   string   `bson:"err" json:"err"`
	Attrs util.Map `bson:"attrs" json:"attrs"`
	Time  int64    `bson:"time" json:"time"`
}

var Indexes = map[string]map[string]mgo.Index{
	CN_ACTION: map[string]mgo.Index{
		"a_nid": mgo.Index{
			Key: []string{"nid"},
		},
		"a_uri": mgo.Index{
			Key: []string{"uri"},
		},
		"a_sub": mgo.Index{
			Key: []string{"sub"},
		},
		"a_type": mgo.Index{
			Key: []string{"type"},
		},
		"a_code": mgo.Index{
			Key: []string{"code"},
		},
		"a_used": mgo.Index{
			Key: []string{"used"},
		},
		"a_len": mgo.Index{
			Key: []string{"len"},
		},
		"a_err": mgo.Index{
			Key: []string{"err"},
		},
		"a_time": mgo.Index{
			Key: []string{"time"},
		},
	},
}
