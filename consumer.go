package main

import (
	"Iris_product/common"
	"Iris_product/rabbitmq"
	"Iris_product/repositories"
	"Iris_product/services"
	"fmt"
)

func main() {
	db, err := common.NewMysqlConn()
	if err != nil {
		fmt.Println(err)
	}

	//创建product数据库操作实例
	product := repositories.NewProductManager("product", db)
	//创建product service
	productService := services.NewProductService(product)
	//创建Order数据库实例
	order := repositories.NewOrderManageRepository("order", db)
	//创建Order service
	orderService := services.NewOrderService(order)

	rabbitmqConsumerSimple := rabbitmq.NewRabbitMQSimple("iris_product")
	rabbitmqConsumerSimple.ConsumeSimple(orderService, productService)
}
