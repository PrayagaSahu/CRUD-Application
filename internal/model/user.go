package model

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Email    string `json:"email" gorm:"uniqueIndex"`
	Phone    string `json:"phone" validate:"required,e164" gorm:"uniqueIndex"`
	Age      int    `json:"age" validate:"gte=0,lte=130"`
	Role     string `json:"role" validate:"required,oneof=admin user viewer"`
	Password string `json:"password,omitempty"` // never return in response

}
