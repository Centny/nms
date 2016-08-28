package nmsrc

import (
	"github.com/Centny/gwf/util"
	"github.com/Centny/nms/nmsdb"
	_ "github.com/Centny/nms/test"
	"testing"
	"time"
)

func TestRc(t *testing.T) {
	var nms_s = NewNMS_S(":8323", "../task", "abc")
	var err = nms_s.L.Run()
	if err != nil {
		t.Error(err.Error())
		return
	}
	var nms_c = NewNMS_C(":8323", "task_c", "xxx", "abc", 1000)
	nms_c.ShowLog = true
	nms_c.R.Start()
	time.Sleep(3 * time.Second)
	ns, err := nmsdb.ListNode()
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(ns) < 1 || ns[0].Id != "task_c" || ns[0].Alias != "xxx" {
		t.Error("error")
		return
	}
	//
	nms_c.ShowLog = false
	nms_c.R.Close()
	time.Sleep(3 * time.Second)
	//
	nms_c.task.Stop()
	time.Sleep(3 * time.Second)
	//
	nms_c.R.VExec_s("nms/record", util.Map{})
	nms_c.R.VExec_s("login_", util.Map{"token": "abc"})
	nms_c.R.VExec_s("nms/record", util.Map{"data": "sdfs"})
	//
	nms_c.R.Stop()
	//
	nms_c = NewNMS_C(":8323", "task_x", "xxx", "abc", 1000)
	nms_c.ShowLog = true
	nms_c.R.Start()
	time.Sleep(3 * time.Second)
	nms_c.R.Stop()
	//
	nms_c = NewNMS_C(":8323", "task_c", "xxx", "abcxx", 1000)
	nms_c.ShowLog = true
	nms_c.R.Start()
	time.Sleep(3 * time.Second)
	//
	nms_c.OnCmd(nil)
}
