package test

import (
	"github.com/Centny/dbm/mgo"
	"github.com/Centny/nms/nmsdb"
)

func init() {
	mgo.AddDefault2("cny:123@loc.w:27017/cny")
	nmsdb.C = mgo.C
	mgo.C(nmsdb.CN_ACTION).RemoveAll(nil)
	mgo.C(nmsdb.CN_NODE).RemoveAll(nil)
}
