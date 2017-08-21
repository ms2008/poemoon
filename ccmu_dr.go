// CCMU Automatic Connector
//
// This tool provides automatic login feature for
// connect to the network authenticated by dr.com
//
// Author: ms2008vip@gmail.com at 2017/8/15 11:10:02


//go:generate goversioninfo -icon=cover.ico

package main

import (
    "fmt"
    _ "net/http"
    _ "net/url"
    "os"
    "bufio"
    _ "io/ioutil"
    "strings"
    _ "strconv"
    _ "math/rand"
    "time"
    "flag"
    _ "github.com/toqueteos/webbrowser"

    "./conf"
    "./utils"
    "./service"
)

const (
    URL = "http://192.168.161.2/"
    FILE = "ids.txt"
)

var buildstamp string = ""
var githash string = ""


func main() {
    defer log.Close()

    args := os.Args
    if len(args)==2 && (args[1]=="--version" || args[1] =="-v") {
        fmt.Printf("Git Commit Hash: %s\n", githash)
        fmt.Printf("UTC Build Time: %s\n", buildstamp)
        return
    }

	var err error

	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}

    // 1. challenge
    svr := service.New(conf.Conf)
	if err = svr.Challenge(svr.ChallengeTimes); err != nil {
		log.Error("drcomSvc.Challenge(%d) error(%v)", svr.ChallengeTimes, err)
		return
	}

	// 2. login
	if err = svr.Login(); err != nil {
		log.Error("drcomSvc.Login() error(%v)", err)
		return
	}

	// 3. keepalive
	count := 0
	for {
		count++
		if err = svr.Alive(); err != nil {
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
