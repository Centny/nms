package task

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/netw/rc"
	"github.com/Centny/gwf/util"
	"github.com/Centny/nms/nmsdb"
	"os"
	"path/filepath"
)

type Tester struct {
	ShowLog bool
}

func NewTester() *Tester {
	return &Tester{}
}

func (t *Tester) Hand(l *rc.RC_Listener_m) {
	l.AddHFunc("tester/echo", t.EchoH)
}

func (t *Tester) EchoH(rc *impl.RCM_Cmd) (interface{}, error) {
	var msg = rc.StrVal("msg")
	log.D("Tester receive message(%v) from %v", msg, rc.RemoteAddr().String())
	return "OK", nil
}

func (t *Tester) EchoSrv(r *rc.RC_Runner_m, msg string) *nmsdb.Action {
	var rbeg = util.Now()
	var _, err = r.Exec_s("tester/echo", util.Map{
		"msg": msg,
	})
	var used = util.Now() - rbeg
	if used < 1 {
		used = 1
	}
	if err == nil {
		return &nmsdb.Action{
			Type: nmsdb.AT_RCC,
			Code: 0,
			Used: 1,
			Len:  len(msg),
			Time: util.Now(),
		}
	} else {
		return &nmsdb.Action{
			Type: nmsdb.AT_RCC,
			Code: 1,
			Used: 0,
			Len:  len(msg),
			Err:  err.Error(),
			Time: util.Now(),
		}
	}
}

func (t *Tester) RunHttp(url string) *nmsdb.Action {
	var run_s = &nmsdb.Action{
		Type:  nmsdb.AT_HTTP,
		Attrs: util.Map{},
		Time:  util.Now(),
	}
	rbeg := util.Now()
	status_code, res, err := util.HGet3("%v", url)
	rused := util.Now() - rbeg
	if rused < 1 {
		rused = 1
	}
	if err == nil && status_code == 200 {
		run_s.Code, run_s.Used, run_s.Len = 0, rused, len(res)
		if t.ShowLog {
			log.D("TestHttp do get by url(%v) success by used(%v)", url, rused)
		}
	} else if err != nil {
		run_s.Code, run_s.Used, run_s.Err = 1, rused, err.Error()
		log.E("TestHttp do get by url(%v) fail with error(%v)", url, err)
	} else {
		run_s.Code, run_s.Used = status_code, rused
		if t.ShowLog {
			log.D("TestHttp do get by url(%v) fail with error(%v)", url, err)
		}
	}
	return run_s
}

func (t *Tester) RunDown(url string) *nmsdb.Action {
	var tmpf = filepath.Join(os.TempDir(), util.UUID()+".tmp")
	defer os.Remove(tmpf)
	var rbeg = util.Now()
	var dlen, err = util.DLoadV(tmpf, "%v", url)
	var used = util.Now() - rbeg
	if used < 1 {
		used = 1
	}
	var res = &nmsdb.Action{
		Type:  nmsdb.AT_DOWN,
		Attrs: util.Map{},
		Time:  util.Now(),
	}
	if err == nil {
		res.Code, res.Len, res.Used, res.Attrs["speed"] = 0, int(dlen), used, dlen/used
		if t.ShowLog {
			log.D("TestHttp do down by url(%v) success by used(%v)", url, used)
		}
	} else {
		res.Code, res.Used, res.Err = 1, used, err.Error()
		log.E("TestHttp do down by url(%v) fail with error(%v)", url, err)
	}
	return res
}
