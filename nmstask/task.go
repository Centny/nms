package nmstask

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/rc"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"github.com/Centny/nms/nmsdb"
	"strings"
	"time"
)

type ActionF func(a *nmsdb.Action, args ...interface{})

func (f ActionF) OnAction(a *nmsdb.Action, args ...interface{}) {
	f(a, args...)
}

type ActionH interface {
	OnAction(a *nmsdb.Action, args ...interface{})
}

type Task struct {
	Tester
	Running bool
	H       ActionH
	//
	rcs *rc.RC_Listener_m
	rcc *rc.RC_Runner_m
}

func NewTask(h ActionH) *Task {
	return &Task{H: h}
}

func (t *Task) Start(fcfg *util.Fcfg) error {
	var tasks = fcfg.Val2("tasks", "")
	if len(tasks) < 1 {
		return util.Err("the loc/tasks is empty")
	}
	for _, task := range strings.Split(tasks, ",") {
		go t.run_c(fcfg, task)
	}
	return nil
}

func (t *Task) run_c(fcfg *util.Fcfg, name string) {
	var err = t.Run(fcfg, name)
	if err != nil {
		log.E("Task run name(%v) fail with error %v", name, err)
	}
}

func (t *Task) Run(fcfg *util.Fcfg, name string) error {
	var typ = fcfg.Val(name + "/type")
	switch typ {
	case "http":
		var url = fcfg.Val2(name+"/url", "")
		if len(url) < 1 {
			return util.Err("the http task(%v) url is empty", name)
		}
		var delay = fcfg.Int64ValV(name+"/delay", 120000)
		t.run_http(name, url, time.Duration(delay))
	case "down":
		var url = fcfg.Val2(name+"/url", "")
		if len(url) < 1 {
			return util.Err("the down task(%v) url is empty", name)
		}
		var delay = fcfg.Int64ValV(name+"/delay", 120000)
		t.run_down(name, url, time.Duration(delay))
	case "rcc":
		var rc_con = fcfg.Val2(name+"/rc_con", "")
		var token = fcfg.Val2(name+"/token", "")
		if len(rc_con) < 1 || len(token) < 1 {
			return util.Err("the rcs task(%v) rc_con or token is empty", name)
		}
		var delay = fcfg.Int64ValV(name+"/delay", 120000)
		t.run_rcc(name, rc_con, token, time.Duration(delay))
	case "rcs":
		var addr = fcfg.Val2(name+"/rc_addr", "")
		var token = fcfg.Val2(name+"/token", "")
		if len(addr) < 1 || len(token) < 1 {
			return util.Err("the rcs task(%v) addr or token is empty", name)
		}
		t.run_rcs(name, addr, token)
	}
	return nil
}

func (t *Task) run_http(name, url string, delay time.Duration) {
	t.Running = true
	for t.Running {
		var res = t.RunHttp(url)
		res.Uri = url
		res.Sub = "test"
		t.H.OnAction(res)
		time.Sleep(time.Millisecond * delay)
	}
	log.D("Task run http for url(%v) done...", url)
}

func (t *Task) run_down(name, url string, delay time.Duration) {
	t.Running = true
	for t.Running {
		var res = t.RunDown(url)
		res.Uri = url
		res.Sub = "test"
		t.H.OnAction(res)
		time.Sleep(time.Millisecond * delay)
	}
}

func (t *Task) run_rcs(name, addr, token string) {
	var h = NewTaskCCH_S(addr, "rcs", t.H)
	t.rcs = rc.NewRC_Listener_m_j(pool.BP, addr, h)
	t.rcs.Name = name
	t.rcs.AddHFunc("tester/echo", t.EchoH)
	t.rcs.AddToken3(token, 1)
	h.T = t
	var err = t.rcs.Run()
	if err != nil {
		log.E("Task run on addr(%v),name(%v) fail with error(%v)", addr, name, err)
		return
	}
	t.rcs.Wait()
}

