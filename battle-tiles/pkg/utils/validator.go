package utils

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
	transMp  = map[string]string{}
)

func InitValidator() {
	//注册翻译器
	zh := zh.New()
	uni = ut.New(zh, zh)

	trans, _ = uni.GetTranslator("zh")

	//获取gin的校验器
	validate = binding.Validator.Engine().(*validator.Validate)

	_ = validate.RegisterTranslation("ltefield", trans, func(ut ut.Translator) error {
		return ut.Add("ltefield", "{0}必须小于或等于{1}", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		param := fe.Param()
		if transMp[param] != "" {
			param = transMp[param]
		}
		t, _ := ut.T("ltefield", fe.Field(), param)
		return t
	})
	_ = validate.RegisterTranslation("eqfield", trans, func(ut ut.Translator) error {
		return ut.Add("eqfield", "{0}必须等于{1}", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		param := fe.Param()
		if transMp[param] != "" {
			param = transMp[param]
		}
		t, _ := ut.T("eqfield", fe.Field(), param)
		return t
	})

	// 注册 gtefield 翻译
	_ = validate.RegisterTranslation("gtefield", trans, func(ut ut.Translator) error {
		return ut.Add("gtefield", "{0} 必须大于或等于 最小值", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		param := fe.Param()
		if transMp[param] != "" {
			param = transMp[param]
		}
		t, _ := ut.T("gtefield", fe.Field(), param)
		return t
	})

	//注册翻译器
	zh_translations.RegisterDefaultTranslations(validate, trans)
	_ = validate.RegisterTranslation("mobile", trans, func(ut ut.Translator) error {
		return ut.Add("mobile", "{0}必须为手机号码格式!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("mobile", fe.Field())
		return t
	})

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return fld.Tag.Get("title")
	})

	// 注册自定义的校验器和翻译
	_ = validate.RegisterTranslation("float64Decimal", trans, func(ut ut.Translator) error {
		return ut.Add("float64Decimal", "{0}的小数位数超过了允许的最大位数 {1} 位!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		params := fe.Param()
		t, _ := ut.T("float64Decimal", fe.Field(), params)
		return t
	})

	_ = validate.RegisterTranslation("excludesall", trans, func(ut ut.Translator) error {
		return ut.Add("excludesall", "{0} 不能包含字符 “{1}” ", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		params := fe.Param()
		t, _ := ut.T("excludesall", fe.Field(), params)
		return t
	})

	_ = validate.RegisterTranslation("positiveNumber", trans, func(ut ut.Translator) error {
		return ut.Add("positiveNumber", "{0} 不能为负数 ", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		params := fe.Param()
		t, _ := ut.T("positiveNumber", fe.Field(), params)
		return t
	})

	validate.RegisterValidation("mobile", validateMobile)
	validate.RegisterValidation("nonNegativeNumber", validateNonNegativeNumber)
	validate.RegisterValidation("float64Decimal", validateFloat64Decimal)
	validate.RegisterValidation("excludesall", excludesAll)
	validate.RegisterValidation("positiveNumber", positiveNumber)
	validate.RegisterValidation("transferEmptyString", TransferString)
}

// Translate 翻译错误信息
func Translate(err error) map[string][]string {

	var result = make(map[string][]string)

	errors := err.(validator.ValidationErrors)

	for _, err := range errors {

		result[err.Field()] = append(result[err.Field()], err.Translate(trans))
	}
	return result
}

/****************************** 自定义验证器 ****************************/

// validateMobile 手机号码的校验规则，用于gin的请求参数自动校验
func validateMobile(fl validator.FieldLevel) bool {
	// 内部通过反射获取mobile的值
	mobile := fl.Field().String()
	//使用正则表达式判断是否合法
	isValid, _ := regexp.MatchString(`^1\d{10}$|^(0\d{2,3}-?|\(0\d{2,3}\))?[1-9]\d{4,7}(-\d{1,8})?$`, mobile)
	return isValid
}

// validateNonNegativeNumber 验证字符串是否为非负数
func validateNonNegativeNumber(fl validator.FieldLevel) bool {
	// 正则表达式匹配非负整数和浮点数
	// 内部通过反射获取mobile的值
	nonNegativeNumber := fl.Field().String()
	nonNegativeNumber = strings.TrimSpace(nonNegativeNumber)

	// 如果字符串为空，则返回 false
	if nonNegativeNumber == "" {
		return false
	}
	// re := regexp.MustCompile(`^\d*\.?\d+$`)
	re := regexp.MustCompile(`^\d+(\.\d+)?$`)
	return re.MatchString(nonNegativeNumber)
}

// 验证float64类型的小数点位数
func validateFloat64Decimal(fl validator.FieldLevel) bool {

	// 获取字段值
	value := fl.Field().Float()

	// 获取自定义的最大小数位数
	param := fl.Param()
	maxDecimalPlaces, err := strconv.Atoi(param)
	if err != nil || maxDecimalPlaces < 0 {
		return false // 如果转换失败或者小于0，返回 false
	}

	// 使用 FormatFloat 去除尾随零
	formattedValue := strconv.FormatFloat(value, 'f', -1, 64)

	// 定义正则表达式匹配最多 maxDecimalPlaces 位小数的浮点数
	re := regexp.MustCompile(`^\d+(\.\d{0,` + strconv.Itoa(maxDecimalPlaces) + `})?$`)
	isValid := re.MatchString(formattedValue)
	return isValid
}

func GetValidator() (*validator.Validate, ut.Translator) {
	return validate, trans
}

// 自定义 excludesall 验证器, 不包含指定的字符串
func excludesAll(fl validator.FieldLevel) bool {
	// 处理特殊字符： 将 "0x2C" 替换为 ","
	param := strings.ReplaceAll(fl.Param(), "0x2C", ",")
	return !strings.ContainsAny(fl.Field().String(), param)
}

func TransferString(fl validator.FieldLevel) bool {
	// 如果字段为空字符串，将其视为有效，返回 true
	return fl.Field().String() == "" || fl.Field().String() != ""
}

// 自定义验证器方法 positiveNumber  输入的字符串需要能转化为非负数的数据
func positiveNumber(fl validator.FieldLevel) bool {
	numStr := fl.Field().String()
	if numStr == "" {
		return false
	}

	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return false
	}
	return num > 0
}
