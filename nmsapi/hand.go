package nmsapi

import (
	"github.com/Centny/gwf/routing"
)

func Hand(pre string, mux *routing.SessionMux) {
	mux.HFunc("^/(index.html)?(\\?.*)?$", IndexHtml)
	mux.HFunc("^/task.html(\\?.*)?$", TaskHtml)
}
