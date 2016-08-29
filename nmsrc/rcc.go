package nmsrc

import (
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/rc"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/tools/timer"
	"github.com/Centny/gwf/util"
	"github.com/Centny/nms/nmsdb"
	"github.com/Centny/nms/nmstask"
	"sync"
)

type NMS_C struct {
	R       *rc.RC_Runner_m
	Lcid    string
	Alias   string
	Token   string
	ShowLog bool

	//
	task  *nmstask.Task
	cache []*nmsdb.Action
	ready bool
	c_lck sync.RWMutex
	idc   uint64
}

func NewNMS_C(addr, lcid, alias, token string, delay int64) *NMS_C {
	var c = &NMS_C{
		Lcid:  lcid,
		Alias: alias,
		Token: token,
	}
	c.R = rc.NewRC_Runner_m_j(pool.BP, addr, c)
	c.R.Name = "NMS_C"
	timer.Register4(delay, c.on_time, true)
	return c
}

func (n *NMS_C) OnConn(c netw.Con) bool {
	go n.onlogin(c)
	return true
}

func (n *NMS_C) OnClose(c netw.Con) {
	n.c_lck.Lock()
	n.ready = false
	n.c_lck.Unlock()
}

func (n *NMS_C) OnCmd(c netw.Cmd) int {
	return 0
}

func (n *NMS_C) onlogin(c netw.Con) {
	var err = n.R.LoginV(n.Token, util.Map{
		"lcid":  n.Lcid,
		"alias": n.Alias,
	})
	if err != nil {
		log.E("NMS_C login by token(%v) fail with error (%v)", n.Token, err)
		n.R.Stop()
		return
	}
	if n.task != nil {
		n.task.Stop()
		n.task = nil
	}
	n.StartTask()
	n.c_lck.Lock()
	n.ready = true
	n.c_lck.Unlock()
}

func (n *NMS_C) StartTask() error {
	var cfg, err = n.R.VExec_s("nms/conf", util.Map{})
	if err != nil {
		log.E("NMS_C load conf fail with error->%v", err)
		return err
	}
	var fcfg = util.NewFcfg3()
	err = fcfg.InitWithData(cfg)
	if err != nil {
		log.E("NMS_C parse conf fail with error->%v ->conf:\n%v", err, cfg)
		return err
	}
	n.task = nmstask.NewTask(n)
	n.task.ShowLog = fcfg.Val2("showlog", "0") == "1"
	err = n.task.Start(fcfg)
	if err != nil {
		log.E("NMS_C start task fail with error->%v ->conf:\n%v", err, cfg)
	}
	return err
}

func (n *NMS_C) DoPush() {
	n.c_lck.Lock()
	if len(n.cache) < 1 || !n.ready {
		n.c_lck.Unlock()
		return
	}
	var cache = n.cache[0:]
	var clen = len(cache)
	n.c_lck.Unlock()
	var _, err = n.R.VExec_s("nms/record", util.Map{
		"data": util.S2Json(cache),
	})
	if err == nil {
		n.c_lck.Lock()
		n.cache = n.cache[clen:]
		n.c_lck.Unlock()
		if n.ShowLog {
			log.D("NMS_C do push %v record success->data:\n%v", clen, util.S2Json(cache))
		} else {
			log.D("NMS_C do push %v record success", clen)
		}
	} else {
		log.E("NMS_C do push fail with error(%v)", err)
	}
}

func (n *NMS_C) on_time(i uint64) error {
	n.DoPush()
	return nil
}

func (n *NMS_C) OnAction(a *nmsdb.Action, args ...interface{}) {
	n.c_lck.Lock()
	n.idc += 1
	a.Id = fmt.Sprintf("a%v", n.idc)
	n.cache = append(n.cache, a)
	n.c_lck.Unlock()
}
