package model

type Admin struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	Name      string `json:"name" gorm:"column:name;type:varchar(64);"`
	Password  string `json:"password" gorm:"column:password;type:varchar(64);"`
	CreatedAt string `json:"created_at" gorm:"column:created_at;type:datetime;"`
	UpdatedAt string `json:"updated_at" gorm:"column:updated_at;type:datetime;"`
}
