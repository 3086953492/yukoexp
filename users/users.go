package users

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"yukoreimburse/database"
	"fmt"
	"strings"
)

type Users struct {
	User_id  int    `gorm:"primarykey"`
	Username string `gorm:"size:30;unique"`
	Password string `gorm:"size:128"`
	Is_admin int    `gorm:"default:0"`
	Lark_id  string `gorm:"size:50;unique"`
}

type UsernameList struct {
    User_id int `gorm:"primarykey"`
    Username string `gorm:"size:30;unique"`
}

func CheckUser(username string, password string) (*Users, error) {
	// 去空格
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	// 合法性验证，不能为空
	if username == "" || password == "" {
		return nil, errors.New("用户名或密码不能为空")
	}

	var user Users

	database.Gdb.Where("username = ?", username).Find(&user)

	if user.Username == "" {
		return nil, errors.New("用户名不存在")
	}

	// 密码要转md5
	hash := md5.Sum([]byte(password))

	if user.Password != hex.EncodeToString(hash[:]) {
		return nil, errors.New("密码错误")
	}
	fmt.Println(user)

	fmt.Println("登录成功")

	return &user, nil
}

func GetUsernameList() ([]UsernameList, error) {
    var usernameList []UsernameList
    err := database.Gdb.Raw("SELECT user_id, username FROM users").Scan(&usernameList).Error
    if err != nil {
        return nil, err
    }
    return usernameList, nil
}
