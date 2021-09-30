package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

/*
	RESTful API风格：通过http请求方式来描述操作
*/

type User struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

const (
	addr = "127.0.0.1:8080"
)

var users []User

func main() {
	// 模拟数据
	users = append(users, User{1, "zs"}, User{2, "ls"}, User{3, "ww"})
	g := gin.Default()

	g.GET("/:id", getHandler)
	g.POST("/add", addHandler)
	g.PUT("/update", updateHandler)
	g.DELETE("/:id", deleteHandler)

	if err := http.ListenAndServe(addr, g); err != nil {
		log.Fatalf("服务异常退出：%v", err)
	}
}

func remove(slice []User, i int) []User {
	// 最后一个元素
	if i == len(slice)-1 {
		return slice[:len(slice)-1]
	}

	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}

// 模拟数据库查询
func query(value interface{}) *User {
	u := new(User)
	for _, user := range users {
		if user.ID == value || user.Name == value {
			*u = user
			return u
		}
	}
	return nil
}

// 模拟数据库更新
func update(value User) {
	for i, user := range users {
		if user.ID == value.ID {
			users[i] = value
		}
	}
}

// 模拟数据库删除
func delete(id uint64) {
	for i, user := range users {
		if user.ID == id {
			users = remove(users, i)
		}
	}
}

func getHandler(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 0)
	user := query(uint64(id))
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": user,
	})
}

func addHandler(c *gin.Context) {
	var user User
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "fail",
		})
	}

	users = append(users, user)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": users,
	})

}

func updateHandler(c *gin.Context) {
	var user User
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "fail",
		})
	}

	update(user)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": users,
	})
}

func deleteHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	delete(uint64(id))
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": users,
	})
}
