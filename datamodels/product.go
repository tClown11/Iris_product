package datamodels

type Product struct {
	ID           int64  `json:"id" sql:"ID" imooc:"ID"`
	ProductName  string `json:"product_name" sql:"productname" imooc:"ProductName"`
	ProductNum   int64  `json:"product_num" sql:"productnum" imooc:"ProductNum"`
	ProductImage string `json:"product_image" sql:"productimage" imooc:"ProductImage"`
	ProductUrl   string `json:"product_url" sql:"producturl" imooc:"ProductUrl"`
}
