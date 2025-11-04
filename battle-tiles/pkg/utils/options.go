package utils

type PageParam struct {
	NotPage  bool   `json:"not_page" form:"not_page" binding:"omitempty" default:"false"`                                         // 是否不分页 需要不分页是传递 true
	PageNo   int32  `json:"page_no" form:"page_no" binding:"required,gte=1" example:"1" default:"1" title:"页码"`                   // 页码
	PageSize int32  `json:"page_size" form:"page_size" binding:"required,gte=1,lte=1000" example:"10"  default:"20" title:"每页大小"` // 每页大小 不可大于1000
	Keyword  string `json:"keyword" form:"keyword" binding:"omitempty" title:"关键字"`                                               // 关键字
	OrderBy  string `json:"order_by" form:"order_by" binding:"omitempty" title:"排序方式"`
}

type PageResult struct {
	Header         interface{} `json:"header,omitempty" description:"头部"`                                                                              // 头部信息
	List           interface{} `json:"list" description:"列表"`                                                                                          // 结果
	ExtendedField1 interface{} `json:"extended_field1,omitempty" description:"额外字段1"`                                                                  // 扩展1
	ExtendedField2 interface{} `json:"extended_field2,omitempty" description:"额外字段2"`                                                                  // 扩展2
	NotPage        bool        `json:"not_page,omitempty" form:"not_page" binding:"omitempty" default:"false"`                                         // 是否不分页
	Total          int64       `json:"total,omitempty" description:"总条数"`                                                                              // 总条数
	PageNo         int32       `json:"page_no,omitempty" form:"page_no" binding:"omitempty" example:"1" default:"1" title:"页码"`                        // 页页码
	PageSize       int32       `json:"page_size,omitempty" form:"page_size" binding:"required,gte=1,lte=1000" example:"10"  default:"20" title:"每页大小"` // 每页大小
}

type PkByStringParam struct {
	Id string `json:"id" form:"id" binding:"required,gte=1" example:"1" default:"1" title:"主键"` // 主键
}

type PkByInt32Param struct {
	Id int32 `json:"id" form:"id" binding:"required,gte=1" example:"1" default:"1" title:"主键"` // 主键
}

type PkByInt32sParam struct {
	Id []int32 `json:"id" form:"id" binding:"required,gte=1" example:"1" default:"1" title:"主键"` // 主键
}

type Option struct {
	Value int32  `json:"value"` // 默认使用value 作为key
	Key   string `json:"key"`   // key
	Label string `json:"label"`
}

type OptionWithPy struct {
	Option
	FirstLetter string `json:"first_letter"` // 首字母
	PinyinCode  string `json:"pinyin_code"`  // 拼音码

}
