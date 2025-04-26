package tools

import (
	"encoding/gob"
	"yukoreimburse/users"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func init() {
	store = sessions.NewCookieStore([]byte("recomend-32bytes-at-least"))
	store.Options = &sessions.Options{
	    Path: "/",
		MaxAge: 0,
		Secure: false,
		HttpOnly: true,

		Domain: "",
		SameSite: http.SameSiteLaxMode,
	}

	gob.Register(users.Users{})
}

func SetSession(ctx *gin.Context, key string, value interface{}) error{

	session, err := store.Get(ctx.Request, "yukoreimburseID")
	if err != nil {
		return err
	}

	session.Values[key] = value
	fmt.Println(session.Values[key])
	return session.Save(ctx.Request, ctx.Writer)
}

func GetSession(ctx *gin.Context, key string) (interface{}) {

    session, err := store.Get(ctx.Request, "yukoreimburseID")
    if err != nil {
        return err
    }

	fmt.Println(session.Values[key])

    return session.Values[key]
}

func DelSession(ctx *gin.Context, key string) error {

	session, err := store.Get(ctx.Request, "yukoreimburseID")
	if err != nil {
		return err
	}

	delete(session.Values, key)
	return session.Save(ctx.Request, ctx.Writer)
}