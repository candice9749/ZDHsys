package salt

import (
    "bytes"
    "fmt"
    "strings"
    "io/ioutil"
    "net/http"
    "net/url"
    "crypto/tls"
    "encoding/json"
    "github.com/astaxie/beego"
    "github.com/bitly/go-simplejson"
)

//成员变量名首字母必须大写
type JsonU struct {
	Username string `json:"username"`        //重命名，因为前面的部分不能小写，所以用这种方法让他在输出的时候是小写的
	Password string `json:"password"`
	Eauth    string `json:"eauth"`
}

type Json struct {
        Perms  []string `json:"perms"`
        Start  float64  `json:"start"`
        Token  string   `json:"token"`
        Expire float64  `json:"expire"`
        User   string   `json:"user"`
        Eauth  string   `json:"eauth"`
}

type Jsonslice struct {
        Return []Json `json:"return"`   //类型是Json结构体的切片
}


/*
 * 返回token
 */
func token() string {
    salt_url := beego.AppConfig.String("salt_url")   //salt服务器的地址

    var js JsonU
    js.Username = beego.AppConfig.String("salt_username")
    js.Password = beego.AppConfig.String("salt_password")
    js.Eauth = "pam"  //身份认证

    b, err := json.Marshal(js) //将数据编码成json字符串
    if err != nil {
        fmt.Println("json err:", err)
    }
    fmt.Println(string(b))   //必须要转换成字符串输出才是正常的，不然就是一堆数字
    var jsonStr = b //salt的用户名和密码和认证方式

    req, err := http.NewRequest("POST", salt_url, bytes.NewBuffer(jsonStr))//返回一个新的request
    req.Header.Set("Accept", "application/json")
    req.Header.Set("Content-Type", "application/json")   //设置头信息

    tr := &http.Transport{  //结构体
        TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},//InsecureSkipVerify为真，TLS接受任何证书
        DisableCompression: true, //禁止传输
    }

    client := &http.Client{Transport: tr}

    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)

    var s Jsonslice
    str := string(body)
    json.Unmarshal([]byte(str), &s)

    var token string
    for _, v := range s.Return {
        token = v.Token
    }
    fmt.Println(token)
    return token
}


/*
 * 公共POST传递func
 */
func exec(data string, accept string, ctype string) string {
    token  := token()
    salt_api_url := beego.AppConfig.String("salt_api_url") //连接
    req, err := http.NewRequest("POST", salt_api_url, strings.NewReader(data))
    req.Header.Set("Accept", accept)
    req.Header.Set("X-Auth-Token", token)
    req.Header.Set("Content-Type", ctype)
    tr := &http.Transport{
        TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
        DisableCompression: true,
    }
    client := &http.Client{Transport: tr}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    return string(body)
}


/*
 * 重新组合tgt
 */
func resetTgt(tgt string) string {
    tmp := strings.Split(tgt, "\n")  //以换行符分隔tgt
    var str []string
    for i := 0; i < len(tmp); i++ {
        if len(tmp[i]) != 0 {
            str = append(str, strings.TrimSpace(tmp[i]))  //tgt的前后端空白都去掉，然后追加
        }
    }
    tgtall := strings.Join(str,",") //将切片组合成字符串，并用，分隔
    return tgtall
}

//minion端的认证管理
type KeyListController struct {
    beego.Controller
}

