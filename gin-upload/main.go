package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()
	/*
		g.MaxMultipartMemory = 8 << 20 // 8 MiB
		上传的文件们按顺序存入内存中，累加大小不得超出 8Mb ，最后累加超出的文件就存入系统的临时文件中。非文件字段部分不计入累加。所以这种情况，文件上传是没有任何限制的。
		要想限制文件大小，其实是要限制整个请求body的大小
	*/

	r := g.Group("/api")
	{
		r.POST("/upload", Upload)
		r.POST("/mult", UploadMult)
	}

	srv := &http.Server{
		Addr:    "0.0.0.0:9999",
		Handler: g,
	}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Could not listen on ", srv.Addr)
	}
}

// 单文件上传
func Upload(c *gin.Context) {

	file, _ := c.FormFile("example")
	log.Println("file name: ", file.Filename)

	// 上传文件到指定路径，一般都是路径+file.Filename
	dst := file.Filename
	err := c.SaveUploadedFile(file, dst)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("upload file err: %s", err.Error()))
	} else {
		c.String(http.StatusOK, "upload file success")
	}

}

// 多文件上传
func UploadMult(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["uploads"]
	for _, file := range files {
		log.Println(file.Filename)

		// 上传文件到指定路径
		// c.SaveUploadedFile(file, dst)
	}
	c.String(http.StatusOK, "upload mult file success")
}
