package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)
//用户管理结构体
type User struct {
	Id int
	Name string
	Pwd string
}
type Role struct {
	Id int
	Username string
	Name string
	Role string
	Line string
	Section string
}


//机器管理结构体
type Server struct {
	Id int
	Serverip string `null`
	Hostname string `orm:"default(localhost)"`
	Mac string `null`
	Os  string `null`
	Stime time.Time `orm:"auto_now"`
	Addr string `null`
	Line string `null`
	Ywman string `orm:"size(20)"`
	Cpu string `orm:"default(英特尔 酷睿 i5-5200U)"`
	Memory string`orm:"default(4GB)"`
	Mainboard string `null`
	Graphics string `orm:"default(AMD Radeon R5 M330)"`
}



func init(){
	//设置数据库基本信息
	orm.RegisterDataBase("default","mysql","root:WWW.1.com@tcp(192.168.122.222:3306)/test?charset=utf8")
	//映射model数据,多个表用，分开
	orm.RegisterModel(new(User),new(Server),new(Role))
	//生成表,default就是我们要用到的数据库别名，force是是否自动更新，当表的字段发生改变的时候，verbose是显示创建表的过程
	orm.RunSyncdb("default",false,true)
}