func (this *KeyListController) Get() {
    token := token()
    fmt.Println(token)

    var para = url.Values{}
    para.Add("client", "wheel")
    para.Add("tgt", "*")
    para.Add("fun", "key.list_all")
    data := para.Encode()
    fmt.Println(data)

    salt_api_url := beego.AppConfig.String("salt_api_url")
    req, err := http.NewRequest("POST", salt_api_url, strings.NewReader(data))
    req.Header.Set("Accept", "application/json")
    req.Header.Set("X-Auth-Token", token)
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    tr := &http.Transport{
        TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
        DisableCompression: true,
    }
    client := &http.Client{Transport: tr}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))

    js, err := simplejson.NewJson(body)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(js)
    pjs         := js.Get("return").GetIndex(0).Get("data").Get("return").Get("minions").MustStringArray()
    minions_pre := js.Get("return").GetIndex(0).Get("data").Get("return").Get("minions_pre").MustStringArray()

    var str []string
    for i := 0; i < len(pjs); i++ {
        str = append(str, pjs[i])
    }
    out := strings.Join(str,",")
    fmt.Println(out)

    this.Data["title"]               = beego.AppConfig.String("title")
    this.Data["culture"]             = beego.AppConfig.String("culture")
    this.Data["navtitle"]            = "minion管理"
    this.Data["salt_active"]         = "active"
    this.Data["salt_active_keylist"] = "active"
    this.Data["minions"]             = pjs
    this.Data["minions_pre"]         = minions_pre
    this.Layout = "admin/layout.tpl"
    this.TplName = "salt/keylist.tpl"
}


type KeyDeleteController struct {
    beego.Controller
}

func (this *KeyDeleteController) Get() {
    match := this.GetString("match")
    token := token()
    fmt.Println(token)

    var para = url.Values{}
    para.Add("client", "wheel")
    para.Add("match", match)
    para.Add("fun", "key.delete")
    data := para.Encode()
    fmt.Println(data)

    salt_api_url := beego.AppConfig.String("salt_api_url")
    req, err := http.NewRequest("POST", salt_api_url, strings.NewReader(data))
    req.Header.Set("Accept", "application/json")
    req.Header.Set("X-Auth-Token", token)
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    tr := &http.Transport{
        TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
        DisableCompression: true,
    }
    client := &http.Client{Transport: tr}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))

    this.Ctx.Redirect(302, "/salt/keylist")
}


type KeyAcceptController struct {
    beego.Controller
}

func (this *KeyAcceptController) Get() {
    match := this.GetString("match")
    token := token()
    fmt.Println(token)

    var para = url.Values{}
    para.Add("client", "wheel")
    para.Add("match", match)
    para.Add("fun", "key.accept")
    data := para.Encode()
    fmt.Println(data)

    salt_api_url := beego.AppConfig.String("salt_api_url")
    req, err := http.NewRequest("POST", salt_api_url, strings.NewReader(data))
    req.Header.Set("Accept", "application/json")
    req.Header.Set("X-Auth-Token", token)
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    tr := &http.Transport{
        TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
        DisableCompression: true,
    }
    client := &http.Client{Transport: tr}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))

    this.Ctx.Redirect(302, "/salt/keylist")
}


/*
 * test.ping 联机状态检测
 */
type PingController struct {
    beego.Controller
}

func (this *PingController) Get() { //表示刚打开页面的时候显示的
    token := token()
    fmt.Println(token)

    var para = url.Values{}
    para.Add("client", "local")
    para.Add("tgt", "*")
    para.Add("fun", "test.ping")
   // para.Add("fun", "grains.item")
   // para.Add("arg", "os")
    fmt.Println("arg:",para)
    data := para.Encode()
    fmt.Println(data)

    body := exec(data, "application/x-yaml", "application/x-www-form-urlencoded")
    this.Data["title"] = beego.AppConfig.String("title")
    this.Data["culture"] = beego.AppConfig.String("culture")
    this.Data["navtitle"] = "联机状态"
    this.Data["testping"] = string(body)
    this.Layout = "admin/layout.tpl"
    this.TplName = "salt/ping.tpl"
}

