package repository

import (
	"codebase-app/internal/module/product/entity"
	"codebase-app/internal/module/product/ports"
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.ProductRepository = &productRepository{}

type productRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *productRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) CreateProduct(ctx context.Context, req *entity.CreateProductRequest) (*entity.CreateProductResponse, error) {
	var resp = new(entity.CreateProductResponse)
	// Your code here
	query := `
		INSERT INTO products (user_id, name, description, shop_id, category_id, price, stock)
		VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id
	`

	err := r.db.QueryRowContext(ctx, r.db.Rebind(query),
		req.UserId,
		req.Name,
		req.Description,
		req.ShopId,
		req.CategoryId,
		req.Price,
		req.Stock).Scan(&resp.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::CreateProduct - Failed to create product")
		return nil, err
	}

	return resp, nil
}

func (r *productRepository) GetProduct(ctx context.Context, req *entity.GetProductRequest) (*entity.GetProductResponse, error) {
	var resp = new(entity.GetProductResponse)
	// Your code here
	query := `
		SELECT a.name, a.description, a.shop_id, b.name shop_name, a.category_id, c.name category_name, a.price, a.stock
		FROM products as a
		INNER JOIN shops b ON b.id = a.shop_id
		INNER JOIN categories c ON c.id = a.category_id
		WHERE a.id = ? AND a.deleted_at is NULL
	`

	err := r.db.QueryRowxContext(ctx, r.db.Rebind(query), req.Id).StructScan(resp)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::GetProduct - Failed to get product")
		return nil, err
	}

	return resp, nil
}

func (r *productRepository) DeleteProduct(ctx context.Context, req *entity.DeleteProductRequest) error {
	query := `
		UPDATE products
		SET deleted_at = NOW()
		WHERE id = ? AND user_id = ?
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), req.Id, req.UserId)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::DeleteProduct - Failed to delete product")
		return err
	}

	return nil
}

func (r *productRepository) UpdateProduct(ctx context.Context, req *entity.UpdateProductRequest) (*entity.UpdateProductResponse, error) {
	var resp = new(entity.UpdateProductResponse)

	query := `
		UPDATE products
		SET name = ?, description = ?, price = ?, updated_at = NOW()
		WHERE id = ? AND user_id = ?
		RETURNING id
	`

	err := r.db.QueryRowxContext(ctx, r.db.Rebind(query),
		req.Name,
		req.Description,
		req.Price,
		req.Id,
		req.UserId).Scan(&resp.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::UpdateProduct - Failed to update product")
		return nil, err
	}

	return resp, nil
}

func (r *productRepository) GetProducts(ctx context.Context, req *entity.ProductRequest) (*entity.ProductResponse, error) {
	type dao struct {
		TotalData int `db:"total_data"`
		entity.ProductItem
	}

	var (
		resp = new(entity.ProductResponse)
		data = make([]dao, 0, req.Paginate)
	)
	resp.Items = make([]entity.ProductItem, 0, req.Paginate)

	query := `
		SELECT
			COUNT(a.id) OVER() as total_data,
			a.id,
			a.name,
			a.description,
			a.shop_id,
			b.name shop_name,
			a.category_id,
			c.name category_name,
			a.price,
			a.stock
		FROM products a
		INNER JOIN shops b ON b.id = a.shop_id
		INNER JOIN categories c ON c.id = a.category_id
		WHERE a.user_id = ? AND a.shop_id = ? AND a.deleted_at is NULL
		LIMIT ? OFFSET ?
	`

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query),
		req.UserId,
		req.ShopId,
		req.Paginate,
		req.Paginate*(req.Page-1),
	)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::GetProducts - Failed to get products")
		return nil, err
	}

	if len(data) > 0 {
		resp.Meta.TotalData = data[0].TotalData
	}

	for _, d := range data {
		resp.Items = append(resp.Items, d.ProductItem)
	}

	resp.Meta.CountTotalPage(req.Page, req.Paginate, resp.Meta.TotalData)

	return resp, nil
}

// func (r *shopRepository) IsShopOwner(ctx context.Context, userId, shopId string) (bool, error) {
// 	panic("implement me")
// }
