package main

import (
	"Iris_product/common"
	"Iris_product/fronted/middlerware"
	"Iris_product/fronted/web/controllers"
	"Iris_product/rabbitmq"
	"Iris_product/repositories"
	"Iris_product/services"
	"context"
	"fmt"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/opentracing/opentracing-go/log"
)

func main() {
	//1.创建iris实例
	app := iris.New()
	//2.设置错误模式
	app.Logger().SetLevel("debug")
	//3.注册模板
	tmplate := iris.HTML("./web/views", ".html").Layout("shared/layout.html").Reload(true)
	fmt.Println(&tmplate)
	app.RegisterView(tmplate)

	//4.设置模板目标
	app.StaticWeb("/public", "./web/public")
	//设置异常跳转页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！ "))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})

	//链接数据库
	db, err := common.NewMysqlConn()
	if err != nil {
		log.Error(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	user := repositories.NewUserRepository("user", db)
	userService := services.NewUserService(user)
	userPro := mvc.New(app.Party("/user"))
	userPro.Register(userService, ctx)
	userPro.Handle(new(controllers.UserController))

	rabbitmq := rabbitmq.NewRabbitMQSimple("iris_product")

	//注册product控制器
	product := repositories.NewProductManager("product", db)
	productService := services.NewProductService(product)
	order := repositories.NewOrderManageRepository("order", db)
	orderService := services.NewOrderService(order)
	proProduct := app.Party("/product")
	pro := mvc.New(proProduct)
	proProduct.Use(middlerware.AuthConProduct)
	pro.Register(productService, orderService, ctx, rabbitmq)
	pro.Handle(new(controllers.ProductController))

	app.Run(
		iris.Addr("0.0.0.0:8000"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
