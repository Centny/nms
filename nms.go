package nms

import (
	"github.com/Centny/dbm/mgo"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"github.com/Centny/nms/nmsapi"
	"github.com/Centny/nms/nmsdb"
	"github.com/Centny/nms/nmsrc"
)

func RunNMS_S(fcfg *util.Fcfg) error {
	mgo.AddDefault2(fcfg.Val2("db_con", ""))
	err := mgo.ChkIdx(mgo.C, nmsdb.Indexes)
	if err != nil {
		return err
	}
	nmsapi.LoadAlias(fcfg)
	nmsapi.WWW = fcfg.Val2("www", ".")
	var showlog = fcfg.Val2("showlog", "0") == "1"
	netw.ShowLog = showlog
	netw.ShowLog_C = showlog
	impl.ShowLog = showlog
	nmsdb.C = mgo.C
	nmsapi.Hand("", routing.Shared)
	//
	var nms_s = nmsrc.NewNMS_S(fcfg.Val2("listen_rc", ""), fcfg.Val2("conf_dir", "."), fcfg.Val2("rc_token", ""))
	err = nms_s.L.Run()
	if err != nil {
		return err
	}
	//
	routing.Shared.Print()
	var listen = fcfg.Val("listen_web")
	log.D("listen web server on %v", listen)
	return routing.ListenAndServe(listen)
}

func RunNMS_C(fcfg *util.Fcfg) error {
	var nms_c = nmsrc.NewNMS_C(
		fcfg.Val2("rc_con", ""), fcfg.Val2("lcid", ""),
		fcfg.Val2("alias", ""), fcfg.Val2("rc_token", ""),
		fcfg.Int64ValV("delay", 12000))
	nms_c.ShowLog = true
	nms_c.R.Start()
	nms_c.R.Wait()
	return nil
}
