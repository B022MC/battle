package basic

const TableNameBasicUserRoleRel = "basic_user_role_rel"

type BasicUserRoleRel struct {
	UserID int32 `gorm:"column:user_id;primaryKey"`
	RoleID int32 `gorm:"column:role_id;primaryKey"`
}

func (*BasicUserRoleRel) TableName() string { return TableNameBasicUserRoleRel }
