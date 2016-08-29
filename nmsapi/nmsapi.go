package nmsapi

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"github.com/Centny/nms/nmsdb"
	"html/template"
	"path/filepath"
	"strings"
)

var WWW = ""
var Alias = util.Map{}

func LoadAlias(fcfg *util.Fcfg) {
	for key, _ := range fcfg.Map {
		if !strings.HasPrefix(key, "alias/uri_") {
			continue
		}
		var id = strings.TrimPrefix(key, "alias/uri_")
		Alias[fcfg.Val2(key, "")] = fcfg.Val2("alias/alias_"+id, "")
	}
	log.D("load %v alias success, alias->\n%v", len(Alias), util.S2Json(Alias))
}

func IndexHtml(hs *routing.HTTPSession) routing.HResult {
	var beg int64
	var err = hs.ValidCheckVal(`
		beg,O|I,R:0;
		`, &beg)
	if err != nil {
		return hs.Printf("%v", err)
	}
	var data = util.Map{}
	nodes, err := nmsdb.ListNode_m()
	if err != nil {
		log.E("IndexHtml list node fail with error(%v)", err)
		return hs.Printf("%v", err)
	}
	if len(nodes) > 0 {
		avg, err := nmsdb.CountActionAvg(beg)
		if err != nil {
			log.E("IndexHtml count action avg fail with error(%v)", err)
			return hs.Printf("%v", err)
		}
		sub, err := nmsdb.CountActionSub(beg)
		if err != nil {
			log.E("IndexHtml count action sub fail with error(%v)", err)
			return hs.Printf("%v", err)
		}
		data = nmsdb.JoinAvgSub(avg, sub)
	}
	tmpl, err := template.New("index.html").Funcs(index_fm).ParseFiles(filepath.Join(WWW, "index.html"))
	if err != nil {
		log.E("IndexHtml parse template file(%v) fail with error(%v)", "index.html", err)
		return hs.Printf("%v", err)
	}
	err = tmpl.Execute(hs.W, util.Map{
		"alias": Alias,
		"nodes": nodes,
		"beg":   beg,
		"data":  data,
	})
	if err == nil {
		return routing.HRES_RETURN
	} else {
		return hs.Printf("%v", err)
	}
}

var index_fm = template.FuncMap{
	"strval": func(v util.Map, p string) string {
		var val = v.StrVal(p)
		if len(val) < 1 {
			val = "-"
		}
		return val
	},
	"strvalp": func(v util.Map, p ...string) string {
		var val = v.StrValP(strings.Join(p, "/"))
		if len(val) < 1 {
			val = "-"
		}
		return val
	},
	"stime": func(t int64) string {
		return util.Time(t).Format("2006-01-02 15:04:05")
	},
	"node_alias": func(nodes map[string]*nmsdb.Node, nid string) string {
		if node, ok := nodes[nid]; ok {
			return node.Alias
		} else {
			return "-"
		}
	},
	"sjson": util.S2Json,
}

func TaskHtml(hs *routing.HTTPSession) routing.HResult {
	var act = &nmsdb.Action{}
	var err = hs.ValidCheckVal(`
		nid,O|S,L:0;
		uri,R|S,L:0;
		sub,O|S,L:0;
		code,O|I,R:-999;
		used,O|I,R:0;
		beg,O|I,R:0;
		`, &act.Nid, &act.Uri, &act.Sub, &act.Code, &act.Used, &act.Time)
	if err != nil {
		return hs.Printf("%v", err)
	}
	nodes, err := nmsdb.ListNode_m()
	if err != nil {
		log.E("TaskHtml list node fail with error(%v)", err)
		return hs.Printf("%v", err)
	}
	actions, err := nmsdb.ListAction(act)
	if err != nil {
		log.E("TaskHtml list action by action(%v) fail with error(%v)", util.S2Json(act), err)
		return hs.Printf("%v", err)
	}
	tmpl, err := template.New("task.html").Funcs(index_fm).ParseFiles(filepath.Join(WWW, "task.html"))
	if err != nil {
		log.E("TaskHtml parse template file(%v) fail with error(%v)", "task.html", err)
		return hs.Printf("%v", err)
	}
	err = tmpl.Execute(hs.W, util.Map{
		"nodes":   nodes,
		"actions": actions,
		"beg":     act.Time,
	})
	if err == nil {
		return routing.HRES_RETURN
	} else {
		return hs.Printf("%v", err)
	}
}
