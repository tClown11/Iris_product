package services

import (
	"Iris_product/datamodels"
	"Iris_product/repositories"
)

type IProductService interface {
	GetProductByID(int64) (*datamodels.Product, error)
	GetAllProduct() ([]*datamodels.Product, error)
	DeleteProductByID(int64) bool
	InsertProduct(*datamodels.Product) (int64, error)
	UpdateProduct(*datamodels.Product) error
	SubNumberOne(productID int64) error
}

type ProductService struct {
	productrepository repositories.IProduct
}

func NewProductService(repository repositories.IProduct) IProductService {
	return &ProductService{repository}
}

func (p *ProductService) GetProductByID(productID int64) (*datamodels.Product, error) {
	return p.productrepository.SelectByKey(productID)
}

func (p *ProductService) GetAllProduct() ([]*datamodels.Product, error) {
	return p.productrepository.SelectAll()
}

func (p *ProductService) DeleteProductByID(productID int64) bool {
	return p.productrepository.Delete(productID)
}

func (p *ProductService) InsertProduct(product *datamodels.Product) (int64, error) {
	return p.productrepository.Insert(product)
}

func (p *ProductService) UpdateProduct(product *datamodels.Product) error {
	return p.productrepository.Update(product)
}

func (p *ProductService) SubNumberOne(productID int64) error {
	return p.productrepository.SubProductNum(productID)
}
