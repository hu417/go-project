
## 分支说明
```bash
开发分支简要:
    DEV分支: 开发环境分支
    TEST分支: 测试环境分支
    UAT分支: 联调环境分支
    PRE分支: 预发布环境分支
    MASTER分支: 生产环境分支

代码合并路线是: DEV->TEST->UAT->PRE->MASTER 然后根据不同的分支判断执行不同环境的部署。
```

## 项目准备
```bash
node        v18.12.1  `npm config set registry="https://registry.npm.taobao.org/"`
npm         8.19.2
vue3        @vue/cli 4.5.12  `npm install vue@3.2.47 -g`
vue/cli     @4.5.12          `npm install @vue/cli@4.5.12 -g`
vscode      插件：Auto Close Tag、Auto Rename Tag、Live Server、Chinese (Simplified)、Vue Language Features (Volar)、Vue 3 Snippets
```

## 创建项目
```bash
$ npm init vite@latest blog_client
√ Select a framework: » Vue
√ Select a variant: » JavaScript

Scaffolding project in D:\study\project\golang\3-gin\go-vue-blogCommunity\blog_client\blog_client...

Done. Now run:

  cd blog_client
  npm install     // 安装依赖
  npm run dev     // 运行项目

// nodejs版本高,则需要: NODE_OPTIONS="--openssl-legacy-provider"
```

## 项目初始化
```bash
// Vite构建工具
官方文档https://cn.vitejs.dev/guide/
安装步骤:
npm config get registry
npm config set registry=https://registry.npm.taobao.org/
npm install vite@4.3.1

创建项目:
npm create vite@latest kubeA-web -- --template vue
  cd kubeA-web
  npm install     // 安装依赖
  npm run dev     // 运行项目

# nodejs版本高,则需要: NODE_OPTIONS="--openssl-legacy-provider"


// vue-router: Vue.js官方的路由插件
npm install vue-router@4

// nprogress 进度条
进度条: https://www.npmjs.com/package/nprogress
安装:
npm install nprogress

// ant-design-vue 组件库
npm install ant-design-vue --save


// axios: 基于promise的HTTP库，用于http请求
官方文档: https://router.vuejs.org/zh/introduction.html
npm install axios



```

## 修改全局格式文件
```bash
// src/style.css
# 全部注释,添加此段
body {
  background-color: #FCFAF7;
  margin: 0;
  padding: 0;
}
```

## git操作
```bash
git clone https://github.com/hu417/Vue3-Vite-ElementPlus.git

// 初始化
git init
git config --global user.name "***"
git config --global user.email ****@qq.com

// 处理换行符
git config --global core.autocrlf true

// ssl认证关闭
git config --global http.sslVerify "false"
git config --global credential.helper manager

// 创建分支
git branch -M main
git remote add origin https://github.com/hu417/Vue3-Vite-ElementPlus.git

// 分支操作 
git checkout -b dev   # 新建分支并切换
git branch            # 查看当前分支
git branch -a         # 查看远程分支
git push --set-upstream origin dev # 建立本地分支和远程分支的关联（创建远程分支）


// 提交项目
echo "# Vue3-Vite-ElementPlus" >> README.md
git add .
git commit -m "fix: Vue3-Vite-ElementPlus项目
【后台全局Layout布局开发】
1、后台页面总布局实现-Container 布局容器
" 
git tag -a v1.0 -m "版本v1.0"
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
git tag -d v0.3                 # 删除本地tag
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