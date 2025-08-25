#!/bin/bash

TAG=$1

if [ -z "$TAG" ]; then
  echo "用法: $0 <tag>"
  exit 1
fi

# 检查 tag 是否已存在
if git rev-parse "$TAG" >/dev/null 2>&1; then
  echo "Tag $TAG 已存在，跳过创建。"
else
  git tag "$TAG"
  git push origin "$TAG"
  echo "已创建并推送 tag: $TAG"
fi