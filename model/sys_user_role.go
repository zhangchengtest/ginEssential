package model

type SysUserRole struct {
	Id     int    `json:"id"`
	UserId string `json:"userId"`
	RoleId int    `json:"roleId"`
}

// 自定义表名
func (SysUserRole) TableName() string {
	return "sys_user_role"
}
