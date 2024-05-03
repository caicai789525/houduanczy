package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建Gin路由
	router := gin.Default()

	router.Static("/uploads", "./uploads")

	router.POST("/inference", func(c *gin.Context) {
		// 从请求中解析文件和参数
		file, err := c.FormFile("driven_audio")
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}
		drivenAudioPath := "uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, drivenAudioPath); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
		sourceImage := c.PostForm("source_image")
		resultDir := c.PostForm("result_dir")

		// 构造命令
		cmd := exec.Command("python", "inference.py",
			"--driven_audio", drivenAudioPath,
			"--source_image", sourceImage,
			"--still",
			"--preprocess", "full",
			"--enhancer", "gfpgan",
		)

		// 执行命令
		out, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}

		// 返回结果
		c.String(http.StatusOK, fmt.Sprintf("Output: %s\n", out))
	})

	// 启动服务器
	port := ":8080"
	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(router.Run(port))
}
