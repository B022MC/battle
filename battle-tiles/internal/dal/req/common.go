package req

import (
	"time"
)

type KeyInfoReq struct {
	Keyword string `json:"keyword" form:"keyword"  mapstructure:"keyword" title:"关键字"` // 关键字
}

type PageInfoReq struct {
	Page     int `json:"page" form:"page" binding:"required,gte=1" example:"1" default:"1"`              // 页码
	PageSize int `json:"page_size" form:"page_size" binding:"required,gte=1" example:"10"  default:"20"` // 每页大小 最大1000
	KeyInfoReq
}

type ReqById struct {
	ID      int32  `json:"id" form:"id" binding:"required,gte=1" title:"唯一主键"` // 唯一主键
	Keyword string `json:"keyword" form:"keyword"`                             // 关键字
}

type ReqByIds struct {
	IDs []int32 `json:"ids" form:"ids" binding:"required" ` // 唯一主键
}

type ReqInfo struct {
	Keyword string `json:"keyword" form:"keyword"` //关键字
}

type PageInfo struct {
	Page     int    `json:"page" form:"page" binding:"required,gte=1" example:"1" default:"1"`        // 页码
	PageSize int    `json:"page_size" form:"page_size" binding:"required" example:"10"  default:"20"` // 每页大小
	Keyword  string `json:"keyword" form:"keyword"`                                                   // 关键字
}

type ReqRefreshToken struct {
	Token string `json:"token" form:"token" binding:"required"` // token
	OrgId int32  `json:"org_id,string" form:"org_id"`           // org_id 用于切换 普通用户不需要传递
}

type ReqByIdAndKeyword struct {
	ID      int32  `json:"id" form:"id" binding:"required,gte=1" title:"唯一主键"` // 唯一主键
	Keyword string `json:"keyword" form:"keyword"  title:"关键字"`                // 关键字
}

type PageInfoByPatientID struct {
	PageInfo
	PatientID int32 `json:"patient_id" form:"patient_id" binding:"required,gt=0" mapstructure:"patient_id" title:"患者"` // 患者ID
}

type FilterByTime struct {
	StartTime time.Time `json:"start_time" form:"start_time" binding:"required" title:"日期" time_format:"2006-01-02"` // 过滤開始时间  格式 2006-01-02
	EndTime   time.Time `json:"end_time" form:"end_time" binding:"required" title:"日期" time_format:"2006-01-02"`     // 过滤結束时间  格式 2006-01-02
}

type FilterAndPageByTime struct {
	FilterByTime
	PageInfo
}

// ===== 圈/群通用请求体 =====
type GroupBaseRequest struct {
	HouseGID int `json:"house_gid" binding:"required"`
}

type GroupBindRequest struct {
	HouseGID  int `json:"house_gid" binding:"required"`
	MessageID int `json:"message_id" binding:"required"`
}

type GroupForbidRequest struct {
	HouseGID  int    `json:"house_gid" binding:"required"`
	Key       string `json:"key" binding:"required"`
	MemberIDs []int  `json:"member_ids"`
}
