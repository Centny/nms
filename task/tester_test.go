package task

import (
	"fmt"
	"github.com/Centny/gwf/netw/rc/rctest"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/routing/httptest"
	"testing"
	"time"
)

func TestTester(t *testing.T) {
	var tc = NewTester()
	tc.ShowLog = true
	var rc_ts = rctest.NewRCTest_j2(":23434")
	tc.Hand(rc_ts.L)
	var res = tc.EchoSrv(rc_ts.R, "abc")
	if res.Code == 0 {
		t.Error("error")
		return
	}
	time.Sleep(time.Second)
	res = tc.EchoSrv(rc_ts.R, "abc")
	if res.Code != 0 {
		t.Error(res.Err)
		return
	}

	var http = httptest.NewMuxServer()
	http.Mux.HFunc("^/kk.*$", func(hs *routing.HTTPSession) routing.HResult {
		return hs.MsgRes("OK")
	})
	res = tc.RunHttp(http.URL + "/kk")
	if res.Code != 0 {
		t.Error("error")
		return
	}
	res = tc.RunHttp(http.URL + "/ll")
	if res.Code == 0 {
		t.Error("error")
		return
	}
	res = tc.RunHttp(http.URL + "2")
	if res.Code == 0 {
		t.Error("error")
		return
	}
	res = tc.RunDown(http.URL + "/kk")
	if res.Code != 0 {
		t.Error("error")
		return
	}
	res = tc.RunDown(http.URL + "2")
	if res.Code == 0 {
		t.Error("error")
		return
	}
	fmt.Println(res)
}
