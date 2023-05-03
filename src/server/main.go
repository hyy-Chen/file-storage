package main

import (
	"encoding/json"
	"file-storage/src/Tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io"
	"mime/multipart"
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
var createFolderPath string
var uploadFolderPath string

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
	createFolderPath = routes["create_folder_path"].(string)
	uploadFolderPath = routes["upload_folder_path"].(string)

	// 创建一个 gin 实例
	r := gin.Default()

	// 设置路由
	r.POST(uploadPath, uploadHandler)
	r.GET(downloadPath, downloadHandler)
	r.POST(fileMovePath, moveFiles)
	r.POST(fileDeletePath, deleteFiles)
	r.POST(createFolderPath, createFolder)
	r.POST(uploadFolderPath, uploadFolder)

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
	filePath := filepath.Join(uploadDir, filename)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.String(http.StatusNotFound, fmt.Sprintf("file '%s' not found", filename))
		return
	}

	// 设置响应头
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/octet-stream")
	c.File(filePath)
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

// 处理创建文件夹请求
func createFolder(c *gin.Context) {
	// 解析请求参数
	var req Tools.CreateFolderRequest
	if err := c.BindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("invalid request: %s", err.Error()))
		return
	}

	// 规范化路径
	path := filepath.Join(uploadDir, filepath.Clean(req.Path))

	// 检查路径是否已存在
	if _, err := os.Stat(path); err == nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("path already exists: %s", req.Path))
		return
	}

	// 创建文件夹
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("failed to create folder: %s", err.Error()))
		return
	}

	// 返回成功响应
	c.String(http.StatusOK, "folder created successfully")
}

// 层级创建文件夹的方法
func uploadFolder(c *gin.Context) {
	// 解析数据，先规定大小
	err := c.Request.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 获取文件信息
	files := c.Request.MultipartForm.File[Tools.UploadFolderFiles]
	// 建立映射关系，文件名映射到文件
	filesMap := make(map[string]*multipart.FileHeader)
	for _, file := range files {
		filesMap[file.Filename] = file
	}
	// 获取文件夹的层级关系
	// 先获取目标目录，查看是否存在
	dirPath := filepath.Join(uploadDir, c.PostForm(Tools.UploadFolderDestinationDirectory))
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		c.String(http.StatusBadRequest, fmt.Sprintf("source path does not exist: %s", c.PostForm(Tools.UploadFolderDestinationDirectory)))
		return
	}
	// 获取json
	jsonStr := c.PostForm(Tools.UploadFolderJson)
	var directory Tools.Directory
	// 解析json
	if err := json.Unmarshal([]byte(jsonStr), &directory); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("failed to parse JSON data: %s", err.Error()))
		return
	}
	var errs []string
	// 进行递归层级建文件夹以及拷贝文件
	var save func(string2 string, directory2 *Tools.Directory)
	// 匿名函数进行层次建立文件夹，存储文件的操作
	save = func(dir string, directory *Tools.Directory) {
		// 如果是目录
		basePath := filepath.Join(dir, directory.Name)
		// 先暂定不能覆盖原文件夹， 即出现相同文件夹名称就发出失败警告，并停止递归当前目录
		path := filepath.Join(dirPath, basePath)
		if directory.Type == Tools.Folder {
			if _, err := os.Stat(path); os.IsExist(err) {
				errs = append(errs, fmt.Sprintf("directory %s is exist", basePath))
				return
			}
			// 如果创建文件夹失败，就递归当前目录并且存储报错信息
			if err := os.Mkdir(path, 0755); err != nil {
				errs = append(errs, fmt.Sprintf("folder:%s error:%s", basePath, err))
				return
			}
			// 继续递归目录
			for _, child := range directory.Children {
				save(basePath, &child)
			}
			// 如果是文件
		} else {
			// 若找不到对应的文件信息，就停止操作并且存储报错信息
			if fileHeader, ok := filesMap[directory.CName]; !ok {
				errs = append(errs, fmt.Sprintf("Unable to find file information:%s", basePath))
				return
			} else {
				// 进行写入
				// 创建对应的文件
				fileHandle, err := os.Create(path)
				if err != nil {
					errs = append(errs, fmt.Sprintf("Create file(%s) fail, err:%s", basePath, err))
					return
				}
				// 打开对应文件io流
				file, err := fileHeader.Open()
				if err != nil {
					errs = append(errs, fmt.Sprintf("Failed to read file(%s) information, err:%s", basePath, err))
					return
				}
				// 进行io流写入
				if _, err := io.Copy(fileHandle, file); err != nil {
					errs = append(errs, fmt.Sprintf("Failed to copy file(%s) information, err:%s", basePath, err))
					return
				}
			}
		}
	}
	save("", &directory)
	if len(errs) != 0 {
		c.String(http.StatusAlreadyReported, fmt.Sprintf("something file upload fail\n%v", errs))
	} else {
		// 返回成功响应
		c.String(http.StatusOK, "folder upload successfully")
	}
}

func main() {
	// 程序入口
	select {}
}
