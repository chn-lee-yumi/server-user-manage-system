# Server User Manage System - 服务器用户管理系统

## 部署

新建目录`/root/sums`，复制`sums`和`index.html`到`/root/sums`。  

复制`sums.service`到`/lib/systemd/system`。  

启动：`systemctl start sums`  

## 注意事项

需要root权限，故当前部署采用root用户运行。
