package response

import (
	"battle-tiles/pkg/utils"
	"battle-tiles/pkg/utils/ecode"
	"battle-tiles/pkg/utils/render"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"

	"mime"
	"net"
	"net/http"
)

func StandardPgError(pgErr *pgconn.PgError) string {
	switch pgErr.Code {
	case "23505": // unique_violation
		return "该记录已存在, 请检查重复项!"
	case "23503": // foreign_key_violation
		return "外键约束失败"
	case "23502": // not_null_violation
		return "字段不能为空"
	case "42703":
		return "字段列不存在"
	}
	return "数据库错误"
}

type Body struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Result(c *gin.Context, code int, data any, err any) {
	if data == nil {
		data = make(map[string]string, 0)
	}
	resp := Body{
		Code: code,
		Data: data,
		Msg:  ecode.Text(code),
	}

	switch e := err.(type) {
	case validator.ValidationErrors:
		errMarshal, _ := json.Marshal(utils.Translate(e))
		errString := string(errMarshal)
		resp.Msg = resp.Msg + ": " + errString
	case *pgconn.PgError:
		resp.Msg = resp.Msg + ": " + StandardPgError(e)
	case *json.UnmarshalTypeError:
		resp.Msg = resp.Msg + ": 解析" + e.Field + "异常"
	case *net.OpError:
		resp.Msg = resp.Msg + ": 网络连接异常"
	case error:
		resp.Msg = resp.Msg + ": " + e.Error()
	case string:
		resp.Msg = resp.Msg + ":" + e
	}
	c.Render(http.StatusOK, render.JSON{Data: resp})
}

func Success(c *gin.Context, data any) {
	Result(c, ecode.Success, data, nil)
}

func SuccessWithOK(c *gin.Context) {
	Result(c, ecode.Success, "ok", nil)
}

// 向前端返回Excel文件
// 参数 content 为上面生成的io.ReadSeeker， name 为返回前端的文件名
func ResponseXls(c *gin.Context, file *excelize.File, name string) {
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", name))
	c.Writer.Header().Set("Content-Type", mime.TypeByExtension(".xlsx"))
	c.Writer.Header().Set("ifdownload", "true")
	if err := file.Write(c.Writer); err != nil {
		Result(c, ecode.Failed, "文件写入失败", err)
	}
}

func Fail(c *gin.Context, code int, err any) {
	Result(c, code, nil, err)
}
