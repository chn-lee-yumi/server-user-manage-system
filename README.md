# Server User Manage System - 服务器用户管理系统

## 介绍

这是一个很简单的系统，实现了Linux用户管理的功能。制作这个项目的背景是学校的社团有一台服务器，为了安全，设置成了密钥登录，因此每个人都需要登记一下公钥。管理员手动登记太麻烦，所以写一个自助登记公钥的系统。

## 功能

- 添加用户：创建用户并登记公钥
- 删除用户（TODO）

## 部署

新建目录`/root/sums`，复制`sums`和`index.html`到`/root/sums`。  

复制`sums.service`到`/lib/systemd/system`。  

启动：`systemctl start sums`  

## 使用

端口号为8000，直接访问就ok。  
管理员生成凭证的网址：`/api/newticket`。没有任何认证，网址需要**妥善保管**。可以修改源码得到一个更加复杂的网址。

## 注意事项

需要root权限，故当前部署采用root用户运行。
