package main

import (
	"fmt"
	"minghuaonline/api"
)

func main() { //主函数
	var an string
ou:
	fmt.Printf("请输入开启密码：")
	fmt.Scanf("%s\n", &an)
	switch an {
	case "666":
		fmt.Printf("密码输入正确,请开始刷课\n")
		api.GetSchool()
		api.Login()
	default:
		fmt.Printf("密码输入错误！\n")
		goto ou
	}
}
