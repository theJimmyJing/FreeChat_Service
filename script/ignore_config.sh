#!/usr/bin/env bash
source ./style_info.cfg

file_urls=(
  "../config/config.yaml"
  "../config/freechat.yaml"
)

echo -e "Do you want to ignore the local config?"
echo -e ${YELLOW_PREFIX}"New file please commit an empty file to github!"
read -t 60 -p "Please input y/n/q:" yes_no

if [[ ${yes_no} == "y" ]]; then
  for i in ${file_urls[*]}; do
    # echo "git update-index --assume-unchanged ${i}" # 忽略本地文件更改，也忽略远程更新
    git update-index --skip-worktree ${i} # 忽略本地文件更改，但远程更新需要覆盖本地
  done
  echo -e ${GREEN_PREFIX}"Success!!!"
fi

if [[ ${yes_no} = "n" ]]; then
  for i in ${file_urls[*]}; do
    # echo "git update-index --no-assume-unchanged ${i}"
      git update-index --no-skip-worktree ${i}
  done
  echo -e ${GREEN_PREFIX}"Success!!!"
fi

if [[ ${yes_no} = "n" ]]; then
  exit 0
fi
