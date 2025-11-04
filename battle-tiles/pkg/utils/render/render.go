package render

import (
	"github.com/gin-gonic/gin/render"
	jsoniter "github.com/json-iterator/go"
	"github.com/liamylian/jsontime/v3"
	"net/http"
	"time"
)

var (
	_               render.Render = (*JSON)(nil)
	json                          = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonContentType               = []string{"application/json; charset=utf-8"}
)

func init() {
	timeExtension := jsontime.NewCustomTimeExtension()
	timeExtension.SetDefaultTimeFormat(time.DateTime, time.Local)
	json.RegisterExtension(timeExtension)
}

type JSON struct {
	Data any
}

func (r JSON) Render(w http.ResponseWriter) error {
	return WriteJSON(w, r.Data)
}

func (r JSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

// WriteJSON marshals the given interface object and writes it with custom ContentType.
func WriteJSON(w http.ResponseWriter, obj any) error {
	writeContentType(w, jsonContentType)
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonBytes)
	return err
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
