# GoJudge
基于Go语言开发的判题机以及调度器

## 部署

依次构建变异环境、判题机、以及调度服务器。当构建完成后会监听某自定义端口(默认7070)

### 编译环境

ComplieEnv中定义了编译、运行所需的基本环境.可通过修改Dockerfile进行自定义配置。通过运行build.sh进行构建

### 判题机

JudgeCore中的config.json定义了编译选项，可以进行自定义配置。通过运行build.sh进行构建.

### 调度服务器

JudgeServer中的build.sh中配置了调度器暴露的端口，可以进行修改。但其余部分用户不应进行更改。通过运行build.sh进行构建

## 使用

通过发送对应判题信息到listenaddress/submit_task进行操作