func (this *PingController) Post() {   //表示提交tgt后显示的
    tgt := this.GetString("tgt")
    if tgt == "*" {
        fmt.Println("tgt = *")
        this.Ctx.Redirect(302, "/salt/ping")
    }

    fmt.Println(resetTgt(tgt))          //输出的是重新组合后的tgt

    var para = url.Values{}
    para.Add("client", "local")
    para.Add("tgt", resetTgt(tgt))
    para.Add("fun", "test.ping")
    para.Add("expr_form", "list")
    data := para.Encode()
    body := exec(data, "application/x-yaml", "application/x-www-form-urlencoded")

    this.Data["title"] = beego.AppConfig.String("title")
    this.Data["culture"] = beego.AppConfig.String("culture")
    this.Data["navtitle"] = "联机状态"
    this.Data["tgt"] = tgt
    this.Data["testping"] = string(body)
    this.Layout = "admin/layout.tpl"
    this.TplName = "salt/ping.tpl"
}


/*
 *执行命令
 */
type CmdRunController struct {
    beego.Controller
}

func (this *CmdRunController) Get() {
    this.Data["title"]               = beego.AppConfig.String("title")
    this.Data["culture"]             = beego.AppConfig.String("culture")
    this.Data["navtitle"]            = "执行命令"
    this.Layout = "admin/layout.tpl"
    this.TplName = "salt/cmdrun.tpl"
}


type jj struct {
    Return []map[string]string         //一个map类型的切片
}

func (this *CmdRunController) Post() {
    tgt := this.GetString("tgt")          //主机列表
    tgtall:=resetTgt(tgt)              //重新组合后的主机列表

    arg    := this.GetString("arg")
    fmt.Println(tgtall)

    var para = url.Values{}
    para.Add("client", "local")
    para.Add("tgt", tgtall)
    para.Add("fun", "cmd.run")
    para.Add("expr_form", "list")
    para.Add("arg", arg)
    data := para.Encode()

    body := exec(data, "application/json", "application/x-www-form-urlencoded")

    js, _ := simplejson.NewJson([]byte(body))            //处理一个很复杂的json的，一层一层嵌套的那种
    fmt.Println(body)
    fmt.Println("response Body:", js.Get("return").GetIndex(0)) //返回return下的数组


    gjs := make(map[string]interface{})
    gjs, _  = js.Get("return").GetIndex(0).Map()
    /*if err != nil {
        fmt.Println(err)
    }*/
    fmt.Println("simplejson: ",gjs)

    var ret bytes.Buffer
    var j jj
    json.Unmarshal([]byte(body), &j)
    for _, i := range j.Return{
        for k, v := range i{
            fmt.Println(k,v)
            ret.WriteString(fmt.Sprintf("【%s】\n",k))
            ret.WriteString(v)
            ret.WriteString(fmt.Sprintf("\n------------------------------\n"))
    }
}

    this.Data["title"]               = beego.AppConfig.String("title")
    this.Data["culture"]             = beego.AppConfig.String("culture")
    this.Data["navtitle"]            = "执行命令"
    this.Data["tgt"]                 = tgt
    this.Data["arg"]                 = arg
    this.Data["return"]              = ret.String()
    this.Layout = "admin/layout.tpl"
    this.TplName = "salt/cmdrun.tpl"
}


type CpGetFileController struct {
    beego.Controller
}

/*
 * curl -k https://127.0.0.1:8000 -H "Accept: application/x-yaml" -H "X-Auth-Token: 56120d270e90588cbb4d35d1ca0e2d1ed679e650" -d client='local' -d tgt='*' -d fun='cp.get_file' -d arg='salt://hello2' -d arg='/tmp/hello2'
 */
func (this *CpGetFileController) Get() {
    this.Data["title"]               = beego.AppConfig.String("title")
    this.Data["culture"]             = beego.AppConfig.String("culture")
    this.Data["navtitle"]            = "分发脚本/文件"
    this.Layout = "admin/layout.tpl"
    this.TplName = "salt/cpgetfile.tpl"
}

