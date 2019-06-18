package controllers

import (
	"ZDHsys/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.TplName = "register.html"
}
func (c *MainController) Post() {
	//1.拿到数据
	userName:=c.GetString("userName")
	pwd:=c.GetString("pwd")
	//2.对数据进行校验
	if userName =="" || pwd==""{
		beego.Info("数据不能为空")
		c.Redirect("/register",302)
		return
	}
	//3.插入数据库
	o:=orm.NewOrm()

	user:=models.User{}
	user.Name=userName
	user.Pwd=pwd
	_,err:=o.Insert(&user)
	if err!=nil{
		beego.Info("插入数据库失败")
		c.Redirect("/register",302)
		return
	}
	//4.返回登陆页面
	c.TplName="login.html"

}
func (c*MainController) ShowLogin(){
	c.TplName="login.html"

}
func (c*MainController) HandleLogin(){
	//1.拿到数据
	userName:=c.GetString("userName")
	pwd:=c.GetString("pwd")
	//2.判断数据是否合法
	if userName=="" ||pwd==""{
		beego.Info("输入数据不合法")
		c.TplName="login.html"
		return
	}
	//3.查询帐号密码是否正确
	o:=orm.NewOrm()
	user:=models.User{}

	user.Name=userName
	err:=o.Read(&user,"Name")
	if err!=nil{
		beego.Info("查询失败")
		c.TplName="login.html"
		return
	}
	//4.跳转
	c.Data["title"] = beego.AppConfig.String("title")
	c.Data["culture"] = beego.AppConfig.String("culture")
	c.Layout = "admin/layout.tpl"
	c.TplName = "admin/index.tpl"

}
//func (c*MainController) ShowIndex(){
//	c.TplName="index.html"
//}
func (c *MainController) ShowIndex() {
	c.Data["title"] = beego.AppConfig.String("title")
	c.Data["culture"] = beego.AppConfig.String("culture")
	c.Data["navtitle"] = "欢迎页"
	c.Data["datetime"] = time.Now().Format("2006-01-02 15:04:05")
	c.Layout = "admin/layout.tpl"
	c.TplName = "admin/index.tpl"
}