package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"strconv"
)

type Todo struct {
	gorm.Model
	Text   string
	Status string
}

// DB初期化
func dbInit() {
	db, err := gorm.Open("sqlite3", "todo.sqlite3")
	if err != nil {
		panic("データベースが開ません(dbInit)")
	}
	db.AutoMigrate(&Todo{})
	defer db.Close()
}

// DB追加
func dbInsert(text string, status string) {
	db, err := gorm.Open("sqlite3", "todo.sqlite3")
	if err != nil {
		panic("データベースが開ません(dbInsert)")
	}
	db.Create(&Todo{Text: text, Status: status})
	defer db.Close()
}

// DB更新
func dbUpdate(id int, text string, status string) {
	db, err := gorm.Open("sqlite3", "todo.sqlite3")
	if err != nil {
		panic("データベースが開ません(dbUpdate)")
	}
	var todo Todo
	db.First(&todo, id)
	todo.Text = text
	todo.Status = status
	db.Save(&todo)
	db.Close()
}

// DB削除
func dbDelete(id int) {
	db, err := gorm.Open("sqlite3", "todo.sqlite3")
	if err != nil {
		panic("データベースが開ません(dbDelete)")
	}
	var todo Todo
	db.First(&todo, id)
	db.Delete(&todo)
	db.Close()
}

// DB全取得
func dbGetAll() []Todo {
	db, err := gorm.Open("sqlite3", "todo.sqlite3")
	if err != nil {
		panic("データベースが開ません(dbGetAll)")
	}
	var todos []Todo
	db.Order("created_at desc").Find(&todos)
	db.Close()
	return todos
}

// DB一つ取得
func dbGetOne(id int) Todo {
	db, err := gorm.Open("sqlite3", "todo.sqlite3")
	if err != nil {
		panic("データベースが開ません(dbGetOne)")
	}
	var todo Todo
	db.First(&todo, id)
	db.Close()
	return todo
}

func main() {
	router := gin.Default()
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*tmpl")

	dbInit()

	// Index
	router.GET("/", func(c *gin.Context) {
		todos := dbGetAll()
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"todos": todos,
		})
	})

	// Create
	router.POST("/new", func(c *gin.Context) {
		text := c.PostForm("text")
		status := c.PostForm("status")
		dbInsert(text, status)
		c.Redirect(http.StatusFound, "/")
	})

	// Detail
	router.GET("detail/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		todo := dbGetOne(id)
		c.HTML(http.StatusOK, "detail.tmpl", gin.H{"todo": todo})
	})

	// Update
	router.POST("/update/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		text := c.PostForm("text")
		status := c.PostForm("status")
		dbUpdate(id, text, status)
		c.Redirect(http.StatusFound, "/")
	})

	// 削除確認
	router.GET("delete_check/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		todo := dbGetOne(id)
		c.HTML(http.StatusOK, "delete.tmpl", gin.H{"todo": todo})
	})

	// Delete
	router.POST("/delete/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		dbDelete(id)
		c.Redirect(http.StatusFound, "/")
	})
	router.Run("localhost:8080")
}
