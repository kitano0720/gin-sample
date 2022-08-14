package main

import (
	"context"
	"database/sql"
	"gin-sample/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func main() {
	// コンテキスト取得
	ctx := context.Background()
	// DB接続
	db, err := sql.Open("mysql", "root:mysql@(localhost:3306)/gin_sample?parseTime=true")
	if err != nil {
		log.Fatalf("Cannot connect database %v", err)
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/users", func(c *gin.Context) {
		// ユーザー一覧取得
		users, err := models.Users().All(ctx, db)
		if err != nil {
			log.Fatalf("Get users error %v", err)
		}

		// テンプレートをロード
		r.LoadHTMLGlob("templates/users/index.html")

		c.HTML(http.StatusOK, "index.html", gin.H{
			"users": users,
		})
	})
	r.GET("/users/create", func(c *gin.Context) {
		// テンプレートをロード
		r.LoadHTMLGlob("templates/users/create.html")

		c.HTML(http.StatusOK, "create.html", gin.H{})
	})
	r.POST("/users", func(c *gin.Context) {
		// ユーザーを登録
		var user models.User
		c.Bind(&user)
		err := user.Insert(ctx, db, boil.Infer())
		if err != nil {
			log.Fatalf("Create user error: %v", err)
		}

		// ユーザー一覧ページにリダイレクト
		c.Redirect(http.StatusFound, "/users")
	})
	r.GET("/users/edit/:userId", func(c *gin.Context) {
		// ユーザーデータを取得
		userId, _ := strconv.Atoi(c.Param("userId"))
		user, _ := models.FindUser(ctx, db, userId)

		// テンプレートをロード
		r.LoadHTMLGlob("templates/users/edit.html")

		c.HTML(http.StatusOK, "edit.html", gin.H{
			"user": user,
		})
	})
	r.POST("/users/edit", func(c *gin.Context) {
		// ユーザーを更新
		var user models.User
		c.Bind(&user)
		_, err := user.Update(ctx, db, boil.Whitelist("user_name", "age", "updated_at"))
		if err != nil {
			log.Fatalf("Update user error: %v", err)
		}

		// ユーザー一覧ページにリダイレクト
		c.Redirect(http.StatusFound, "/users")
	})
	r.POST("/users/delete", func(c *gin.Context) {
		// ユーザーデータを取得
		userId, _ := strconv.Atoi(c.PostForm("user_id"))
		user, _ := models.FindUser(ctx, db, userId)

		// ユーザーを削除
		_, err := user.Delete(ctx, db)
		if err != nil {
			log.Fatalf("Delete user error: %v", err)
		}

		// ユーザー一覧ページにリダイレクト
		c.Redirect(http.StatusFound, "/users")
	})
	r.Run(":9000") // 0.0.0.0:9000 でサーバーを立てます。
}
