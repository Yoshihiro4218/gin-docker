# gin-docker
Go with gin on Docker


# APIサンプル ( `http://sekitaka-1214.hatenablog.com/entry/2016/08/11/153816` 様 )

```
アを使う。
    authorized.Use(AuthRequired())
    {
        authorized.POST("/login", loginEndpoint)
        authorized.POST("/submit", submitEndpoint)
        authorized.POST("/read", readEndpoint)

        // ネストされたグループ
        testing := authorized.Group("testing")
        testing.GET("/analytics", analyticsEndpoint)
    }

    // Listen and server on 0.0.0.0:8080
    r.Run(":8080")
}
モデルバインディングとバリデーション
リクエストボディを型にバインドするには、モデルバインディングを使用してください。 GinはJSON,XMLと標準的なフォームの値(foo=bar&hoge=fuga)をサポートしています。

// JSONからバインド
type Login struct {
    User     string `form:"user" json:"user" binding:"required"`
    Password string `form:"password" json:"password" binding:"required"`
}

func main() {
    router := gin.Default()


    // JSONをバインディングする例({"user": "manu", "password": "123"})
    router.POST("/loginJSON", func(c *gin.Context) {
        var json Login // Login型変数
        if c.BindJSON(&json) == nil { // c.BindJSON関数を使ってバインドする BindJSONの戻り値はerrorオブジェクトかな
            if json.User == "manu" && json.Password == "123" {
                c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
            } else {
                c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
            }
        }
    })

    // HTML formをバインドする例(user=manu&password=123)
    router.POST("/loginForm", func(c *gin.Context) {
        var form Login // Login型変数
    // content-typeにより元データの形式(JSON,XML,form)を予想してバインドする
        if c.Bind(&form) == nil {
            if form.User == "manu" && form.Password == "123" {
                c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
            } else {
                c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
            }
        }
    })

    // Listen and server on 0.0.0.0:8080
    router.Run(":8080")
}
Multipart/Urlencoded(ポスト) のバインディング
同様にできる

XML,JSON,YAMLの結果を返す
func main() {
    r := gin.Default()

    // gin.H は map[string]interface{} へのショートカット
    r.GET("/someJSON", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
    })

    r.GET("/moreJSON", func(c *gin.Context) {
        // 構造体も使える
        var msg struct {
            Name    string `json:"user"`
            Message string
            Number  int
        }
        msg.Name = "Lena"
        msg.Message = "hey"
        msg.Number = 123
    // この msg.Name は JSONないでは "user"となる
    // 出力{"user": "Lena", "Message": "hey", "Number": 123}
        c.JSON(http.StatusOK, msg)
    })

    r.GET("/someXML", func(c *gin.Context) {
        c.XML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
    })

    r.GET("/someYAML", func(c *gin.Context) {
        c.YAML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
    })

    // Listen and server on 0.0.0.0:8080
    r.Run(":8080")
}
静的ファイル
静的コンテンツは以下のようにディレクトリやファイルごとにマッピングしておける。

func main() {
    router := gin.Default()
    router.Static("/assets", "./assets")
    router.StaticFS("/more_static", http.Dir("my_file_system"))
    router.StaticFile("/favicon.ico", "./resources/favicon.ico")

    // Listen and server on 0.0.0.0:8080
    router.Run(":8080")
}
HTMLを返却する
LoadHTMLTemplates()関数を使う。

func main() {
    router := gin.Default()
    router.LoadHTMLGlob("templates/*") // 事前にテンプレートをロード
    //router.LoadHTMLFiles("templates/template1.html", "templates/template2.html") // ファイル指定でロード
    router.GET("/index", func(c *gin.Context) {
        // テンプレートを使って、値を置き換えてHTMLレスポンスを応答
        c.HTML(http.StatusOK, "index.tmpl", gin.H{
            "title": "Main website",
        })
    })
    router.Run(":8080")
}
templates/index.html
<html>
    <h1>
        {{ .title }}
    </h1>
</html>
同じ名前でディレクトリが異なるテンプレートの使い方
func main() {
    router := gin.Default()
    router.LoadHTMLGlob("templates/**/*")
    router.GET("/posts/index", func(c *gin.Context) {
        c.HTML(http.StatusOK, "posts/index.tmpl", gin.H{
            "title": "Posts",
        })
    })
    router.GET("/users/index", func(c *gin.Context) {
        c.HTML(http.StatusOK, "users/index.tmpl", gin.H{
            "title": "Users",
        })
    })
    router.Run(":8080")
}
リダイレクト
r.GET("/test", func(c *gin.Context) {
    c.Redirect(http.StatusMovedPermanently, "http://www.google.com/")
})
独自ミドルウェアの使い方
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        t := time.Now()

        // Set example variable
        c.Set("example", "12345")

        // リクエスト前

        c.Next()

        // リクエスト後
        latency := time.Since(t)
        log.Print(latency)

        // 送信したリクエストにアクセスできる
        status := c.Writer.Status()
        log.Println(status)
    }
}

func main() {
    r := gin.New()
    r.Use(Logger())

    r.GET("/test", func(c *gin.Context) {
        example := c.MustGet("example").(string)

        // it would print: "12345"
        log.Println(example)
    })

    // Listen and server on 0.0.0.0:8080
    r.Run(":8080")
}
ミドルウェア内でのゴルーチン
ミドルウェアやハンドラー内でゴルーチンを使う場合、オリジナルのcontextは使うべきでない。リードオンリーコピーを使うようにしなさい。

func main() {
    r := gin.Default()

    r.GET("/long_async", func(c *gin.Context) {
        // create copy to be used inside the goroutine
        cCp := c.Copy()
        go func() {
            // simulate a long task with time.Sleep(). 5 seconds
            time.Sleep(5 * time.Second)

            // note that you are using the copied context "cCp", IMPORTANT
            log.Println("Done! in path " + cCp.Request.URL.Path)
        }()
    })

    r.GET("/long_sync", func(c *gin.Context) {
        // simulate a long task with time.Sleep(). 5 seconds
        time.Sleep(5 * time.Second)

        // since we are NOT using a goroutine, we do not have to copy the context
        log.Println("Done! in path " + c.Request.URL.Path)
    })

    // Listen and server on 0.0.0.0:8080
    r.Run(":8080")
}
独自のHTTPの設定
http.ListenAndServe()を以下のように直接使う

func main() {
    router := gin.Default()
    http.ListenAndServe(":8080", router)
}
または

func main() {
    router := gin.Default()

    s := &http.Server{
        Addr:           ":8080",
        Handler:        router,
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }
    s.ListenAndServe()
}
グレースフルリスタート
グレースフルリスタートしたいですか？いくつかの方法があります。 fvbock/endlessでデフォルトのListenAndServeを置換できます。

router := gin.Default()
router.GET("/", handler)
// [...]
endless.ListenAndServe(":4242", router)
```
