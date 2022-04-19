# bighelper

## 功能

通过连接[贝壳物联平台](https://www.bigiot.net)实现远程关机等功能，还可以将[贝壳物联](https://www.bigiot.net)账号绑定到天猫精良、小爱同学、小度等实现语音远程对电脑进行注销、关机、重启、休眠等操作还支持自定义指令。

*支持windowns系统和linux系统（树莓派等）*

## 代码结构
```bash
bighelper
├── action
│   ├── driver
│   │   ├── driver_linux.go
│   │   └── driver_windows.go
│   └── action.go
├── bigiot
│   └── bitgiot.go
├── bin
├── config
│   └── config.go
├── service
│   └── service.go
├── LICENSE
├── README.md
├── bighelper.ini  # 配置文件
├── go.mod
├── go.sum
└── main.go        # 程序入口   
```

## 编译

这里提供在linux环境下编译不同系统或架构的方法，在golang、idea等IDE环境下的编译方法请自行查找。

### 编译用于在windows下运行二进制

```bash
cd bighelper

# 设置编译环境，并编译
GOOS=windows GOARCH=amd64 go build

# 编译后生成可执行文件bighelper.exe
```

### 编译用于在linux下运行的程序

```bash
cd bighelper

# 设置编译环境，并编译。GOARCH设置为amd64
GOOS=linux GOARCH=amd64 go build

# 编译后生成可执行文件bighelper
```

### 编译用于在树莓派下运行的程序

```bash
cd bighelper

# 设置编译环境，并编译。GOARCH设置为arm
GOOS=linux GOARCH=arm go build

# 编译后生成可执行文件bighelper
# 如果编译后在树莓派系统下运行失败可能是架构不兼容引起。可以将源码下载到树莓派上，在树莓派上安装go编译环境重新编译。
```

## 准备bigiot账号

- 已注册[贝壳物联](https://www.bigiot.net)账号并创建好产品。在"智能设备">"设备列表"中创建或选择一个合适的设备并记录其ID和APIKEY
- 修改bighelper.ini配置文件

```ini
# 设备ID，修改为自己的实际值
device_id = 12345

# APIKey，修改为自己的实际值
api_key   = 1s2f3h4k5

# 其他配置可根据自己需要更改
```

## 安装

### windows系统上安装

#### 以服务方式运行（建议）

在windows下支持将应用安装为服务（服务和普通程序的区别请自行了解），随系统自动启动。

- 编译生成bighelper.exe可执行文件
- 安装应用

```bash
# 1.创建一个合适的目录将应用放置进去
# 比如创建一个目录
C:\Program Files\bighepler

# 2. 将编译好的bighelper.exe可执行文件和bighelper.ini配置文件放置进去
#    bighelper.ini配置文件中的device_id和APIKey需要修改为自己的值
```

- 以管理员方式执行`install.bat`批处理指令

```bash
# 安装执行install.bat批处理

# 卸载执行uninstall.bat批处理
```

- 确认服务安装成功及运行状态

```bash
# 方式一
# 打开“任务管理器” > 找到“服务标签”，查看服务的状态，启动或者关闭服务。


# 方式二
# 桌面找到“此电脑” > 鼠标右键单击“此电脑” > 选择"管理" > 点击“服务和应用程序” > 点击“服务”，在服务列表中根据服务名称“bighepler”可以查看服务的状态，启动或者关闭服务。
```

#### 以普通应用方式运行（不建议）

- 将bighelper.exe可执行文件和bighelper.ini配置文件放置到合适的位置，比如新建目录"C:\Program Files\bighepler"放置进去
- bighelper.ini配置文件中的device_id和APIKey需要修改为自己的值
- 双击bighelper.exe即可运行（这种方式会有黑色的运行窗口，并输出运行日志）

如果想开机后应用自动启动可通过以下方式

```bash
# 1，“win + R” 键打开“运行”窗口，输入“shell:startup”

# 2，此时会自动打开一个文件夹，将bighelper.exe可执行文件和bighelper.ini配置文件放到此文件夹中

# 3，下次开机待用户登录完成后应用会自动运行

# 限制：
# 1，必须等待用户登录后才会启动
# 2，会有黑色窗口驻留，并显示日志
```

### linux系统上安装

- 下载安装包bighelper-linux.tar

```bash
# 1. 解压安装包
tar -xvf bighelper-linux.tar

# 2. 安装包包含以下文件 
bighelper
├── bighelper           # 主程序，如果想升级可直接将编译好的二进制替换掉
├── bighelper.ini       # 配置文件
├── bighelper.service   # service文件
├── install.sh          # 安装、卸载脚本
└── readme.txt

# 修改bighelper.ini配置文件
# device_id和APIKey需要修改为自己的值
```

- 运行`install.sh`安装脚本

```bash
# 安装命令
./install.sh install

# 设置开机自启动
systemctl enable bighelper

# 启动服务
systemctl start bighelper

# 查看服务状态
systemctl status bighelper
```

- 其他命令，可参考`readme.txt`

```
# 安装命令
./install.sh install

# 卸载命令
./install.sh uninstall

# 设置开机自启动
systemctl enable bighelper

# 关闭开机自启动
systemctl disable bighelper

# 启动服务
systemctl start bighelper

# 停止服务
systemctl stop bighelper

# 查看服务状态
systemctl status bighelper

# 查看服务日志
journalctl -u bighelper

# 实时查看服务日志
journalctl -fu bighelper
```

## 控制指令
通过[贝壳物联](https://www.bigiot.net)网页端或者贝壳物联公众号向设备发送指令控制设备。支持以下指令

### 内置的指令

- 命令：shutdown

```bash
# 功   能：一分钟后关机，可在一分钟内取消执行
# 限   制: 可能导致未保存的文件丢失
# linux ：shutdown -h 1
# windows：shutdown -s -t 60
```

- 命令：reboot

```bash
# 功   能：一分钟后重启，可在一分钟内取消执行
# 限  制: 可能导致未保存的文件丢失
# linux：shutdown -r 1
# windows：shutdown -r -t 60
```

- 命令：sleep

```bash
# 功  能：休眠
# 限  制: 仅windowns生效
# windows：shutdown -h
```

- 命令：cancel

```bash
功能：关机
说明：相当于执行命令“shutdown -c”，当发出“shutdown”和“reboot”后15秒内发送该命令可以取消关机或重启
```

### 指令别名

作用是给内置的命令起一个别名，以达到和内置指令相同效果的目的。

举2个例子：

1，贝壳物联是没有`shutdown`指令，但是有`stop`指令，可以将`stop`作为内置指令`shutdown`的别名。

2，贝壳物联有一个指令是`pause`,可以将`pause`作为内置指令`sleep`的别名。具体实现方式如下：

```bash
# 修改bighelper.ini的配置文件，找到`[alias]`配置块

[alias]
# 执行stop和执行shutdown指令有相同的效果
shutdown = stop
sleep = pause
```

### 自定义指令

```bash
[command]
# 自定义命令
# 功  能：执行自定义指令
custom = "touch /test.txt"        # 执行指令
bash = "bash /home/example.sh"    # 执行脚本
```

