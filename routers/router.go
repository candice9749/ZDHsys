package routers

import (
	"ZDHsys/controllers"
	"ZDHsys/controllers/server"
	"ZDHsys/controllers/user"
	"github.com/astaxie/beego"
	"ZDHsys/controllers/salt"
	"ZDHsys/controllers/ssh"

)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.Router("/register", &controllers.MainController{})
	beego.Router("/login", &controllers.MainController{},"get:ShowLogin;post:HandleLogin")
	//beego.Router("/index", &controllers.MainController{},"get:ShowIndex")
	beego.Router("/index", &controllers.MainController{},"get:ShowIndex")

	beego.Router("/server/srvinfo", &server.InfoController{})
	beego.Router("/server/srvadd", &server.AddController{})
	beego.Router("/server/hardinfo", &server.HardController{})
	beego.Router("/server/srvupdate", &server.UpdateController{})
	beego.Router("/server/delete", &server.DeleteController{})
	beego.Router("/server/srvexist", &server.ExistController{})
	beego.Router("/server/srvinput", &server.InputController{})
	beego.Router("/server/error1", &server.Error1Controller{})
	beego.Router("/server/error2", &server.Error2Controller{})

	beego.Router("/user/userinfo", &user.InfoController{})
	beego.Router("/user/roleadd", &user.RoleController{})
	beego.Router("/user/userupdate", &user.UpdateController{})
	beego.Router("/user/delete", &user.DeleteController{})

	beego.Router("/ssh/index", &ssh.SshController{})

	beego.Router("/salt/ping", &salt.PingController{})
	beego.Router("/salt/keylist", &salt.KeyListController{})
	beego.Router("/salt/keydelete", &salt.KeyDeleteController{})
	beego.Router("/salt/keyaccept", &salt.KeyAcceptController{})
	beego.Router("/salt/cmdrun", &salt.CmdRunController{})
	beego.Router("/salt/cpgetfile", &salt.CpGetFileController{})
	beego.Router("/salt/deploy", &salt.DeployController{})
}
