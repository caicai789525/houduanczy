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
			log.Printf("Error getting driven_audio file: %s", err)
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}
		drivenAudioPath := "uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, drivenAudioPath); err != nil {
			log.Printf("Error saving driven_audio file: %s", err)
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
		sourceImage := c.PostForm("source_image")

		// 获取可选参数 ref_eyeblink 和 ref_pose
		refEyeblink := c.PostForm("ref_eyeblink")
		refPose := c.PostForm("ref_pose")

		// 构造命令
		cmd := exec.Command("python", "inference.py",
			"--driven_audio", drivenAudioPath,
			"--source_image", sourceImage,
			"--preprocess", "full",
			"--enhancer", "gfpgan",
		)

		// 如果存在 ref_eyeblink 参数，则添加到命令中
		if refEyeblink != "" {
			cmd.Args = append(cmd.Args, "--ref_eyeblink", refEyeblink)
		}

		// 如果存在 ref_pose 参数，则添加到命令中
		if refPose != "" {
			cmd.Args = append(cmd.Args, "--ref_pose", refPose)
		}

		// 执行命令
		out, err := cmd.Output()
		if err != nil {
			log.Printf("Error executing command: %s", err)
			c.String(http.StatusInternalServerError, fmt.Sprintf("error executing command: %s", err))
			return
		}

		// 返回结果
		c.String(http.StatusOK, fmt.Sprintf("Output: %s\n", out))
	})

	// 启动服务器
	port := ":8080"
	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(router.Run(port))
}
