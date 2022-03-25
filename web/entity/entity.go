package entity

type UserEntity struct {
	Id       string `bson:"id" json:"id"`
	Username string `bson:"username" json:"username" form:"username" binding:"required"`
	Password string `bson:"password" json:"password" form:"password" binding:"required"`
	Role     string `bson:"role" json:"role"`
}
