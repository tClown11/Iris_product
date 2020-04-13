package repositories

import (
	"Iris_product/common"
	"Iris_product/datamodels"
	"database/sql"
	"fmt"
	"strconv"
)

//第一步，先开发对应的接口
//第二步实现接口

type IProduct interface {
	//链接数据库
	Conn() error
	Insert(*datamodels.Product) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Product) error
	SelectByKey(int64) (*datamodels.Product, error)
	SelectAll() ([]*datamodels.Product, error)
	SubProductNum(productID int64) error
}

type ProductManager struct {
	table     string
	mysqlconn *sql.DB
}

func NewProductManager(table string, db *sql.DB) IProduct {
	return &ProductManager{table: table, mysqlconn: db}
}

func (p *ProductManager) Conn() (err error) {
	if p.mysqlconn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		p.mysqlconn = mysql
	}
	if p.table == "" {
		p.table = "product"
	}
	return
}

func (p *ProductManager) Insert(product *datamodels.Product) (productId int64, err error) {
	//1.判断链接是否存在
	if err = p.Conn(); err != nil {
		return
	}

	//2.准备sql
	sql := "INSERT product SET productname=?,productnum=?,productimage=?,producturl=?"
	stmt, err := p.mysqlconn.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		return 0, err
	}

	//3.传入参数
	result, err := stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (p *ProductManager) Delete(productID int64) bool {
	if err := p.Conn(); err != nil {
		return false
	}

	sql := "DELETE from product where ID=?"

	stmt, err := p.mysqlconn.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		return false
	}
	_, err = stmt.Exec(strconv.FormatInt(productID, 10))
	if err != nil {
		return false
	}
	return true
}

func (p *ProductManager) Update(product *datamodels.Product) error {
	if err := p.Conn(); err != nil {
		return err
	}

	sql := "Update product set productname=?,productnum=?,productimage=?,producturl=? where ID=" + strconv.FormatInt(product.ID, 10)
	stemt, err := p.mysqlconn.Prepare(sql)
	defer stemt.Close()
	if err != nil {
		return err
	}

	_, err = stemt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if err != nil {
		return err
	}

	return nil
}

//根据商品ID查询商品
func (p *ProductManager) SelectByKey(productID int64) (productResult *datamodels.Product, err error) {
	if err = p.Conn(); err != nil {
		return &datamodels.Product{}, err
	}

	sql := "SELECT * from " + p.table + " where ID=" + strconv.FormatInt(productID, 10)
	row, err := p.mysqlconn.Query(sql)
	defer row.Close()
	if err != nil {
		return &datamodels.Product{}, err
	}

	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Product{}, err
	}
	productResult = &datamodels.Product{}
	fmt.Println(productResult)

	common.DataToStructByTagSql(result, productResult)
	return
}

func (p *ProductManager) SelectAll() (productArray []*datamodels.Product, errproduct error) {
	if err := p.Conn(); err != nil {
		return nil, err
	}

	sql := "Select * from " + p.table
	rows, err := p.mysqlconn.Query(sql)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, nil
	}

	for _, v := range result {
		product := &datamodels.Product{}
		common.DataToStructByTagSql(v, product)
		productArray = append(productArray, product)
	}
	return
}

func (p *ProductManager) SubProductNum(productID int64) error {
	if err := p.Conn(); err != nil {
		return nil
	}

	sql := "UPDATE " + p.table + " set " + " productNum=productNum-1 where ID=" + strconv.FormatInt(productID, 10)
	stmt, err := p.mysqlconn.Prepare(sql)
	if err != nil {
		fmt.Println(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
