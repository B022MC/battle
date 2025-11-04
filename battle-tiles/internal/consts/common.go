package consts

const (
	CloudPlatformDB   = "cloud-cloud"
	SmsLoginKeyPrefix = "login:sms:code:"
	IS_EXIST_USER_SQL = `select exists(select 1 from base_user where username like '%s')`
)

func SmsLoginKey(phone string) string {
	return SmsLoginKeyPrefix + phone
}
