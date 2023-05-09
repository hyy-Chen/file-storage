package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	//传入目录地址
	dirPath := "./folder"

	//获取目录结构JSON对象
	dirJSON := getDirJSON(dirPath)

	//将JSON对象转换成字符串并打印
	jsonString, _ := json.MarshalIndent(dirJSON, "", "  ")
	fmt.Println(string(jsonString))
}

type FileNode struct {
	Name string `json:"name"`
	Type int    `json:"type"`
}

type Node struct {
	Name     string      `json:"name"`
	Type     int         `json:"type"`
	Children []*FileNode `json:"children"`
}

func getDirJSON(dirPath string) *Node {
	info, err := os.Stat(dirPath)
	if err != nil {
		panic(err)
	}

	//确认是否为目录
	if !info.IsDir() {
		panic("传入的不是目录")
	}

	//创建根节点Node
	dirJSON := Node{
		Name:     filepath.Base(dirPath),
		Type:     1,
		Children: []*FileNode{},
	}

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}

	//遍历目录下的所有文件（包括子目录）
	for _, file := range files {
		//如果是目录，递归调用本函数获取其Node，从顶层开始向下把Node加入当前节点的Children列表中
		if file.IsDir() {
			dirJSON.Children = append(dirJSON.Children, &FileNode{
				Name: file.Name(),
				Type: 1,
			})
		} else { //如果是普通文件，直接加入当前节点的Children列表中
			dirJSON.Children = append(dirJSON.Children, &FileNode{
				Name: file.Name(),
				Type: 2,
			})
		}
	}

	return &dirJSON
}
