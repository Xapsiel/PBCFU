package dewu

type User struct {
	ID             int    `json:"id" db:"id"`
	Login          string `json:"login" binding:"required"`
	Email          string `json:"email"`
	Password       string `json:"password" binding:"required"`
	RepeatPassword string `json:"repeatpassword"`
	LastClick      int    `json:"LastClick"`
}
