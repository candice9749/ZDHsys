package server

import (
	"ZDHsys/models"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/sbinet/go-python"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

//主机总体信息页面
type InfoController struct {
	beego.Controller
}

func (this *InfoController) Get() {
	//有orm对象
	o := orm.NewOrm()
	//查询对象
	var servers []models.Server                    //定义一个结构体数组
	_, err := o.QueryTable("Server").All(&servers) //高级查询
	if err != nil {
		beego.Info("查询所有信息失败")
		return
	}
	//传给视图
	this.Data["servers"] = servers
	this.Data["navtitle"] = "主机信息"
	this.Data["title"] = "ZDHYW平台"
	this.Layout = "admin/layout.tpl"
	this.TplName = "server/srvinfo.html"
}


//现有机器页面
type ExistController struct {
	beego.Controller
}

func (this *ExistController) Get() {
	//有orm对象
	o := orm.NewOrm()
	//查询对象
	var servers []models.Server                    //定义一个结构体数组
	_, err := o.QueryTable("Server").All(&servers) //高级查询  //Server是表名
	if err != nil {
		beego.Info("查询所有信息失败")
		return
	}
	this.Data["servers"] = servers
	this.Data["navtitle"] = "现有机器"
	this.Data["title"] = "ZDHYW平台"
	this.Layout = "admin/layout.tpl"
	this.TplName = "server/srvexist.html"
}


//手动录入机器页面
type InputController struct {
	beego.Controller
}

func (this *InputController) Get() {
	this.Data["navtitle"] = "现有机器"
	this.Data["title"] = "ZDHYW平台"
	this.Layout = "admin/layout.tpl"
	this.TplName = "server/srvinput.html"
}

func (this *InputController) Post() {
	//1.拿到数据
	mac := this.GetString("mac")
	addr := this.GetString("addr")

	//2.判断mac地址是否存在
	o:=orm.NewOrm()
	server:=models.Server{}
	server.Mac=mac
	err:=o.Read(&server,"Mac")
	if err==nil{
		beego.Info("mac地址已存在！！！")
		this.Layout = "admin/layout.tpl"
		this.TplName = "server/error1.html"
	}else {
		//3.插入数据
		o := orm.NewOrm()         //申请orm对象
		server := models.Server{} //有一个要插入的结构体
		server.Mac = mac
		server.Addr = addr

		_, err := o.Insert(&server)
		if err != nil {
			beego.Info("插入数据库错误", err)
			return
		}

		//4.返回主机界面
		this.Layout = "admin/layout.tpl"
		this.Redirect("/server/srvexist.html", 302)
	}
}


//为现有机器创建system页面
type AddController struct {
	beego.Controller
}

func (this *AddController) Get() {
	//1.获取主机ip
	mac := this.GetString("mac") //注意要和前端页面里面的name是相同的
	//fmt.Println(mac)

	//2.查询数据库获取信息
	o := orm.NewOrm()
	server := models.Server{Mac: mac}
	err := o.Read(&server, "Mac")
	if err != nil {
		beego.Info("查询错误", err)
		return
	}

	//3.传递数据给视图
	this.Data["server"] = server
	this.Data["navtitle"] = "创建主机"
	this.Data["title"] = "ZDHYW平台"
	this.Layout = "admin/layout.tpl"
	this.TplName = "server/srvadd.html"
}

func (this *AddController) Post() {
	//1.拿到数据
	mac := this.GetString("mac")
	serverip := this.GetString("server_ip")
	os := this.GetString("os")
	servername := this.GetString("servername")
	line := this.GetString("line")
	ywman := this.GetString("ywman")


	//2.判断数据是否合法
	if serverip == "" || os == "" || servername == "" || mac == "" {
		beego.Info("添加主机数据错误")
		return
	}

	o := orm.NewOrm()
	server:=models.Server{}
	server.Serverip=serverip
	err:=o.Read(&server,"Serverip")
	if err==nil{
		beego.Info("ip地址已存在！！！")
		this.Layout = "admin/layout.tpl"
		this.TplName = "server/error2.html"
	}else {
		//3.更新操作
		o := orm.NewOrm()                 //申请orm对象
		server := models.Server{Mac: mac} //有一个要插入的结构体
		err := o.Read(&server, "Mac")     //注意！！！只要不是Id就必须加cols
		if err != nil {
			beego.Info("查询数据错误！")
			return
		}

		server.Serverip = serverip
		server.Os = os
		server.Hostname = servername
		server.Line = line
		server.Ywman = ywman

		_, err = o.Update(&server)
		if err != nil {
			beego.Info("更新数据显示错误")
			return
		}

		//调用cobbler

		//组合interfacestring
		interfaceString := fmt.Sprintf("{'macaddress-eth0':'%s','ipaddress-eth0':'%s','Gateway-eth0':'192.168.100.1','subnet-eth0':'255.255.255.0','static-eth0':1}", mac, serverip)
		fmt.Println(interfaceString)

		ks := "/var/lib/cobbler/kickstarts/" + os + ".ks"
		fmt.Println(ks)

		name:=mac+"-"+serverip+"-"+os
		Cobbler(name, servername, os, ks, interfaceString)

		//4.返回主机界面
		this.Layout = "admin/layout.tpl"
		this.Redirect("/server/srvinfo.html", 302)
	}
}

func Cobbler(name, hostname, profile, ks_meta, modify_interface string) {
	err := python.Initialize()
	if err != nil {
		panic(err.Error())
	}
	m := python.PyImport_ImportModule("sys")
	if m == nil {
		fmt.Println("import error")
		return
	}
	path := m.GetAttrString("path")
	if path == nil {
		fmt.Println("get path error")
		return
	}

	//导入路径，找python的包
	currentDir := python.PyString_FromString("")
	python.PyList_Insert(path, 0, currentDir)

	m = python.PyImport_ImportModule("cobbler")
	if m == nil {
		fmt.Println("import cobbler error")
		return
	}
	g := m.GetAttrString("cobbler")
	if g == nil {
		fmt.Println("get cobbler error")
		return
	}

	res := g.CallFunction(name, hostname, profile, ks_meta, modify_interface)
	if res == nil {
		fmt.Println("callfunction error")
		return
	}
}


//显示硬件信息---调用salt
type JsonU struct {
	Username string `json:"username"` //重命名
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
	Return []Json `json:"return"` //类型是Json结构体的切片
}

/*
* 返回token
 */
func token() string {
	salt_url := beego.AppConfig.String("salt_url") //salt服务器的地址

	var js JsonU
	js.Username = beego.AppConfig.String("salt_username")
	js.Password = beego.AppConfig.String("salt_password")
	js.Eauth = "pam" //身份认证

	b, err := json.Marshal(js) //将数据编码成json字符串
	if err != nil {
		fmt.Println("json err:", err)
	}
	fmt.Println(string(b))
	var jsonStr = b

	req, err := http.NewRequest("POST", salt_url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

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
	token := token()
	salt_api_url := beego.AppConfig.String("salt_api_url")
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

//对返回数据的处理
func handlebody(item string, server *models.Server) (a string) {
	var para = url.Values{}
	para.Add("client", "local")
	para.Add("tgt", "*")
	//para.Add("fun", "test.ping")
	para.Add("fun", "grains.item")
	para.Add("arg", item)
	fmt.Println("arg:", para)
	data := para.Encode()
	fmt.Println(data)

	body := exec(data, "application/x-yaml", "application/x-www-form-urlencoded")

	//对body进行切割
	body1 := strings.Trim(body, "return:") //将return去掉
	body2 := strings.TrimSpace(body1)      //将所有空白都去掉
	body3 := strings.Trim(body2, "-")      //把-去掉
	body3 = strings.Trim(body3, " ")       //把-去掉
	//body4 := strings.Replace(body3," ","",-1) //把全文中所有的空白去掉

	body4 := strings.Replace(body3, "\n", "", -1)
	body4 = strings.Replace(body4, "   ", "", -1)
	body5 := strings.Split(body4, "  ") //以os：为分隔线把他们都分开，返回切片
	//fmt.Println(body4)
	////
	//fmt.Println(body5)
	//fmt.Println("============")

	////对数据进行处理，根据不同的主机返回不同的信息
	servername := server.Hostname
	fmt.Print(servername)
	fmt.Println("============")
	for i, data := range body5 {
		fmt.Println(i, data)
	}
	//fmt.Printf("%T",body5[1])
	fmt.Println("============")
	for i := 0; i < len(body5); i++ {
		a := strings.HasPrefix(body5[i], servername)
		fmt.Println(a)
		if a == true {
			body6 := body5[i]
			body7 := strings.Split(body6, ":")
			body6 = body7[2]
			fmt.Println(body6)
			return body6
		}
	}
	return
}

//显示硬件信息的详情页
type HardController struct {
	beego.Controller
}

func (this *HardController) Get() {
	//1.获取主机ip
	serverip := this.GetString("serverip")

	//2.查询数据库获取信息
	o := orm.NewOrm()
	server := models.Server{Serverip: serverip}
	err := o.Read(&server, "Serverip")
	if err != nil {
		beego.Info("查询错误", err)
		return
	}
	//3.传递数据给视图

	//server.Hostname=handlebody("host",&server)
	server.Cpu = handlebody("cpu_model", &server)
	server.Os = handlebody("osfinger", &server)
	server.Memory = handlebody("mem_total", &server)
	Kernelversion := handlebody("kernelrelease", &server)
	this.Data["title"] = beego.AppConfig.String("title")
	this.Data["culture"] = beego.AppConfig.String("culture")
	//this.Data["server.Hostname"] = string(host)

	//fmt.Printf("数据类型为：%T",body)

	this.Data["server"] = server
	this.Data["Kernelversion"] = Kernelversion

	this.Data["navtitle"] = "硬件信息"
	this.Data["title"] = "ZDHYW平台"
	this.Layout = "admin/layout.tpl"
	this.TplName = "server/hardinfo.html"
}


//主机的编辑页面
type UpdateController struct {
	beego.Controller
}

func (this *UpdateController) Get() {
	//1.获取主机ip
	serverip := this.GetString("serverip") //注意要和前端页面里面的name是相同的
	//beego.Info("server ip id ",serverip)
	//2.查询数据库获取信息
	o := orm.NewOrm()
	server := models.Server{Serverip: serverip}
	err := o.Read(&server, "Serverip")
	if err != nil {
		beego.Info("查询错误", err)
		return
	}
	//3.传递数据给视图
	this.Data["server"] = server
	this.Data["title"] = "ZDHYW平台"
	this.Layout = "admin/layout.tpl"
	this.TplName = "server/srvupdate.html"
}

//处理更新页面的数据
func (this *UpdateController) Post() {
	//1.拿到数据
	id, _ := this.GetInt("id")
	mac := this.GetString("mac")
	serverip := this.GetString("server_ip")
	os := this.GetString("os")
	servername := this.GetString("servername")
	line := this.GetString("line")
	ywman := this.GetString("ywman")

	//2.对数据进行一个处理
	if serverip == "" || os == "" || servername == "" {
		beego.Info("添加主机数据错误")
		return
	}
	//3.更新操作
	o := orm.NewOrm()
	server := models.Server{Id: id}
	err := o.Read(&server)
	if err != nil {
		beego.Info("查询数据错误！")
		return
	}
	server.Serverip = serverip
	server.Os = os
	server.Hostname = servername
	server.Mac = mac
	server.Line = line
	server.Ywman = ywman

	_, err = o.Update(&server)
	if err != nil {
		beego.Info("更新数据显示错误")
		return
	}
	interfaceString := fmt.Sprintf("{'macaddress-eth0':'%s','ipaddress-eth0':'%s','Gateway-eth0':'192.168.100.1','subnet-eth0':'255.255.255.0','static-eth0':1, 'dnsname-eth0':'114.114.114.114'}", mac, serverip)
	fmt.Println(interfaceString)

	ks := "/var/lib/cobbler/kickstarts/" + os + ".ks"
	fmt.Println(ks)

	Cobbler(servername, servername, os, ks, interfaceString)

	//4.返回列表页面

	this.Layout = "admin/layout.tpl"
	this.Redirect("/server/srvinfo.html", 302)

}



//删除操作
type DeleteController struct {
	beego.Controller
}

func (this *DeleteController) Get() {
	//1.拿到数据
	serverip := this.GetString("serverip")
	//2.执行删除操作
	o := orm.NewOrm()
	server := models.Server{Serverip: serverip}
	err := o.Read(&server, "Serverip")
	if err != nil {
		beego.Info("查询错误", err)
		return
	}
	o.Delete(&server)
	//返回列表页面

	this.Layout = "admin/layout.tpl"
	this.Redirect("/server/srvinfo.html", 302)

}


//error page
type Error1Controller struct {
	beego.Controller
}
func (this *Error1Controller) Get() {
	this.Layout = "admin/layout.tpl"
	this.Redirect("/server/error1.html", 302)
}

type Error2Controller struct {
	beego.Controller
}
func (this *Error2Controller) Get() {
	this.Layout = "admin/layout.tpl"
	this.Redirect("/server/srvexist.html", 302)
}
