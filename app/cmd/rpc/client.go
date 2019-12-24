package rpc

import (
	"encoding/hex"
	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"net/http"
)

func Dispatch(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !cors(&w, r) {
		return
	}
	d := types.SessionHeader{}
	if err := PopModel(w, r, ps, d); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.QueryDispatch(d)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	j, er := app.GetCodec().MarshalJSON(res)
	if er != nil {
		WriteErrorResponse(w, 400, er.Error())
		return
	}
	WriteResponse(w, string(j), r.URL.Path, r.Host)
}

func Relay(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var relay = types.Relay{}
	if !cors(&w, r) {
		return
	}
	if err := PopModel(w, r, ps, relay); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.QueryRelay(relay)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	j, er := app.GetCodec().MarshalJSON(res)
	if er != nil {
		WriteErrorResponse(w, 400, er.Error())
		return
	}
	WriteResponse(w, string(j), r.URL.Path, r.Host)
}

type sendRawTxParams struct {
	Addr        string `json:"address"`
	RawHexBytes string `json:"raw_hex_bytes"`
}

func SendRawTx(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = sendRawTxParams{}
	if !cors(&w, r) {
		return
	}
	if err := PopModel(w, r, ps, params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	bz, err := hex.DecodeString(params.RawHexBytes)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.SendRawTx(params.Addr, bz)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	j, er := app.GetCodec().MarshalJSON(res)
	if er != nil {
		WriteErrorResponse(w, 400, er.Error())
		return
	}
	WriteResponse(w, string(j), r.URL.Path, r.Host)
}