func (t *Task) run_rcc(name, con, token string, delay time.Duration) {
	var h = NewTaskCCH_C(token, con, "rcs", t.H)
	t.rcc = rc.NewRC_Runner_m_j(pool.BP, con, h)
	t.rcc.Name = name
	t.rcc.L.Dailer.OnDailFail = h.OnDailFail
	h.Delay, h.Runner, h.T, h.Token = delay, t.rcc, t, token
	t.rcc.Start()
	t.Running = true
	t.rcc.Wait()
}

func (t *Task) Stop() {
	t.Running = false
	if t.rcs != nil {
		t.rcs.Close()
		t.rcs = nil
	}
	if t.rcc != nil {
		t.rcc.Stop()
		t.rcc = nil
	}
}

type TaskCCH struct {
	T    *Task
	Uri  string
	Type string
	H    ActionH
}

func NewTaskCCH(uri, typ string, h ActionH) *TaskCCH {
	return &TaskCCH{
		Uri:  uri,
		Type: typ,
		H:    h,
	}
}

func (t *TaskCCH) OnConn(c netw.Con) bool {
	c.SetWait(true)
	if t.T != nil && t.T.Running {
		t.H.OnAction(&nmsdb.Action{
			Uri:  t.Uri,
			Code: 0,
			Type: t.Type,
			Sub:  "conn",
			Attrs: util.Map{
				"addr": c.RemoteAddr().String(),
			},
		}, t, c)
	}
	return true
}

func (t *TaskCCH) OnDailFail(addr string, err error) {
	if t.T != nil && t.T.Running {
		t.H.OnAction(&nmsdb.Action{
			Uri:  t.Uri,
			Code: 1,
			Type: t.Type,
			Sub:  "fail",
			Err:  err.Error(),
			Attrs: util.Map{
				"addr": addr,
			},
		}, t, err)
	}
}

//see ConHandler
func (t *TaskCCH) OnClose(c netw.Con) {
	if t.T != nil && t.T.Running {
		t.H.OnAction(&nmsdb.Action{
			Uri:  t.Uri,
			Code: 1,
			Type: t.Type,
			Sub:  "close",
			Attrs: util.Map{
				"addr": c.RemoteAddr().String(),
			},
		}, t, c)
	}
}

//see CmdHandler
func (t *TaskCCH) OnCmd(c netw.Cmd) int {
	return 0
}

type TaskCCH_S struct {
	*TaskCCH
}

func NewTaskCCH_S(uri, typ string, h ActionH) *TaskCCH_S {
	return &TaskCCH_S{
		TaskCCH: NewTaskCCH(uri, typ, h),
	}
}

type TaskCCH_C struct {
	*TaskCCH
	Delay  time.Duration
	Token  string
	Runner *rc.RC_Runner_m
	//
	running bool
}

func NewTaskCCH_C(token, uri, typ string, h ActionH) *TaskCCH_C {
	return &TaskCCH_C{
		Token:   token,
		TaskCCH: NewTaskCCH(uri, typ, h),
	}
}

func (t *TaskCCH_C) OnConn(c netw.Con) bool {
	go t.run_c(c)
	return t.TaskCCH.OnConn(c)
}

//see ConHandler
func (t *TaskCCH_C) OnClose(c netw.Con) {
	t.running = false
	t.TaskCCH.OnClose(c)
}

func (t *TaskCCH_C) run_c(c netw.Con) {
	var err = t.Runner.Login_(t.Token)
	if err != nil {
		log.E("TaskCCH_C do login by token(%v) fail with error(%v)", t.Token, err)
		time.Sleep(3 * time.Second)
		c.Close()
		return
	}
	t.running = true
	for t.running {
		var res = t.T.EchoSrv(t.Runner, "abc")
		res.Uri = t.Uri
		res.Sub = "test"
		t.H.OnAction(res)
		time.Sleep(t.Delay * time.Millisecond)
	}
	t.running = false
}
