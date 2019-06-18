package user

import (
	"ZDHsys/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type InfoController struct {
	beego.Controller
}
func (this *InfoController) Get() {

	//有orm对象
	o := orm.NewOrm()
	//查询对象
	var roles []models.Role           //定义一个结构体数组
	_,err := o.QueryTable("Role").All(&roles)        //高级查询
	if err != nil {
		beego.Info("查询所有信息失败")
		return
	}
	//传给视图
	this.Data["roles"]=roles
	this.Data["navtitle"] = "用户信息"
	this.Data["title"] = "ZDHYW平台"
	this.Layout = "admin/layout.tpl"
	this.TplName = "user/userinfo.html"
}



type RoleController struct {
	beego.Controller
}
func (this *RoleController) Get() {
	this.Data["title"] = "ZDHYW平台"
	this.Data["navtitle"] = "添加用户"
	this.Layout = "admin/layout.tpl"
	this.TplName = "user/roleadd.html"
}

func (this *RoleController) Post() {
	//1.拿到数据
	username := this.GetString("username")
	name := this.GetString("name")
	role := this.GetString("role")
	line := this.GetString("line")
	section := this.GetString("section")


	//2.判断数据是否合法
	if username == "" || name == "" || role==""||line=="" || section==""{
		beego.Info("添加用户数据错误")
		return
	}

	//3.插入数据
	o := orm.NewOrm()         //申请orm对象
	roles := models.Role{} //有一个要插入的结构体
	roles.Username=username
	roles.Name=name
	roles.Role=role
	roles.Line = line
	roles.Section=section


	_, err := o.Insert(&roles)
	if err != nil {
		beego.Info("插入数据库错误", err)
		return
	}

	//4.返回主机界面

	this.Layout = "admin/layout.tpl"
	this.Redirect("/user/userinfo.html",302)



}


//编辑页,更新
type UpdateController struct {
	beego.Controller
}
func (this *UpdateController) Get() {
	//1.获取主机ip
	id ,_:=this.GetInt("id") //注意要和前端页面里面的name是相同的
	//2.查询数据库获取信息
	o := orm.NewOrm()
	roles := models.Role{Id:id}
	err := o.Read(&roles)
	if err!= nil {
		beego.Info("查询错误",err)
		return
	}
	//3.传递数据给视图
	this.Data["roles"]=roles
	this.Data["title"] = "ZDHYW平台"
	this.Layout = "admin/layout.tpl"
	this.TplName = "user/userupdate.html"
}

//处理更新页面的数据
func (this *UpdateController) Post() {
	//1.拿到数据
	id ,_:=this.GetInt("id")
	username := this.GetString("username")
	name := this.GetString("name")
	role := this.GetString("role")
	line := this.GetString("line")
	section := this.GetString("section")

	//2.对数据进行一个处理
	if username == "" || name == "" || role==""||line=="" || section==""{
		beego.Info("添加用户数据错误")
		return
	}

	//3.更新操作
	o:=orm.NewOrm()
	roles := models.Role{Id:id}
	err:=o.Read(&roles)
	if err !=nil{
		beego.Info("查询数据错误！")
		return
	}
	roles.Username=username
	roles.Name=name
	roles.Role=role
	roles.Line=line
	roles.Section=section


	_,err=o.Update(&roles)
	if err !=nil{
		beego.Info("更新数据显示错误")
		return
	}
	//4.返回列表页面
	this.Layout = "admin/layout.tpl"
	this.Redirect("/user/userinfo.html",302)

}

//删除操作
type DeleteController struct {
	beego.Controller
}
func (this *DeleteController) Get() {
	//1.拿到数据
	id ,_:=this.GetInt("id")
	//2.执行删除操作
	o := orm.NewOrm()
	roles := models.Role{Id:id}
	err := o.Read(&roles)
	if err!= nil {
		beego.Info("查询错误",err)
		return
	}
	o.Delete(&roles)
	//返回列表页面
	this.Layout = "admin/layout.tpl"
	this.Redirect("/user/userinfo.html",302)


}