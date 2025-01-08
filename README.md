# generate_code

## 功能说明

+ 生成每个项目中都需要拥有的、重复性较高、能通用的代码
+ 依赖初始化
+ cobra初始化

## 文件目录结构


```
generate_code
├─ README.md
├─ code
│  ├─ create_file
│  │  └─ create_file.go
│  ├─ generated_code
│  │  └─ generated_code.go
│  ├─ go.mod
│  ├─ go.sum
│  ├─ main.go
│  └─ readme.md
├─ srv_test.sh
└─ test
   ├─ LICENSE
   ├─ cmd
   │  └─ root.go
   ├─ common
   │  ├─ be_config
   │  │  ├─ config.go
   │  │  └─ config_test.go
   │  ├─ log
   │  │  ├─ log.go
   │  │  └─ log_test.go
   │  ├─ pg_conn
   │  │  ├─ pgconn.go
   │  │  └─ pgconn_test.go
   │  └─ redis_conn
   │     ├─ redis_conn.go
   │     └─ redis_conn_test.go
   ├─ go.mod
   ├─ go.sum
   ├─ main.go
   └─ output_binary

```

> 说明
>
> + code：存放“生成代码”功能的代码
> + test：用于测试“生成代码”功能，保存功能输出的结果
> + srv_test.sh：测试脚本，会运行代码，并将`go build`得到的编译后的文件移动到文件夹`test`，并运行

