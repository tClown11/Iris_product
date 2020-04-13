package main

import (
	"Iris_product/backend/web/controllers"
	"Iris_product/common"
	"Iris_product/repositories"
	"Iris_product/services"
	"context"

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
	app.RegisterView(tmplate)

	//4.设置模板目标
	app.StaticWeb("/assets", "./web/assets")
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
	//5.注册控制器
	//productRepository := repositories.NewProductManager("product", db)
	//productService := services.NewProductService(productRepository)
	//productParty := app.Party("/product")
	//product := mvc.New(productParty)
	//product.Register(ctx,productService)
	//product.Handle(new(controllers.ProductController))
	controllers.RegisterProduct(ctx, "product", db, app)

	orderRepository := repositories.NewOrderManageRepository("order", db)
	orderService := services.NewOrderService(orderRepository)
	orderParty := app.Party("/order")
	order := mvc.New(orderParty)
	order.Register(ctx, orderService)
	order.Handle(new(controllers.OrderController))

	//6.启动控制器
	app.Run(
		iris.Addr("0.0.0.0:8000"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
