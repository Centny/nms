package nmsrc

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/netw/rc"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"github.com/Centny/nms/nmsdb"
	"gopkg.in/mgo.v2/bson"
	"path/filepath"
	"sync"
)

type NMS_S struct {
	L     *rc.RC_Listener_m
	CDir  string
	n_lck sync.RWMutex
}

func NewNMS_S(port, cdir, token string) *NMS_S {
	var nms = &NMS_S{}
	var l = rc.NewRC_Listener_m_j(pool.BP, port, netw.NewDoNotH())
	l.Name = "NMS_S"
	l.AddToken3(token, 1)
	l.LCH = nms
	nms.L, nms.CDir = l, cdir
	nms.Hand(l)
	return nms
}

func (n *NMS_S) Hand(l *rc.RC_Listener_m) {
	l.AddHFunc("nms/conf", n.ConfH)
	l.AddHFunc("nms/record", n.RecordH)
}

func (n *NMS_S) OnLogin(rc *impl.RCM_Cmd, token string) (string, error) {
	var lcid, alias string
	var err = rc.ValidF(`
		lcid,R|S,L:0;
		alias,R|S,L:0;
		`, &lcid, &alias)
	if err != nil {
		return "", err
	}
	err = nmsdb.FOI_Node(&nmsdb.Node{
		Id:    lcid,
		Alias: alias,
	})
	if err != nil {
		return "", err
	}
	rc.Kvs().SetVal("lcid", lcid)
	rc.Kvs().SetVal("alias", alias)
	log.D("NMS_S login by lcid(%v),alias(%v) success", lcid, alias)
	return n.L.RCH.OnLogin(rc, token)
}

func (n *NMS_S) ConfH(rc *impl.RCM_Cmd) (interface{}, error) {
	var cid = rc.Kvs().StrVal("lcid")
	var cf = filepath.Join(n.CDir, cid+".properties")
	var def = filepath.Join(n.CDir, "default.properties")
	var bys, err = util.FRead(cf, def)
	if err != nil {
		err = util.Err("NMS_S read configure file for client(%v) fail with error(%v)", cid, err)
		log.E("%v", err)
	}
	return string(bys), err
}

func (n *NMS_S) RecordH(rc *impl.RCM_Cmd) (interface{}, error) {
	var cid = rc.Kvs().StrVal("lcid")
	var data = rc.StrVal("data")
	if len(data) < 1 {
		return nil, util.Err("the data arguments is empty")
	}
	var rs = []*nmsdb.Action{}
	var err = util.Json2Ss(data, &rs)
	if err != nil {
		return nil, err
	}
	for _, r := range rs {
		r.Id = bson.NewObjectId().Hex()
		r.Nid = cid
	}
	err = nmsdb.AddAction(rs...)
	if err == nil {
		return "OK", nil
	} else {
		return "", err
	}
}
