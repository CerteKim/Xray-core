package handler

import (
	"net/http"
	"log"

	"github.com/julienschmidt/httprouter"
	"github.com/xtls/xray-core/app/web/client"
	xlog "github.com/xtls/xray-core/common/log"
)

func AddOutboundHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	res, err, m := Convert(r)
	if err != nil {
		log.Println("[Web] %v", err)
		return
	}
	xlog.Record(&xlog.AccessMessage{
                From:   "Web",
                To:     "AddOutboundHandler",
                Status: xlog.AccessAccepted,
                Detour: m["tag"].(string),
        })
	client.Client.AddOutbound(res)
}

func RemoveOutboundHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	data := ps.ByName("tag")
	xlog.Record(&xlog.AccessMessage{
		From:   "Web",
		To:     "RemoveOutboundHandler",
		Status: xlog.AccessAccepted,
		Detour: data,
	})
	client.Client.RemoveOutbound(data)
}
