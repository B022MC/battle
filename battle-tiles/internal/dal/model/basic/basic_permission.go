package basic

import (
	basex "battle-tiles/pkg/plugin/gormx/base"
)

const TableNameBasicPermission = "basic_permission"

// BasicPermission 基础权限表
type BasicPermission struct {
	basex.Model[int32]

	Code        string `gorm:"column:code;type:varchar(100);not null;comment:权限编码" json:"code"`
	Name        string `gorm:"column:name;type:varchar(255);not null;comment:权限名称" json:"name"`
	Category    string `gorm:"column:category;type:varchar(50);not null;comment:权限分类" json:"category"`
	Description string `gorm:"column:description;type:text;comment:权限描述" json:"description"`
	IsDeleted   bool   `gorm:"column:is_deleted;type:bool;not null;default:false;comment:是否删除" json:"is_deleted"`
}

func (*BasicPermission) TableName() string {
	return TableNameBasicPermission
}

const TableNameBasicRolePermissionRel = "basic_role_permission_rel"
const TableNameBasicRoleMenuRel = "basic_role_menu_rel"

// BasicRolePermissionRel 角色权限关联表
type BasicRolePermissionRel struct {
	RoleID       int32 `gorm:"column:role_id;type:int4;primaryKey;comment:角色ID" json:"role_id"`
	PermissionID int32 `gorm:"column:permission_id;type:int4;primaryKey;comment:权限ID" json:"permission_id"`
}

func (*BasicRolePermissionRel) TableName() string {
	return TableNameBasicRolePermissionRel
}

const TableNameBasicMenuButton = "basic_menu_button"

// BasicMenuButton 菜单按钮表
type BasicMenuButton struct {
	basex.Model[int32]

	MenuID          int32  `gorm:"column:menu_id;type:int4;not null;comment:所属菜单ID" json:"menu_id"`
	ButtonCode      string `gorm:"column:button_code;type:varchar(100);not null;comment:按钮编码" json:"button_code"`
	ButtonName      string `gorm:"column:button_name;type:varchar(255);not null;comment:按钮名称" json:"button_name"`
	PermissionCodes string `gorm:"column:permission_codes;type:varchar(500);not null;comment:所需权限码" json:"permission_codes"`
}

func (*BasicMenuButton) TableName() string {
	return TableNameBasicMenuButton
}
