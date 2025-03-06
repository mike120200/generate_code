package main

import (
	"fmt"
	"generate/create_file"
	"generate/generated_code"
	"os/exec"
)

func main() {
	//输入项目的名称
	fmt.Println("Please enter a project name:")
	var prjName string = "default_prj"
	_, err := fmt.Scanln(&prjName)
	if err != nil {
		fmt.Println("enter Error:", err)
		return
	}
	fmt.Println("add file to project")
	//添加文件
	instance := create_file.NewCreateFileSrv("./common")

	//日志模块
	if err := instance.AddFile("log/log.go", generated_code.LogCode.ReplaceProjectName(prjName)); err != nil {
		fmt.Println("AddFile Error:", err)
		return
	}
	if err := instance.AddFile("log/log_test.go", generated_code.LogCodeTest.ReplaceProjectName(prjName)); err != nil {
		fmt.Println("AddFile Error:", err)
		return
	}

	//pg数据库模块
	//选择pg的驱动
	fmt.Println("Please enter a pg drive(1:gorm,2:pgx/v4,default:pgx/v4):")
	pgDrive := ""
	_, err = fmt.Scanln(&pgDrive)
	if err != nil {
		fmt.Println("enter Error:", err)
		return
	}
	if pgDrive == "1" {
		if err := instance.AddFile("pg_conn/pg_conn.go", string(generated_code.PgconnCode_gorm.ReplaceProjectName(prjName))); err != nil {
			fmt.Println("AddFile Error:", err)
			return
		}
		if err := instance.AddFile("pg_conn/pg_conn_test.go", string(generated_code.PgconnCodeTest_gorm.ReplaceProjectName(prjName))); err != nil {
			fmt.Println("AddFile Error:", err)
			return
		}
	} else {

		if err := instance.AddFile("pg_conn/pgconn.go", generated_code.PgconnCode.ReplaceProjectName(prjName)); err != nil {
			fmt.Println("AddFile Error:", err)
			return
		}

		if err := instance.AddFile("pg_conn/pgconn_test.go", generated_code.PgconnCodeTest.ReplaceProjectName(prjName)); err != nil {
			fmt.Println("AddFile Error:", err)
			return
		}
	}

	//redis模块
	if err := instance.AddFile("redis_conn/redis_conn.go", generated_code.RedisCode.ReplaceProjectName(prjName)); err != nil {
		fmt.Println("AddFile Error:", err)
		return
	}
	if err := instance.AddFile("redis_conn/redis_conn_test.go", generated_code.RedisCodeTest.ReplaceProjectName(prjName)); err != nil {
		fmt.Println("AddFile Error:", err)
		return
	}

	//配置文件模块
	if err := instance.AddFile("config/config.go", generated_code.ConfigCode.ReplaceProjectName(prjName)); err != nil {
		fmt.Println("AddFile Error:", err)
		return
	}
	if err := instance.AddFile("config/config_test.go", generated_code.ConfigCodeTest.ReplaceProjectName(prjName)); err != nil {
		fmt.Println("AddFile Error:", err)
		return
	}

	//结果模块
	if err := instance.AddFile("result/response.go", generated_code.ResultGeneratedCode.ReplaceProjectName(prjName)); err != nil {
		fmt.Println("AddFile Error:", err)
		return
	}
	if err := instance.AddFile("result/response_test.go", generated_code.ResultCodeTest.ReplaceProjectName(prjName)); err != nil {
		fmt.Println("AddFile Error:", err)
		return
	}

	fmt.Println("generate files......")
	//生成文件
	if err := instance.Generate(); err != nil {
		fmt.Println("Generate Error:", err)
		return
	}
	fmt.Println("Generate Success!")
	//golang项目初始化
	fmt.Println("project init......")
	//定义命令
	initCmd := exec.Command("go", "mod", "init", prjName)
	//执行命令
	output, err := initCmd.Output()
	if err != nil {
		fmt.Println("golang project init Error:", err)
		return
	}
	// 打印命令的输出
	fmt.Printf("golang init output:\n%s\n", string(output))

	//golang依赖初始化
	depCmd := exec.Command("go", "mod", "tidy")
	output, err = depCmd.Output()
	if err != nil {
		fmt.Println("dep Error:", err)
		return
	}
	fmt.Printf("golang dep output:\n%s\n", string(output))

	//cobra初始化
	cobraCmd := exec.Command("cobra-cli", "init")
	output, err = cobraCmd.Output()
	if err != nil {
		fmt.Println("cobra init Error:", err)
		return
	}
	fmt.Printf("cobra init output:\n%s\n", string(output))

	fmt.Println("project init success!")

}
