package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type C struct {
	Id int
	Name string
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ping",
		})
	})

	r.GET("/hello/:num", func(c *gin.Context) {
		key := c.Param("num")
		c.JSON(200, gin.H{
			"hello": key,
		})
	})

	// このリクエストは /welcome?firstname=Jane&lastname=Doe へ返答する
	r.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest") // Geustはデフォルト値?
		lastname := c.Query("lastname") // c.Request.URL.Query().Get("lastname") のショートカット

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})

	r.LoadHTMLGlob("templates/*.tmpl")
	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"a": "a",
			"b": []string{"b_todo1","b_todo2"},
			"c": []C{{1,"c_mika"},{2,"c_risa"}},
			"d": C{3,"d_mayu"},
			"e": true,
			"f": false,
			"h": true,
		})
	})
	// ポートを設定しています。
	r.Run(":3001")
}
