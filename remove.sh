#!/bin/bash

# 1. 定义要清理的目录（在这里添加路径，支持绝对和相对路径）
TARGET_DIRS=(
    "./.inbox"
    "./.tasks"
    "./.team"
    "./sandbox"
)

# 2. 核心清理逻辑
for dir in "${TARGET_DIRS[@]}"; do
    # 检查目录是否存在且确实是一个目录
    if [ -d "$dir" ]; then
        echo "Cleaning: $dir"
        
        # 使用 find 命令递归删除目录下所有文件和子目录，但保留根目录本身
        # -mindepth 1 确保不删除 $dir 文件夹本身
        # -delete 直接由内核层执行删除，比 xargs rm 更快
        find "$dir" -mindepth 1 -delete
        
        if [ $? -eq 0 ]; then
            echo "Successfully cleaned $dir"
        else
            echo "Error: Failed to clean $dir (Check permissions)"
        fi
    else
        echo "Skip: $dir does not exist or is not a directory"
    fi
done