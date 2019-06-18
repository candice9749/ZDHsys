package server

import (
	"fmt"
	"github.com/sbinet/go-python"
	"testing"
)
//为了单独测试server.go这个文件的，不然只有main才能执行
func TestCobbler(t *testing.T) {
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
	interfaceString := "{'macaddress-eth0':'66:66:62:66:66:66','ipaddress-eth0':'192.168.142.106','static-eth0':1}"
	res := g.CallFunction("wangyutongj3","la","centos7.5-x86_64","/var/lib/cobbler/kickstarts/centos7_5.ks",interfaceString)
	if res == nil {
		fmt.Println("callfunction error")
		return
	}
	//item := python.PyList_GetItem(res, 0)
	//d := python.PyDict_GetItemString(res, "X")
	//dd := python.PyString_AsString(d)
	//fmt.Println(dd)
}
