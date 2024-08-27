package entity

import "codebase-app/pkg/types"

type CreateProductRequest struct {
	UserId string `validate:"uuid" db:"user_id"`

	Name        string `json:"name" validate:"required" db:"name"`
	Description string `json:"description" validate:"required,max=255" db:"description"`
	ShopId      string `json:"shop_id" validate:"required" db:"shop_id"`
	CategoryId  string `json:"category_id" validate:"required" db:"category_id"`
	Price       int    `json:"price" validate:"required" db:"price"`
	Stock       int    `json:"stock" validate:"required" db:"stock"`
}

type CreateProductResponse struct {
	Id string `json:"id" db:"id"`
}

type GetProductRequest struct {
	Id string `validate:"uuid" db:"id"`
}

type GetProductResponse struct {
	Name         string  `json:"name" db:"name"`
	Description  string  `json:"description" db:"description"`
	ShopId       string  `json:"shop_id" db:"shop_id"`
	ShopName     string  `json:"shop_name" db:"shop_name"`
	CategoryId   string  `json:"category_id" db:"category_id"`
	CategoryName string  `json:"category_name" db:"category_name"`
	Price        float64 `json:"price" db:"price"`
	Stock        int     `json:"stock" db:"stock"`
}

type DeleteProductRequest struct {
	UserId string `prop:"user_id" validate:"uuid" db:"user_id"`

	Id string `validate:"uuid" db:"id"`
}

type UpdateProductRequest struct {
	UserId string `prop:"user_id" validate:"uuid" db:"user_id"`

	Id          string  `params:"id" validate:"uuid" db:"id"`
	Name        string  `json:"name" validate:"required" db:"name"`
	Description string  `json:"description" validate:"required" db:"description"`
	Price       float64 `json:"price" validate:"required" db:"price"`
}

type UpdateProductResponse struct {
	Id string `json:"id" db:"id"`
}

type ProductRequest struct {
	UserId   string `prop:"user_id" validate:"uuid"`
	ShopId   string `query:"shop_id" validate:"uuid"`
	Page     int    `query:"page" validate:"required"`
	Paginate int    `query:"paginate" validate:"required"`
}

func (r *ProductRequest) SetDefault() {
	if r.Page < 1 {
		r.Page = 1
	}

	if r.Paginate < 1 {
		r.Paginate = 10
	}
}

type ProductItem struct {
	Id           string  `json:"id" db:"id"`
	Name         string  `json:"name" db:"name"`
	Description  string  `json:"description" db:"description"`
	ShopId       string  `json:"shop_id" db:"shop_id"`
	ShopName     string  `json:"shop_name" db:"shop_name"`
	CategoryId   string  `json:"category_id" db:"category_id"`
	CategoryName string  `json:"category_name" db:"category_name"`
	Price        float64 `json:"price" db:"price"`
	Stock        int     `json:"stock" db:"stock"`
}

type ProductResponse struct {
	Items []ProductItem `json:"items"`
	Meta  types.Meta    `json:"meta"`
}
