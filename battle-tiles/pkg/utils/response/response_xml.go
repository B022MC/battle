package response

import (
	"bytes"
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"text/template"
)

type OutParam struct {
	XMLName   xml.Name  `xml:"OutParam"`
	OutResult OutResult `xml:"OutResult"`
}

type OutResult struct {
	Status int    `xml:"Status"`
	Msg    string `xml:"Msg"`
	Data   any    `xml:"Data"`
}

type TemplateData struct {
	Action string
	Body   string
}

func XmlResult(c *gin.Context, action string, status int, msg string, txt string, data any) {
	tpl, _ := template.New("resp").Parse(txt)
	xmlData, _ := xml.MarshalIndent(OutParam{
		OutResult: OutResult{
			Status: status,
			Msg:    msg,
			Data:   data,
		},
	}, "", "")

	buff := new(strings.Builder)
	xml.EscapeText(buff, xmlData)
	tempData := TemplateData{
		Action: action,
		Body:   buff.String(),
	}
	var buf bytes.Buffer

	tpl.Execute(&buf, tempData)

	// 设置 Content-Type 并返回响应
	c.Header("Content-Type", "text/xml")
	c.String(http.StatusOK, buf.String())
}

func XmlSuccess(c *gin.Context, action, txt string) {
	XmlResult(c, action, 1, "Success", txt, nil)
}

func XmlSuccessWithData(c *gin.Context, action, txt string, data any) {
	XmlResult(c, action, 1, "Success", txt, data)
}

func XmlFail(c *gin.Context, action, txt string, err error) {
	XmlResult(c, action, 0, err.Error(), txt, nil)
}

func CustomXml(c *gin.Context, action, result, txt string) {
	tpl, _ := template.New("resp").Parse(txt)
	buff := new(strings.Builder)
	xml.EscapeText(buff, []byte(result))
	data := TemplateData{
		Action: action,
		Body:   buff.String(),
	}
	var buf bytes.Buffer

	tpl.Execute(&buf, data)
	c.Header("Content-Type", "text/xml")
	c.String(http.StatusOK, buf.String())
}
