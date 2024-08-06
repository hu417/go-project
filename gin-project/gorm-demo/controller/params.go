package controller

// req

type User struct {
	Name  string `json:"name" binding:"required"`
	Age   int    `json:"age" binding:"required"`
	Email string `json:"email" binding:"required"`
}

type UserRecord struct {
    UserID int64 `gorm:"int64"`
    Amount int64 `gorm:"amount"`
}