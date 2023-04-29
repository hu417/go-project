## 此项目是参考阳明的<go训练营>课程的项目【kubeA】

## git相关操作
```bash
git clone https://github.com/hu417/go-project.git

// 初始化
git init
git config --global user.name "***"
git config --global user.email ****@qq.com

// 处理换行符
git config --global core.autocrlf true

// ssl认证关闭
git config --global http.sslVerify "false"
git config --global credential.helper manager

// 新建分支
git branch -M main
// 添加远程分支
git remote add origin https://github.com/hu417/go-project.git

// 分支操作 
git checkout -b dev   # 新建分支并切换
git branch            # 查看当前分支
git branch -a         # 查看远程分支
git push --set-upstream origin dev # 建立本地分支和远程分支的关联（创建远程分支）
git add .

// 提交项目
echo "# kubeA项目" >> README.md
git add .
git commit -m "fix: kubeA项目
【kubeA-web】
1、初步完成layout布局
【kubeA-server】
1、初始化k8s-client配置
" 
git tag -a v0.1 -m "版本v0.1"
git push
git push --tags
git status
git log


// 根据tag创建分支
git origin fetch                        # 获得最新
# git branch <new-branch-name> <tag-name> # 会根据tag创建新的分支
git branch newbranch v1.0 . 会以tag v1.0创建新的分支newbranch。

git checkout newbranch                  # 切换到新的分支。
git push origin newbranch               # 把本地创建的分支提交到远程仓库

// 其他操作
git tag -d v0.1                 # 删除本地tag
git push --delete origin v0.3   # 删除远程tag
git checkout main               # 切换到main分支
git merge dev                   # 合并 dev 到 mian 分支
git reset HEAD                  # 撤销已经add，但是没有commit
git reset --soft HEAD^          # 撤销已经commit，但是没有push到远端的文件（仅撤销commit 保留add操作）
git reset --hard commit_id      # 回退到指定的commit id
# 撤销已经push到远端的文件
// 切换到指定分支
git checkout 分支名
// 撤回到需要的版本
git reset --soft 需要回退到的版本号
//提交撤销动作到服务器，强制提交当前版本号
git push origin 分支名 --force


// .gitignore文件
node_modules/
dist/


```

