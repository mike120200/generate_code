#!/bin/bash

# 遇到错误立即终止
set -e
# 设置目录和文件名
TEST_DIR="./test"
CODE_DIR="./code"
OUTPUT_FILE="output_binary"

# 1. 清空 test 文件夹
echo "清空 test 文件夹内容..."
rm -rf "$TEST_DIR"
mkdir -p "$TEST_DIR" # 重新创建 test 文件夹

if [ $? -ne 0 ]; then
    echo "清空 test 文件夹失败！"
    exit 1
fi
echo "test 文件夹已清空并重建"

# 2. 编译 code 文件夹中的代码
echo "开始编译代码..."
cd "$CODE_DIR" || exit 1
go build -o "$OUTPUT_FILE"

if [ $? -ne 0 ]; then
    echo "代码编译失败！"
    exit 1
fi
echo "代码编译成功"

# 3. 将生成的二进制文件移动到 test 文件夹
echo "移动编译后的文件到 test 文件夹..."
mv "$OUTPUT_FILE" "../$TEST_DIR/"

if [ $? -ne 0 ]; then
    echo "移动文件失败！"
    exit 1
fi
echo "文件已成功移动到 $TEST_DIR"

# 4. 运行编译后的文件
echo "运行编译后的文件..."
cd "../$TEST_DIR" || exit 1
./"$OUTPUT_FILE"

if [ $? -ne 0 ]; then
    echo "运行文件失败！"
    exit 1
fi
echo "文件运行结束"