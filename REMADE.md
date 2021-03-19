### 使用
`go get github.com/wuyan94zl/mysql`

```go

package main

import (
	"fmt"
	"github.com/wuyan94zl/mysql"
)
type blog struct {
	Id      uint64 `json:"id"`
	Title   string `json:"title"validate:"required||min:12||max:32"fieldName:"博客标题"`
	Content string `json:"content"validate:"required"fieldName:"博客内容"`
	View    uint64 `json:"view"validate:"numeric"fieldName:"浏览数"`
}

func init() {
	c := mysql.Config{
		// 必填配置
		Username: "root",
		Password: "123456",
		Database: "blog",

		// 可选默认配置
		//Host:           "127.0.0.1",
		//Port:           3306,
		//Charset:        "utf8mb4",

		// 可选配置（连接池）
		//MaxConnect:     100,
		//MaxIdleConnect: 25,
		//MaxLifeSeconds: 300,
	}
	mysql.ConMysql(c)

	// 表结构迁移
	MigrateStruct := make(map[string]interface{})
	MigrateStruct["blog"] = &blog{}
	mysql.AutoMigrate(MigrateStruct)
}
func main() {
	c := blog{Title: "测试", Content: "测试数据", View: 66}
	_ = mysql.GetInstance().Create(&c)
	var u blog
	var w = make(map[string]interface{})
	w["id"] = mysql.Where{Way: "=", Value: c.Id}
	_ = mysql.GetInstance().Where(w).One(&u)
	fmt.Println("查询添加数据：", u)

	w["id"] = mysql.Where{Way: "=", Value: c.Id + 1}
	err := mysql.GetInstance().Where(w).One(&u)
	if err != nil {
		fmt.Println("查询数据错误：", err)
	}
}

```