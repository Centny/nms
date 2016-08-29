package nmsdb

import (
	"github.com/Centny/dbm/mgo"
	"testing"
)

func init() {
	func() {
		defer func() {
			recover()
		}()
		C("ss")
	}()
	mgo.AddDefault2("cny:123@loc.w:27017/cny")
	C = mgo.C
	C(CN_NODE).RemoveAll(nil)
	C(CN_ACTION).RemoveAll(nil)
}

func TestNode(t *testing.T) {
	var err = FOI_Node(&Node{
		Id:    "xx",
		Alias: "xx2",
	})
	if err != nil {
		t.Error(err)
		return
	}
	ns, err := ListNode_m()
	if err != nil {
		t.Error(err)
		return
	}
	if len(ns) != 1 {
		t.Error("error")
		return
	}
}

func TestAction(t *testing.T) {
	var err = AddAction(&Action{Nid: "n0"})
	if err != nil {
		t.Error(err)
		return
	}
	as, err := ListNodeAction("n0")
	if err != nil {
		t.Error(err)
		return
	}
	if len(as) != 1 {
		t.Error("error")
		return
	}
}

func TestCount(t *testing.T) {

}
