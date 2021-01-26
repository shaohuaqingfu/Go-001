package api

import (
	"Week02/src/model"
	"Week02/src/service"
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

const (
	Address = "127.0.0.1:8080"
)

func Init() error {
	registerHandler()
	err := http.ListenAndServe(Address, nil)
	return errors.Wrap(err, "listen failed")
}

type GeekErrHandler struct {
	svc *service.GeekErrService
}

func (h *GeekErrHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("error: %+v", errors.Wrap(err, "parse form failed"))
		writeResponse(w, resp(http.StatusBadRequest, "参数解析失败", nil))
		return
	}
	id := r.PostFormValue("id")
	var user *model.User
	user, err = h.svc.GetById(id)
	if err != nil {
		log.Printf("error: %+v", err)
		writeResponse(w, resp(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	writeResponse(w, resp(http.StatusOK, "成功", user))
}

func writeResponse(w http.ResponseWriter, response *Response) {
	w.WriteHeader(response.Code)
	buf, _ := json.Marshal(response)
	_, _ = w.Write(buf)
}

type Response struct {
	Code int
	Msg  string
	Data interface{}
}

func resp(code int, msg string, data interface{}) *Response {
	r := &Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	return r
}

func registerHandler() {
	geekErrHandler := &GeekErrHandler{
		svc: &service.GeekErrService{},
	}
	http.Handle("/geekErr", geekErrHandler)
}
