package services

import (
	"Iris_product/datamodels"
	"Iris_product/repositories"
)

type IOrderService interface {
	GetOrderByID(int64) (*datamodels.Order, error)
	DeleteOrderByID(int64) bool
	UpdateOrder(*datamodels.Order) error
	InsertOrder(*datamodels.Order) (int64, error)
	GetAllOrder() ([]*datamodels.Order, error)
	GetAllOrderInfo() (map[int]map[string]string, error)
	InsertOrderBymessage(message *datamodels.Message) (productId int64, err error)
}

func NewOrderService(repositories repositories.IOrderRepository) IOrderService {
	return &OrderService{OrderRepository: repositories}
}

type OrderService struct {
	OrderRepository repositories.IOrderRepository
}

func (o *OrderService) GetOrderByID(orderID int64) (order *datamodels.Order, err error) {
	return o.OrderRepository.SelectByKey(orderID)
}

func (o *OrderService) DeleteOrderByID(orderID int64) bool {
	return o.OrderRepository.Delete(orderID)
}

func (o *OrderService) UpdateOrder(order *datamodels.Order) (err error) {
	return o.OrderRepository.Update(order)
}

func (o *OrderService) InsertOrder(order *datamodels.Order) (orderID int64, err error) {
	return o.OrderRepository.Insert(order)
}

func (o *OrderService) GetAllOrder() (orderArray []*datamodels.Order, err error) {
	return o.OrderRepository.SelectAll()
}

func (o *OrderService) GetAllOrderInfo() (map[int]map[string]string, error) {
	return o.OrderRepository.SelectWithInfo()
}

//根据消息创建订单
func (o *OrderService) InsertOrderBymessage(messgae *datamodels.Message) (orderID int64, err error) {
	order := &datamodels.Order{
		UserId:      messgae.UserID,
		ProductId:   messgae.ProductID,
		OrderStatus: datamodels.OrderSuccess,
	}
	return o.InsertOrder(order)
}
