package ssh

import (
    "fmt"
     _ "io"
    "bytes"
    "strings"
    "time"
    "golang.org/x/crypto/ssh"
    "github.com/astaxie/beego"
    "sync"
)

var j = 0
var m map[string]string

//ssh主机列表
type SshController struct {
    beego.Controller
}

func (this *SshController) Get() {

    this.Data["title"] = beego.AppConfig.String("title")
    this.Data["culture"] = beego.AppConfig.String("culture")
    this.Data["navtitle"] = "主机列表"
    this.Layout = "admin/layout.tpl"
    this.TplName = "ssh/index.tpl"
}

//sync.WaitGroup是为了阻塞goroutine的，只有等到goroutine都执行完毕，才将主函数停止。
func ssh_cmd(ip_port, user, password, ip, cmd string, wg *sync.WaitGroup) {
    defer wg.Done()
    //ClientConfig是包ssh里面的一个结构体，AuthMethod是一个接口，里面有ssh.password这个函数，传递的值是password
    Conf := ssh.ClientConfig{User: user, Auth: []ssh.AuthMethod{ssh.Password(password)},HostKeyCallback: ssh.InsecureIgnoreHostKey()}
    //ssh连接客户端
    Client, err := ssh.Dial("tcp", ip_port, &Conf)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer Client.Close()
    //NewSession为该客户端打开一个新会话。(会话是程序的远程执行。)
    if session, err := Client.NewSession(); err == nil {
        defer session.Close()
        j = j+1
        //CombinedOutput在远程主机上运行cmd，并返回其组合的标准输出和标准错误。
        buf, _ := session.CombinedOutput(cmd)
        m[ip] = string(buf)
    }
}


func (this *SshController) Post(){
    //取数据
    iplist := this.GetString("iplist")
    //想要执行命令的主机个数
    tmp := strings.Split(iplist, "\n")
    cmd := this.GetString("cmd")
    m = make(map[string]string)
    start := time.Now()

    wg := sync.WaitGroup{}
    wg.Add(len(tmp))

    for i := 0; i < len(tmp); i++ {
        if len(tmp[i]) != 0 {
            split := strings.Split(tmp[i], ":")

            var ip_port string
            ip_port   = fmt.Sprintf("%s:%s",split[0],split[1])
            user     := split[2]
            password := strings.TrimSpace(split[3])
        
            go ssh_cmd(ip_port, user, password, split[0], cmd, &wg)
        }
    }
    wg.Wait()
    fmt.Println(len(tmp))
    //定义一个buffer缓冲
    var ret bytes.Buffer
    for k, v := range m {
        fmt.Println(k,":",v)
        ret.WriteString(fmt.Sprintf("%s\n%s", k, v))
        ret.WriteString(fmt.Sprintf("------------------------------\n"))
    }
    
    runtime := time.Now().Sub(start).Seconds()
    get_j := j
    j = 0

    //从app.conf文件里面获取内容
    this.Data["title"] = beego.AppConfig.String("title")
    this.Data["culture"] = beego.AppConfig.String("culture")
    this.Data["navtitle"] = "执行命令"
    this.Data["iplist"] = iplist
    this.Data["rcmd"] = ret.String()
    this.Data["tips"] = fmt.Sprintf("%s 主机数:%d 执行时间:%ds", cmd, get_j, int(runtime))
    this.Data["runtime"] = runtime
    this.Layout = "admin/layout.tpl"
    this.TplName = "ssh/index.tpl"
}

