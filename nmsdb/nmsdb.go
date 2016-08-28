package nmsdb

import (
	"github.com/Centny/gwf/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var C = func(name string) *mgo.Collection {
	panic("the nmsdb is not initial")
}

const (
	CN_NODE   = "nms_node"
	CN_ACTION = "nms_action"
)

func AddAction(as ...*Action) error {
	var vs = []interface{}{}
	for _, a := range as {
		vs = append(vs, a)
	}
	return C(CN_ACTION).Insert(vs...)
}

func FOI_Node(n *Node) error {
	_, err := C(CN_NODE).Upsert(bson.M{"_id": n.Id}, bson.M{
		"$set": bson.M{
			"alias": n.Alias,
			"time":  util.Now(),
		},
	})
	return err
}

func ListNode() (ns []*Node, err error) {
	err = C(CN_NODE).Find(nil).All(&ns)
	return
}

func ListNodeAction(nid string) (as []*Action, err error) {
	err = C(CN_ACTION).Find(bson.M{"nid": nid}).All(&as)
	return
}