func (this *CpGetFileController) Post() {
    tgt := this.GetString("tgt")
    tmp := strings.Split(tgt, "\n")
    var str []string
    for i := 0; i < len(tmp); i++ {
        if len(tmp[i]) != 0 {
            str = append(str, strings.TrimSpace(tmp[i]))
        }
    }
    tgtall := strings.Join(str,",")
    sarg    := this.GetString("sarg")
    darg    := this.GetString("darg")
    fmt.Println(tgtall)
    fmt.Println(sarg)
    fmt.Println(darg)
    var para = url.Values{}
    para.Add("client", "local")
    para.Add("tgt", tgtall)
    para.Add("fun", "cp.get_file")
    para.Add("expr_form", "list")
    para.Add("arg", sarg)
    para.Add("arg", darg)
    data := para.Encode()
    fmt.Println(data)
    body := exec(data, "application/x-yaml", "application/x-www-form-urlencoded")
    this.Data["title"]               = beego.AppConfig.String("title")
    this.Data["culture"]             = beego.AppConfig.String("culture")
    this.Data["navtitle"]            = "分发脚本/文件"
    this.Data["sarg"]                = sarg
    this.Data["darg"]                = darg
    this.Data["return"]              = body
    this.Layout = "admin/layout.tpl"
    this.TplName = "salt/cpgetfile.tpl"
}

//部署应用
type DeployController struct {
    beego.Controller
}

func (this *DeployController) Get() {
    this.Data["title"]               = beego.AppConfig.String("title")
    this.Data["culture"]             = beego.AppConfig.String("culture")
    this.Data["navtitle"]            = "部署应用"
    this.Layout = "admin/layout.tpl"
    this.TplName = "salt/deploy.tpl"
}

func (this *DeployController) Post() {
    tgt := this.GetString("tgt") //获取主机列表
    tmp := strings.Split(tgt, "\n")         //没一台主机，返回的是字符串数组
    var str []string
    for i := 0; i < len(tmp); i++ {        //循环每一个主机
        if len(tmp[i]) != 0 {
            str = append(str, strings.TrimSpace(tmp[i]))  //去掉空白，添加到str这个切片里面
        }
    }
    tgtall := strings.Join(str,",")  //将IP组成一个字符串，用，分隔

    arg    := this.GetString("app")  //获取数据
    if arg==""{
        arg    = this.GetString("app1")             //如果没从给定的软件里面选择的话就从文本框获取
    }
    fmt.Println("=================")
    fmt.Println(arg)
    //fmt.Println(tgtall)
    fmt.Println("==================")
    var para = url.Values{}
    para.Add("client", "local")
    para.Add("tgt", tgtall)
    para.Add("fun", "state.sls")
    para.Add("expr_form", "list")       //为列表传递参数
    para.Add("arg", arg)
    data := para.Encode()

    body := exec(data, "application/json", "application/x-www-form-urlencoded")
    js, _ := simplejson.NewJson([]byte(body))
    fmt.Println(body)
    fmt.Println("response Body:", js.Get("return").GetIndex(0))
    gjs := make(map[string]interface{})
    gjs, _  = js.Get("return").GetIndex(0).Map()
    /*if err != nil {
        fmt.Println(err)
    }*/
    fmt.Println("simplejson: ",gjs)

    var out bytes.Buffer
    err := json.Indent(&out, []byte(body), "", "    ")

    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(out.String())

    var ret bytes.Buffer
    var j jj
    json.Unmarshal([]byte(body), &j)
    for _, i := range j.Return{
        for k, v := range i{
            fmt.Println(k,v)
            ret.WriteString(fmt.Sprintf("【%s】\n",k))
            ret.WriteString(v)
            ret.WriteString(fmt.Sprintf("\n------------------------------\n"))
        }
    }

    this.Data["title"]               = beego.AppConfig.String("title")
    this.Data["culture"]             = beego.AppConfig.String("culture")
    this.Data["navtitle"]            = "部署应用执行结果"
    this.Data["tgt"]                 = tgt
    this.Data["return"]              = out.String() 
    this.Layout = "admin/layout.tpl"
    this.TplName = "salt/deploy.tpl"
}
