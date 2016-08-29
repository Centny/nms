package nmsapi

import (
	"fmt"
	"github.com/Centny/gwf/routing/httptest"
	"github.com/Centny/gwf/util"
	"github.com/Centny/nms/nmsdb"
	"github.com/Centny/nms/nmsrc"
	_ "github.com/Centny/nms/test"
	"testing"
	"time"
)

func init() {
	var fcfg = util.NewFcfg3()
	fcfg.InitWithFilePath("../nms_s.properties")
	LoadAlias(fcfg)
}

func TestApi(t *testing.T) {
	//create data
	var nms_s = nmsrc.NewNMS_S(":8323", "../nmstask", "abc")
	var err = nms_s.L.Run()
	if err != nil {
		t.Error(err.Error())
		return
	}
	var nms_c = nmsrc.NewNMS_C(":8323", "task_c", "xxx", "abc", 1000)
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
	WWW = "../www"
	var ts = httptest.NewMuxServer()
	ts.Mux.ShowLog = true
	Hand("", ts.Mux)
	fmt.Println(ts.G(""))
	fmt.Printf(ts.G("/task.html?uri=%v", "http://www.baidu.com"))
}
