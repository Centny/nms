package nmsdb

import (
	"github.com/Centny/gwf/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var C = func(name string) *mgo.Collection {
	panic("the nmsdb is not initial")
}

const (
	CN_NODE   = "nms_node"
	CN_ACTION = "nms_action"
)

func AddAction(as ...*Action) error {
	var vs = []interface{}{}
	for _, a := range as {
		vs = append(vs, a)
	}
	return C(CN_ACTION).Insert(vs...)
}

func FOI_Node(n *Node) error {
	_, err := C(CN_NODE).Upsert(bson.M{"_id": n.Id}, bson.M{
		"$set": bson.M{
			"alias": n.Alias,
			"time":  util.Now(),
		},
	})
	return err
}

func ListNode() (ns []*Node, err error) {
	err = C(CN_NODE).Find(nil).All(&ns)
	return
}

func ListNode_m() (ns_m map[string]*Node, err error) {
	var ns []*Node
	ns, err = ListNode()
	if err == nil {
		ns_m = map[string]*Node{}
		for _, n := range ns {
			ns_m[n.Id] = n
		}
	}
	return
}

func ListNodeAction(nid string) (as []*Action, err error) {
	err = C(CN_ACTION).Find(bson.M{"nid": nid}).All(&as)
	return
}

func ListAction(a *Action) (as []*Action, err error) {
	var query = bson.M{}
	if len(a.Nid) > 0 {
		query["nid"] = a.Nid
	}
	if len(a.Uri) > 0 {
		query["uri"] = a.Uri
	}
	if len(a.Type) > 0 {
		query["type"] = a.Type
	}
	if len(a.Err) > 0 {
		query["code"] = bson.M{
			"$ne": 0,
		}
	}
	if a.Used > 0 {
		query["used"] = bson.M{
			"$gte": a.Used,
		}
	}
	if a.Time > 0 {
		query["time"] = bson.M{
			"$gte": a.Time,
		}
	}
	for k, v := range a.Attrs {
		query["attrs."+k] = v
	}
	err = C(CN_ACTION).Find(query).All(&as)
	return
}

func count_ms2m_sub(ms []util.Map) util.Map {
	var res = util.Map{}
	for _, m := range ms {
		delete(m, "_id")
		var node = res.MapVal(m.StrVal("uri"))
		if node == nil {
			node = util.Map{}
		}
		var sub = node.MapVal(m.StrVal("nid"))
		if sub == nil {
			sub = util.Map{}
		}
		sub[m.StrVal("sub")] = m
		node[m.StrVal("nid")] = sub
		res[m.StrVal("uri")] = node
	}
	return res
}

func CountActionAvg(beg int64) (util.Map, error) {
	var pipe = []bson.M{
		bson.M{
			"$match": bson.M{
				"time": bson.M{
					"$gte": beg,
				},
				"code": 0,
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"uri": "$uri",
					"nid": "$nid",
					"sub": "$sub",
				},
				"avg": bson.M{
					"$avg": "$used",
				},
				"min": bson.M{
					"$min": "$used",
				},
				"max": bson.M{
					"$max": "$used",
				},
				"len": bson.M{
					"$sum": 1,
				},
			},
		},
		bson.M{
			"$project": bson.M{
				"uri": "$_id.uri",
				"nid": "$_id.nid",
				"sub": "$_id.sub",
				"avg": bson.M{
					"$ceil": "$avg",
				},
				"min": "$min",
				"max": "$max",
				"len": "$len",
			},
		},
	}
	var ms = []util.Map{}
	var err = C(CN_ACTION).Pipe(pipe).All(&ms)
	return count_ms2m_sub(ms), err
}

func CountActionSub(beg int64) (util.Map, error) {
	var pipe = []bson.M{
		bson.M{
			"$match": bson.M{
				"time": bson.M{
					"$gte": beg,
				},
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"uri": "$uri",
					"nid": "$nid",
					"sub": "$sub",
				},
				"suc": bson.M{
					"$sum": bson.M{
						"$cond": bson.M{
							"if": bson.M{
								"$eq": []interface{}{"$code", 0},
							},
							"then": 1,
							"else": 0,
						},
					},
				},
				"err": bson.M{
					"$sum": bson.M{
						"$cond": bson.M{
							"if": bson.M{
								"$ne": []interface{}{"$code", 0},
							},
							"then": 1,
							"else": 0,
						},
					},
				},
			},
		},
		bson.M{
			"$project": bson.M{
				"uri": "$_id.uri",
				"nid": "$_id.nid",
				"sub": "$_id.sub",
				"suc": "$suc",
				"err": "$err",
			},
		},
	}
	var ms = []util.Map{}
	var err = C(CN_ACTION).Pipe(pipe).All(&ms)
	return count_ms2m_sub(ms), err
}

func JoinAvgSub(avg, sub util.Map) util.Map {
	var res = util.Map{}
	join_map(res, avg)
	join_map(res, sub)
	return res
}
func join_map(res, val util.Map) {
	for uri, _ := range val {
		var node_m = res.MapVal(uri)
		if node_m == nil {
			node_m = util.Map{}
		}
		var node_v = val.MapVal(uri)
		for node, _ := range node_v {
			var sub_m = node_m.MapVal(node)
			if sub_m == nil {
				sub_m = util.Map{}
			}
			var sub_v = node_v.MapVal(node)
			for sub, _ := range sub_v {
				var data_m = sub_m.MapVal(sub)
				if data_m == nil {
					data_m = util.Map{}
				}
				var data_v = sub_v.MapVal(sub)
				for k, v := range data_v {
					data_m[k] = v
				}
				sub_m[sub] = data_m
			}
			node_m[node] = sub_m
		}
		res[uri] = node_m
	}
}

// func CountActionErr(used, beg int64) (util.Map, error) {
// 	var pipe = []bson.M{
// 		bson.M{
// 			"$match": bson.M{
// 				"time": bson.M{
// 					"$gte": beg,
// 				},
// 			},
// 		},
// 		bson.M{
// 			"$group": bson.M{
// 				"_id": bson.M{
// 					"uri": "$uri",
// 					"sub": "$sub",
// 				},
// 				"suc": bson.M{
// 					"$sum": bson.M{
// 						"$cond": bson.M{
// 							"if": bson.M{
// 								"$eq": []interface{}{"$code", 0},
// 							},
// 							"then": 1,
// 							"else": 0,
// 						},
// 					},
// 				},
// 				"err": bson.M{
// 					"$sum": bson.M{
// 						"$cond": bson.M{
// 							"if": bson.M{
// 								"$ne": []interface{}{"$code", 0},
// 							},
// 							"then": 1,
// 							"else": 0,
// 						},
// 					},
// 				},
// 				"used": bson.M{
// 					"$sum": bson.M{
// 						"$cond": bson.M{
// 							"if": bson.M{
// 								"$gte": []interface{}{"$used", used},
// 							},
// 							"then": 1,
// 							"else": 0,
// 						},
// 					},
// 				},
// 			},
// 		},
// 		bson.M{
// 			"$project": bson.M{
// 				"uri":  "$_id.uri",
// 				"sub":  "$_id.sub",
// 				"nid":  "$_id.nid",
// 				"suc":  "$suc",
// 				"err":  "$err",
// 				"used": "$used",
// 			},
// 		},
// 	}
// 	var ms = []util.Map{}
// 	var err = C(CN_ACTION).Pipe(pipe).All(&ms)
// 	return count_ms2m_test(ms), err
// }
