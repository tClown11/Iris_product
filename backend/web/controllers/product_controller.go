package controllers

import (
	"Iris_product/common"
	"database/sql"
	"Iris_product/datamodels"
	"Iris_product/services"
	"Iris_product/repositories"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"strconv"
	"context"
)

type ProductController struct {
	Ctx iris.Context
	ProductService services.IProductService
}

func (p *ProductController) GetAll() mvc.View {
	productArray, _ := p.ProductService.GetAllProduct()
	return mvc.View{
		Name: "product/view.html",
		Data: iris.Map{
			"productArray": productArray,
		},
	}
}

func (p *ProductController) PostUpdate (){
	product := &datamodels.Product{}
	p.Ctx.Request().ParseForm()
	dec := common.NewDecoder(&common.DecoderOptions{TagName:"imooc"})
	if err := dec.Decode(p.Ctx.Request().Form, product); err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	err := p.ProductService.UpdateProduct(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	p.Ctx.Redirect("/product/all")
}

func (p *ProductController) GetAdd() mvc.View {
	return mvc.View{
		Name:"product/add.html",
	}
}


func (p *ProductController) PostAdd(){
	product := &datamodels.Product{}
	//解析from表单, parseFrom将提交的from表单添加到p.From里面
	p.Ctx.Request().ParseForm()
	//选取datamodels里面的tag标签作为上面生成的product的key
	dec := common.NewDecoder(&common.DecoderOptions{TagName:"imooc"})
	//将from的value值与对应的key放到product里面
	if err := dec.Decode(p.Ctx.Request().Form, product);err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	//调用Service的Insert函数进行数据插入
	_, err := p.ProductService.InsertProduct(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	p.Ctx.Redirect("/product/all")
}

func (p *ProductController) GetManager() mvc.View{
	idstring := p.Ctx.URLParam("id")
	id, err := strconv.ParseInt(idstring, 10, 16)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	product, err := p.ProductService.GetProductByID(id)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	return mvc.View{
		Name: "product/manager.html",
		Data: iris.Map{
			"product": product,
		},
	}
}


func (p *ProductController) GetDelete() {
	idString := p.Ctx.URLParam("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	isOK := p.ProductService.DeleteProductByID(id)
	if isOK {
		p.Ctx.Application().Logger().Debug("删除商品成功ID为："+ idString)
	}else {
		p.Ctx.Application().Logger().Debug("删除商品失败ID为: " + idString)
	}
	p.Ctx.Redirect("/product/all")
	}

func RegisterProduct(ctx context.Context, repositoriesName string, db *sql.DB, app *iris.Application){
	productRepository := repositories.NewProductManager(repositoriesName, db)
	ProductService := services.NewProductService(productRepository)
	productParty := app.Party("/prodcut")
	product := mvc.New(productParty)
	product.Register(ctx,ProductService)
	product.Handle(new(ProductController))
}
