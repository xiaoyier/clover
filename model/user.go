package model

import "time"

func init() {
	RegisterAutoMigrates(&User{})
}

type UserRegisterReq struct {
	UserName   string `json:"user_name" binding:"required,min=6,max=20"`
	Password   string `json:"password" binding:"required,min=6,max=30"`
	RePassword string `json:"re_password" binding:"required,min=6,max=30,eqfield=Password"`
}

type UserLoginReq struct {
	UserName string `json:"user_name" binding:"required,min=6,max=20"`
	Password string `json:"password" binding:"required,min=6,max=30"`
}

type UserLoginRsp struct {
	UserID       int64  `json:"user_id,string"`
	UserName     string `json:"user_name"`
	LoginToken   string `json:"login_token"`
	RefreshToken string `json:"refresh_token"`
}

type User struct {
	ID         uint      `gorm:"primary_key;column:id"`
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime(6)"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:datetime(6)"`
	UserID     int64     `gorm:"unique;column:user_id;type:bigint(20)"`
	UserName   string    `gorm:"unique_index;column:user_name;type:varchar(50)"`
	UserPasswd string    `gorm:"column:user_passwd;type:varchar(50)"`
	Email      string    `gorm:"column:email;type:varchar(30)"`
}

func (u *User) TableName() string {
	return "user"
}

func QueryUserByUserName(userName string) (user *User, err error) {
	user = new(User)
	d := db.Where("user_name = ?", userName).First(&user)
	if d.RecordNotFound() {
		return nil, nil
	}

	return user, d.Error
}

func QueryUserByUserID(userID int64) (user *User, err error) {
	user = new(User)
	d := db.Where("user_id = ?", userID).First(&user)
	if d.RecordNotFound() {
		return nil, nil
	}

	return user, d.Error
}

func (u *User) Insert() error {

	d := db.Save(u)
	return d.Error
}

func (u *User) Update() {
	db.Model(u).Where(User{UserID: u.UserID}).Update(u)
}
