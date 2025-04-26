package tools

import (
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/rand"
)

func GetCaptcha(ctx *gin.Context) string {
	var captcha string

	for i := 0; i < 5; i++ {
		key := rand.Intn(3)
		var ascii int
		if key == 0 {
		    // 数字
			ascii = 49 + rand.Intn(9)
		}else if key == 1 {
		    // 大写字母
			ascii = 65 + rand.Intn(26)
		}else {
		    // 小写字母
			ascii = 97 + rand.Intn(26)
		}
		captcha += string(ascii)
	}
	SetSession(ctx, "captcha", captcha)
	return captcha
}

func CheckCaptcha(ctx *gin.Context, captcha string) bool {
	// 不区分大小写
    sessionCaptcha := GetSession(ctx, "captcha").(string)
	return strings.EqualFold(sessionCaptcha, captcha)
}