package handler

import (
	"net/http"
	"log"

	"github.com/julienschmidt/httprouter"
	"github.com/xtls/xray-core/app/web/client"
	xlog "github.com/xtls/xray-core/common/log"
)

//"Content-Type: application/json"
func AddInboundHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	res, err, m := Convert(r)
	if err != nil {
		log.Println("[Web] %v", err)
		return
	}
	xlog.Record(&xlog.AccessMessage{
                From:   "Web",
                To:     "AddInboundHandler",
                Status: xlog.AccessAccepted,
                Detour: m["tag"].(string),
        })
	client.Client.AddInbound(res)
}

func RemoveInboundHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	data := ps.ByName("tag")
	xlog.Record(&xlog.AccessMessage{
		From:   "Web",
		To:     "RemoveInboundHandler",
		Status: xlog.AccessAccepted,
		Detour: data,
	})
	client.Client.RemoveInbound(data)
}
