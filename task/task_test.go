package task

import (
	"fmt"
	"github.com/Centny/gwf/util"
	"github.com/Centny/nms/nmsdb"
	"testing"
	"time"
)

var mv = map[string]int{}

func echo(a *nmsdb.Action, args ...interface{}) {
	fmt.Println(util.S2Json(a))
	mv[a.Name] += 1
}
func TestTask(t *testing.T) {
	var fcfg = util.NewFcfg3()
	fcfg.InitWithFilePath("task_c.properties")
	var task *Task
	//test rc
	fcfg.SetVal("tasks", "loc_rcs,loc_rcc")
	task = NewTask(ActionF(echo))
	task.Start(fcfg)
	time.Sleep(time.Second * 3)
	task.rcs.Close()
	time.Sleep(time.Second * 3)
	task.Stop()
	if mv["loc_rcc"] < 1 {
		t.Error("error")
		return
	}
	if mv["loc_rcs"] < 1 {
		t.Error("error")
		return
	}
	fmt.Println("test rc done...\n\n")
	//test http
	fcfg.SetVal("tasks", "baidu,baidu_d")
	task = NewTask(ActionF(echo))
	task.Start(fcfg)
	time.Sleep(time.Second * 3)
	task.Stop()
	if mv["baidu"] < 1 {
		t.Error("error")
		return
	}
	if mv["baidu_d"] < 1 {
		t.Error("error")
		return
	}
	fmt.Println("test http done...\n\n")
	//
	fcfg.SetVal("loc_rcc/token", "sdkfsk")
	fcfg.SetVal("tasks", "loc_rcs,loc_rcc")
	task = NewTask(ActionF(echo))
	task.Start(fcfg)
	time.Sleep(time.Second)
	task.Stop()
	//
	fcfg.SetVal("loc_rcs/rc_addr", ":2kdd")
	fcfg.SetVal("tasks", "loc_rcs")
	task = NewTask(ActionF(echo))
	task.Start(fcfg)
	time.Sleep(time.Second)
	task.Stop()
	//
	fcfg.SetVal("loc_rcs/token", "")
	fcfg.SetVal("loc_rcc/token", "")
	fcfg.SetVal("baidu/url", "")
	fcfg.SetVal("baidu_d/url", "")
	fcfg.SetVal("tasks", "loc_rcs,loc_rcc,baidu,baidu_d")
	task = NewTask(ActionF(echo))
	task.Start(fcfg)
	time.Sleep(time.Second)
	task.Stop()
	//
	fcfg.SetVal("tasks", "")
	task = NewTask(ActionF(echo))
	task.Start(fcfg)
	time.Sleep(time.Second)
	task.Stop()
	//
	NewTaskCCH("name", "typ", nil).OnCmd(nil)
}
