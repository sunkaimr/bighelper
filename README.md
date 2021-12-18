# bighelper

## 功能

通过连接[贝壳物联平台](https://www.bigiot.net)实现远程关机等功能，还可以将[贝壳物联](https://www.bigiot.net)账号绑定到天猫精良、小爱同学、小度等实现语音远程注销、关机、重启电脑等操作。

*支持windowns系统和linux系统（树莓派等）*

## 代码结构
```bash
bighelper
├── LICENSE
├── README.md
├── linux                  # linux下使用此下面代码
│   ├── bighelper.ini       # 配置文件
│   ├── go.mod
│   └── main.go            # 主程序
└── win                    # windows下使用此下面代码
    ├── bighelper.ini      # 配置文件
    ├── go.mod
    └── main.go            # 主程序             
```

## 编译

这里提供在linux环境下编译不同系统或架构的方法，在golang、idea等IDE环境下的编译方法请自行查找。

### 编译用于在windows下运行二进制

```bash
# 切换至win目录下
cd bighelper/win

# 设置编译环境，并编译
GOOS=windows GOARCH=amd64 go build

# 编译后生成可执行文件bighelper.exe
```

### 编译用于在linux下运行的程序

```bash
# 切换至linux目录下
cd bighelper/linux

# 设置编译环境，并编译。GOARCH设置为amd64
GOOS=linux GOARCH=amd64 go build

# 编译后生成可执行文件bighelper
```

### 编译用于在树莓派下运行的程序

```bash
# 切换至linux目录下
cd bighelper/linux

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
```

## 安装

### windows系统上安装

#### 以服务方式运行

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

- 利用windows的sc工具将应用发布成服务

```bash
# 1. 在windows左下角搜索框输入“cmd”，找到“命令提示符”程序，以管理员身份运行
#    如果不知道怎么操作请自行百度解决，但是“命令提示符”程序务必是以管理员身份启动，否则接下来安装服务会失败

# 2. 安装服务。其中“binpath”的值填写自己bighelper.exe安装的绝对路径。注意“=”后边有一个空格
sc create bighelper binpath= "C:\Program Files\bighepler\bighelper.exe" start= auto displayname= "bighelper"

# 3. 启动服务
net start bighelper

# 其他命令
net stop bighelper      # 停止服务
sc  delete bighelper    # 卸载服务
```

- 确认服务安装成功及运行状态

```bash
# 方式一
# 打开“任务管理器” > 找到“服务标签”，查看服务的状态，启动或者关闭服务。


# 方式二
# 桌面找到“此电脑” > 鼠标右键单击“此电脑” > 选择"管理" > 点击“服务和应用程序” > 点击“服务”，在服务列表中根据服务名称“bighepler”可以查看服务的状态，启动或者关闭服务。
```

#### 以普通应用方式运行

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

### windows支持的命令

- 命令：logoff
```bash
功能：注销登录
限制: 只能在用户登录后启动的应用生效;如果有其他软件阻止，则无法注销
      比如“以普通应用方式运行”才能生效，“以服务方式运行”时此命令无效，因为服务不属于任何用户所以不存在注销一说
```

- 命令：forcelogoff
```bash
功能：强制注销登录
功能：注销登录
限制: 只能在用户登录后启动的应用生效;可能导致未保存的文件丢失。
      比如“以普通应用方式运行”才能生效，“以服务方式运行”时此命令无效，因为服务不属于任何用户所以不存在注销一说
```

- 命令：shutdown
```bash
功能：关机
限制: 所有文件都已写入磁盘，所有软件都已关闭。如果有其他软件阻止，则无法关闭
```
- 命令：forceshutdown
```bash
功能：强制关机
限制: 可能导致未保存的文件丢失
```
- 命令：reboot
```bash
功能：重启
限制: 所有文件都已写入磁盘，所有软件都以关闭。 如果有其他软件阻止，则无法重启
```

- 命令：forcereboot
```bash
功能：强制重启
限制: 可能导致未保存的文件丢失
```

### linux下支持的命令

- 命令：shutdown

```bash
功能：关机
说明：相当于执行命令“shutdown -H -t 15”，收到命令后15秒后开始关机
```

- 命令：reboot

```bash
功能：重启
说明：相当于执行命令“shutdown -r -t 15”，收到命令后15秒后开始重启
```

- 命令：cancel

```bash
功能：关机
说明：相当于执行命令“shutdown -c”，当发出“shutdown”和“reboot”后15秒内发送该命令可以取消关机或重启
```

