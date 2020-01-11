package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

var (
	isDecode    *bool
	file        *os.File
	base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()abcdefghjklmnopqrs"
)

func main() {
	isDecode = flag.Bool("D", false, "decode text")
	port := flag.String("p", "", "port")
	output := flag.String("o", "", "output to file")
	decodeFile := flag.String("df", "", "decode file")
	flag.Parse()
	if *output != "" {
		var err error
		file, err = os.OpenFile(*output, os.O_APPEND|os.O_CREATE|os.O_RDWR|os.O_SYNC, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()
	}
	// 解密文件
	if *decodeFile != "" {
		srcFile, err := os.Open(*decodeFile)
		if err != nil {
			panic(err)
		}
		all, err := ioutil.ReadAll(srcFile)
		if err != nil {
			panic(err)
		}
		for _, s := range strings.Split(string(all), "\t\r\n") {
			if file == nil {
				fmt.Println(decode(s))
			} else {
				file.WriteString(decode(s))
			}
		}
		if file != nil {
			file.Sync()
		}
	}
	if *port != "" {
		ln, err := net.Listen("tcp", *port)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Listen to  %s\n", *port)
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Println("FROM: ", conn.RemoteAddr().String())
			go handleConnection(conn)
		}
	}
	if *port == "" && *decodeFile == "" {
		flag.Usage()
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	res := make([]byte, 0)
	buffer := make([]byte, 1024)
	var count int
	for {
		length, err := conn.Read(buffer)
		if err != nil {
			log.Println(err)
			break
		}
		count += length
		res = append(res, buffer[:length]...)
		if string(res[(count-3):count]) == "\t\r\n" { // 协议 tcp包 结束标识
			content := string(res[:count])
			if *isDecode {
				cs := strings.Split(content, "\t\r\n")
				for _, c := range cs {
					c = fmt.Sprintf("%s\t\r\n", decode(c))
					if file != nil {
						file.WriteString(c)
					}
					fmt.Print(c)
				}
			} else {
				if file != nil {
					file.WriteString(content)
				}
				fmt.Print(content)
			}
			if file != nil {
				file.Sync()
			}
			// 重置
			res = res[:0]
			count = 0
		}
	}
}

// 解密文本
func decode(s string) string {
	coder := base64.NewEncoding(base64Table)
	res, _ := coder.DecodeString(s)
	return string(res)
}
