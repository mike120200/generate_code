package create_file

import (
	"fmt"
	"os"
	"path/filepath"
)

// CreateFileSrv 生成需要的代码文件服务
type CreateFileSrv struct {
	fileMap map[string]string
	Dir     string
}

func NewCreateFileSrv(dir string) *CreateFileSrv {
	return &CreateFileSrv{
		fileMap: make(map[string]string),
		Dir:     dir,
	}
}

// AddFile 添加需要创建的文件
func (srv *CreateFileSrv) AddFile(fileName, code string) error {
	if fileName == "" {
		fmt.Printf("文件名不能为空")
		return fmt.Errorf("文件名不能为空")
	}
	if code == "" {
		fmt.Printf("代码不能为空")
		return fmt.Errorf("代码不能为空")
	}
	srv.fileMap[fileName] = code
	return nil
}

// CreateFile 生成需要的代码文件
func (srv *CreateFileSrv) Generate() error {

	//查看文件夹是否存在
	if _, err := os.Stat(srv.Dir); !os.IsNotExist(err) {
		fmt.Println("文件夹已存在")
	} else {
		fmt.Println("文件夹不存在,创建文件夹")
		err := os.Mkdir(srv.Dir, os.ModePerm)
		if err != nil {
			fmt.Printf("创建文件夹失败: %v\n", err)
			return err
		}
	}

	for k, v := range srv.fileMap {
		// 获取完整文件路径
		filePath := filepath.Join(srv.Dir, k)
		// 提取目录部分
		dirPath := filepath.Dir(filePath)

		// 检查并创建目录
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			err := os.MkdirAll(dirPath, os.ModePerm) // 创建多级目录
			if err != nil {
				fmt.Printf("创建目录失败: %v\n", err)
				return err
			}
			fmt.Printf("目录已创建: %s\n", dirPath)
		}
		// 创建并写入文件
		err := os.WriteFile(srv.Dir+"/"+k, []byte(v), 0644)
		if err != nil {
			fmt.Printf("写入文件失败: %v\n", err)
			return err
		}
	}

	return nil
}
