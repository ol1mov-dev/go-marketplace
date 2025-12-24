package handlers

import (
	"context"
	"database/sql"
	"log"
	"time"

	productV1 "github.com/ol1mov-dev/protos/pkg/product/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Product struct {
	Id               int32
	Name             string
	Sku              string
	ShortDescription string
	Description      string
	Price            float32
	PriceOld         float32
	Discount         float32
	Quantity         int32
	IsActive         bool
	Brand            string
	Rating           float32
	CategoryId       int32
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type ProductServer struct {
	productV1.UnimplementedProductV1ServiceServer
	DB *sql.DB
}

func (s *ProductServer) CreateProduct(ctx context.Context, req *productV1.CreateProductRequest) (*productV1.CreateProductResponse, error) {
	var productId uint32

	query := `INSERT INTO products (name, sku, short_description, description, price, price_old, discount, quantity, is_active, brand, rating, category_id) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`

	log.Println("Creating product...")
	err := s.DB.QueryRow(
		query,
		req.Name,
		req.Sku,
		req.ShortDescription,
		req.Description,
		req.Price,
		req.PriceOld,
		req.Discount,
		req.Quantity,
		req.IsActive,
		req.Brand,
		req.Rating,
		req.CategoryId,
	).Scan(&productId)

	if err != nil {
		log.Println("Error creating product")
		return nil, err
	}

	log.Println("Product created!")
	return &productV1.CreateProductResponse{ProductId: productId}, nil
}

func (s *ProductServer) UpdateProduct(ctx context.Context, req *productV1.UpdateProductRequest) (*productV1.UpdateProductResponse, error) {

	query := `
		UPDATE products
		SET
			name = $1,
			sku = $2,
			short_description = $3,
			description = $4,
			price = $5,
			price_old = $6,
			discount = $7,
			quantity = $8,
			is_active = $9,
			brand = $10,
			rating = $11,
			category_id = $12,
			updated_at = NOW()
		WHERE id = $13
	`

	// 3. Выполнение
	res, err := s.DB.ExecContext(
		ctx,
		query,
		req.Name,
		req.Sku,
		req.ShortDescription,
		req.Description,
		req.Price,
		req.PriceOld,
		req.Discount,
		req.Quantity,
		req.IsActive,
		req.Brand,
		req.Rating,
		req.CategoryId,
		req.ProductId,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db error: %v", err)
	}

	// 4. Проверка, что продукт найден
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "rows affected error: %v", err)
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, "product not found")
	}

	// 5. Ответ
	return &productV1.UpdateProductResponse{
		ProductId: req.ProductId,
	}, nil
}

func (s *ProductServer) DeleteProduct(ctx context.Context, req *productV1.DeleteProductRequest) (*productV1.DeleteProductResponse, error) {
	var productId uint32

	log.Println("Deleting product...")
	query := `DELETE FROM products WHERE id = $1 RETURNING id`
	err := s.DB.QueryRow(query, req.ProductId).Scan(&productId)

	if err != nil {
		log.Println("Error deleting product")
		return nil, err
	}

	log.Println("Product deleted!")
	return &productV1.DeleteProductResponse{ProductId: productId}, nil
}

func (s *ProductServer) GetAllProductsByQuery(ctx context.Context, req *productV1.GetAllProductsByQueryRequest) (*productV1.GetAllProductsByQueryResponse, error) {
	log.Println("Getting all products by query...")
	query := `SELECT id, name, sku, short_description, description, price, price_old, discount, quantity, is_active, brand, rating, category_id
    			FROM products WHERE (is_active = TRUE) AND name LIKE CONCAT('%', $1::text, '%')`

	rows, err := s.DB.Query(query, req.Query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*productV1.Product

	for rows.Next() {
		product := &productV1.Product{}

		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Sku,
			&product.ShortDescription,
			&product.Description,
			&product.Price,
			&product.PriceOld,
			&product.Discount,
			&product.Quantity,
			&product.IsActive,
			&product.Brand,
			&product.Rating,
			&product.CategoryId,
		)
		if err != nil {
			return nil, err
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	log.Println("Getting products by query done!")
	return &productV1.GetAllProductsByQueryResponse{
		Products: products,
	}, nil
}
