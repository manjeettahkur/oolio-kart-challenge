package repositories

import (
	"context"
	"ooliokartchallenge/internal/domain/entities"
	"ooliokartchallenge/internal/domain/errors"
	"ooliokartchallenge/internal/domain/interfaces"
)

type ProductRepository struct {
	products []entities.Product
}

func NewProductRepository() interfaces.ProductRepository {
	return &ProductRepository{
		products: getSampleProducts(),
	}
}

func (r *ProductRepository) GetAll(ctx context.Context) ([]entities.Product, error) {
	products := make([]entities.Product, len(r.products))
	copy(products, r.products)
	return products, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id string) (*entities.Product, error) {
	for _, product := range r.products {
		if product.ID == id {
			productCopy := product
			return &productCopy, nil
		}
	}

	return nil, errors.ErrProductNotFound
}

func getSampleProducts() []entities.Product {
	return []entities.Product{
		{
			ID:       "10",
			Name:     "iPhone 15 Pro",
			Price:    999.99,
			Category: "Phone",
			Image: entities.Image{
				Thumbnail: "https://orderfoodonline.deno.dev/public/images/image-waffle-thumbnail.jpg",
				Mobile:    "https://orderfoodonline.deno.dev/public/images/image-waffle-mobile.jpg",
				Tablet:    "https://orderfoodonline.deno.dev/public/images/image-waffle-tablet.jpg",
				Desktop:   "https://orderfoodonline.deno.dev/public/images/image-waffle-desktop.jpg",
			},
		},
		{
			ID:       "11",
			Name:     "Samsung Galaxy S24",
			Price:    849.99,
			Category: "Phone",
			Image: entities.Image{
				Thumbnail: "https://orderfoodonline.deno.dev/public/images/image-waffle-thumbnail.jpg",
				Mobile:    "https://orderfoodonline.deno.dev/public/images/image-waffle-mobile.jpg",
				Tablet:    "https://orderfoodonline.deno.dev/public/images/image-waffle-tablet.jpg",
				Desktop:   "https://orderfoodonline.deno.dev/public/images/image-waffle-desktop.jpg",
			},
		},
		{
			ID:       "12",
			Name:     "iPad Pro 12.9",
			Price:    1099.99,
			Category: "Tablet",
			Image: entities.Image{
				Thumbnail: "https://orderfoodonline.deno.dev/public/images/image-waffle-thumbnail.jpg",
				Mobile:    "https://orderfoodonline.deno.dev/public/images/image-waffle-mobile.jpg",
				Tablet:    "https://orderfoodonline.deno.dev/public/images/image-waffle-tablet.jpg",
				Desktop:   "https://orderfoodonline.deno.dev/public/images/image-waffle-desktop.jpg",
			},
		},
		{
			ID:       "13",
			Name:     "MacBook Pro 14",
			Price:    1999.99,
			Category: "Laptop",
			Image: entities.Image{
				Thumbnail: "https://orderfoodonline.deno.dev/public/images/image-waffle-thumbnail.jpg",
				Mobile:    "https://orderfoodonline.deno.dev/public/images/image-waffle-mobile.jpg",
				Tablet:    "https://orderfoodonline.deno.dev/public/images/image-waffle-tablet.jpg",
				Desktop:   "https://orderfoodonline.deno.dev/public/images/image-waffle-desktop.jpg",
			},
		},
		{
			ID:       "14",
			Name:     "Dell XPS 13",
			Price:    1299.99,
			Category: "Laptop",
			Image: entities.Image{
				Thumbnail: "https://orderfoodonline.deno.dev/public/images/image-waffle-thumbnail.jpg",
				Mobile:    "https://orderfoodonline.deno.dev/public/images/image-waffle-mobile.jpg",
				Tablet:    "https://orderfoodonline.deno.dev/public/images/image-waffle-tablet.jpg",
				Desktop:   "https://orderfoodonline.deno.dev/public/images/image-waffle-desktop.jpg",
			},
		},
	}
}
