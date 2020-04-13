package repositories

import (
	"Iris_product/common"
	"Iris_product/datamodels"
	"database/sql"
	"strconv"
)

type IOrderRepository interface {
	Conn() error
	Insert(*datamodels.Order) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Order) error
	SelectByKey(int64) (*datamodels.Order, error)
	SelectAll() ([]*datamodels.Order, error)
	SelectWithInfo() (map[int]map[string]string, error)
}

func NewOrderManageRepository(table string, db *sql.DB) IOrderRepository{
	return &OrderManageRepository{table:table, mysqlConn:db}
}

type OrderManageRepository struct {
	table string
	mysqlConn *sql.DB
}

func (o *OrderManageRepository) Conn() error {
	if o.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		o.mysqlConn = mysql
	}

	if o.table == "" {
		o.table = "order"
	}
	return nil
}

func (o *OrderManageRepository) Insert(order *datamodels.Order) (productID int64, err error) {
	if err = o.Conn(); err != nil {
		return
	}

	sql := "INSERT `order` SET userID=?,productID=?,orderStatus=?"
	stmt, err := o.mysqlConn.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		return productID, err
	}
	result, err := stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	if err != nil {
		return productID, err
	}

	return result.LastInsertId()
}

func (o *OrderManageRepository) Delete(productID int64) (isOK bool) {
	if err := o.Conn(); err != nil {
		return
	}

	sql := "DELETE from " + o.table + " where ID=?"
	stemt, err := o.mysqlConn.Prepare(sql)
	defer stemt.Close()
	if err != nil {
		return
	}

	_, err = stemt.Exec(productID)
	if err != nil {
		return
	}
	return true
}

func (o *OrderManageRepository) Update(order *datamodels.Order) (err error) {
	if err := o.Conn(); err != nil {
		return err
	}

	sql := "UPDATE " + o.table + " set userID=?,productID=?,orderStatus=? where ID=" + strconv.FormatInt(order.ID, 10)
	stmt, err := o.mysqlConn.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		return err
	}

	_, err = stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	return
}

func (o *OrderManageRepository) SelectByKey(orderID int64) (order *datamodels.Order, err error) {
	if err := o.Conn(); err != nil {
		return &datamodels.Order{}, err
	}
	sql := "SELECT * From " + o.table + " where ID=" + strconv.FormatInt(orderID, 10)
	row, err := o.mysqlConn.Query(sql)
	defer row.Close()
	if err != nil {
		return &datamodels.Order{}, err
	}

	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Order{}, err
	}

	order = &datamodels.Order{}
	common.DataToStructByTagSql(result, order)
	return
}

func (o *OrderManageRepository) SelectAll() (orderArray []*datamodels.Order, err error) {
	if err := o.Conn(); err != nil {
		return nil, err
	}

	sql := "SELECT * From " + o.table
	rows, err := o.mysqlConn.Query(sql)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, err
	}

	for _, v := range result {
		order := &datamodels.Order{}
		common.DataToStructByTagSql(v, order)
		orderArray = append(orderArray, order)
	}
	return
}

func (o *OrderManageRepository) SelectWithInfo() (OrderMap map[int]map[string]string, err error) {
	if err := o.Conn(); err != nil {
		return nil, err
	}

	sql := "Select o.ID,p.productName,o.orderStatus From text.order as o left join product as p on o.productID=p.ID"
	rows, err := o.mysqlConn.Query(sql)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	return common.GetResultRows(rows), err
}
