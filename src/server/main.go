package main

import (
	"file-storage/src/Tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"path/filepath"
)

var port string
var uploadDir string
var uploadPath string
var downloadPath string
var fileMovePath string
var fileDeletePath string

func init() {
	// 读取配置文件
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// 从配置文件中读取上传目录、端口号和路由路径
	server := viper.GetStringMap("server")
	port = server["port"].(string)
	routes := viper.GetStringMap("routes")
	uploadDir = routes["upload_dir"].(string)
	uploadPath = routes["upload_path"].(string)
	downloadPath = routes["download_path"].(string)
	fileMovePath = routes["file_move_path"].(string)
	fileDeletePath = routes["file_delete_path"].(string)

	// 创建一个 gin 实例
	r := gin.Default()

	// 设置路由
	r.POST(uploadPath, uploadHandler)
	r.GET(downloadPath, downloadHandler)
	r.POST(fileMovePath, moveFiles)
	r.POST(fileDeletePath, deleteFiles)

	// 启动 HTTP 服务器
	go func() {
		if err := r.Run(port); err != nil {
			panic(fmt.Errorf("Fatal error server: %s \n", err))
		}
	}()
}

// 处理文件上传请求
func uploadHandler(c *gin.Context) {
	file, err := c.FormFile("file_path")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	// 保存文件到本地
	filename := filepath.Join(uploadDir, file.Filename)
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

// 处理文件下载请求
func downloadHandler(c *gin.Context) {
	var req Tools.DownloadRequest
	if err := c.BindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("invalid request: %s", err.Error()))
		return
	}

	filename := req.Filename
	filepath := filepath.Join(uploadDir, filename)

	// 检查文件是否存在
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		c.String(http.StatusNotFound, fmt.Sprintf("file '%s' not found", filename))
		return
	}

	// 设置响应头
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/octet-stream")
	c.File(filepath)
}

// 处理文件移动请求
func moveFiles(c *gin.Context) {
	// 解析请求参数
	var req Tools.MoveFilesRequest
	if err := c.BindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("invalid request: %s", err.Error()))
		return
	}

	// 规范化路径
	sourcePath := filepath.Join(uploadDir, filepath.Clean(req.SourcePath))
	destinationPath := filepath.Join(uploadDir, filepath.Clean(filepath.Join(req.DestinationPath, filepath.Base(req.SourcePath))))

	// 检查源路径是否存在
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		c.String(http.StatusBadRequest, fmt.Sprintf("source path does not exist: %s", req.SourcePath))
		return
	}

	// 检查目标路径是否存在
	if _, err := os.Stat(filepath.Join(uploadDir, req.DestinationPath)); os.IsNotExist(err) {
		c.String(http.StatusBadRequest, fmt.Sprintf("destination path does not exist: %s", req.DestinationPath))
		return
	}

	// 移动文件
	if err := os.Rename(sourcePath, destinationPath); err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("failed to move file: %s", err.Error()))
		return
	}

	// 返回成功响应
	c.String(http.StatusOK, "file moved successfully")
}

// 处理删除文件/文件夹请求
func deleteFiles(c *gin.Context) {
	// 解析请求参数
	var req Tools.DeleteFilesRequest
	if err := c.BindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("invalid request: %s", err.Error()))
		return
	}

	// 规范化路径
	path := filepath.Join(uploadDir, filepath.Clean(req.Path))

	// 检查路径是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		c.String(http.StatusBadRequest, fmt.Sprintf("path does not exist: %s", req.Path))
		return
	}

	// 删除文件或文件夹
	if err := os.RemoveAll(path); err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("failed to delete file or directory: %s", err.Error()))
		return
	}

	// 返回成功响应
	c.String(http.StatusOK, "file or directory deleted successfully")
}

func main() {
	// 程序入口
	select {}
}
