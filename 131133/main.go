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

	// 设置静态文件目录
	router.Static("/uploads", "./uploads")

	router.POST("/inference", func(c *gin.Context) {
		// 从请求中解析文件和参数
		drivenAudioFile, err := c.FormFile("driven_audio")
		if err != nil {
			log.Printf("Error getting driven_audio file: %s", err)
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}
		drivenAudioPath := "uploads/" + drivenAudioFile.Filename
		if err := c.SaveUploadedFile(drivenAudioFile, drivenAudioPath); err != nil {
			log.Printf("Error saving driven_audio file: %s", err)
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		sourceImage := c.PostForm("source_image")
		refEyeblink := c.PostForm("ref_eyeblink")
		refPose := c.PostForm("ref_pose")

		// 上传 ppt 文件
		pptFile, err := c.FormFile("ppt")
		if err != nil {
			log.Printf("Error getting ppt file: %s", err)
			c.String(http.StatusBadRequest, fmt.Sprintf("get ppt file err: %s", err.Error()))
			return
		}
		pptPath := "uploads/" + pptFile.Filename
		if err := c.SaveUploadedFile(pptFile, pptPath); err != nil {
			log.Printf("Error saving ppt file: %s", err)
			c.String(http.StatusBadRequest, fmt.Sprintf("upload ppt file err: %s", err.Error()))
			return
		}

		// 上传背景图片
		bgImageFile, err := c.FormFile("bg_image")
		if err != nil {
			log.Printf("Error getting bg_image file: %s", err)
			c.String(http.StatusBadRequest, fmt.Sprintf("get bg_image file err: %s", err.Error()))
			return
		}
		bgImagePath := "uploads/" + bgImageFile.Filename
		if err := c.SaveUploadedFile(bgImageFile, bgImagePath); err != nil {
			log.Printf("Error saving bg_image file: %s", err)
			c.String(http.StatusBadRequest, fmt.Sprintf("upload bg_image file err: %s", err.Error()))
			return
		}

		// 构造命令
		cmd := exec.Command("python", "inference.py",
			"--driven_audio", drivenAudioPath,
			"--source_image", sourceImage,
			"--ppt", pptPath,
			"--bg_image", bgImagePath,
			"--preprocess", "full",
			"--enhancer", "gfpgan",
		)
		if refEyeblink != "" {
			cmd.Args = append(cmd.Args, "--ref_eyeblink", refEyeblink)
		}
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
