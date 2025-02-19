package main

import (
	"context"
	"fmt"
	"log"

	"main/internal/application/order"
	"main/internal/application/product"
	"main/internal/infrastructure/persistence/mongodb"
	"main/utility/mongodb"
)

func main() {
	ctx := context.Background()

	// 1. 初始化 MongoDB 配置
	mongoConfig := mongodb.Config{
		URI:      "mongodb://localhost:27017",
		Database: "ddd_example",
	}

	// 2. 创建仓储实例
	productRepo, err := mongodb.NewProductRepository(ctx, mongoConfig)
	if err != nil {
		log.Fatalf("Failed to create product repository: %v", err)
	}

	orderRepo, err := mongodb.NewOrderRepository(ctx, mongoConfig)
	if err != nil {
		log.Fatalf("Failed to create order repository: %v", err)
	}

	// 3. 创建应用服务
	productService := product.NewProductService(productRepo)
	orderService := order.NewOrderService(orderRepo, productRepo)

	// 4. 创建商品
	product1, err := productService.CreateProduct(ctx, product.CreateProductCommand{
		Name:        "iPhone 15",
		Description: "Latest iPhone model",
		Price:       7999.00,
		Stock:       100,
	})
	if err != nil {
		log.Fatalf("Failed to create product: %v", err)
	}
	fmt.Printf("Created product: %s\n", product1.Name)

	product2, err := productService.CreateProduct(ctx, product.CreateProductCommand{
		Name:        "AirPods Pro",
		Description: "Wireless earbuds",
		Price:       1999.00,
		Stock:       50,
	})
	if err != nil {
		log.Fatalf("Failed to create product: %v", err)
	}
	fmt.Printf("Created product: %s\n", product2.Name)

	// 5. 创建订单
	newOrder, err := orderService.CreateOrder(ctx, order.CreateOrderCommand{
		UserId: "user123",
		Items: []order.OrderItemCommand{
			{
				ProductId: product1.Id,
				Quantity: 1,
			},
			{
				ProductId: product2.Id,
				Quantity: 2,
			},
		},
		Remark: "Gift package needed",
	})
	if err != nil {
		log.Fatalf("Failed to create order: %v", err)
	}
	fmt.Printf("Created order: %s\n", newOrder.Id)

	// 6. 支付订单
	err = orderService.PayOrder(ctx, order.PayOrderCommand{
		OrderId:        newOrder.Id,
		Amount:         newOrder.TotalAmount.Amount(),
		PaymentMethod: "Alipay",
		PaymentChannel: "APP",
		TradeNo:        "2024021912345678",
	})
	if err != nil {
		log.Fatalf("Failed to pay order: %v", err)
	}
	fmt.Printf("Order paid: %s\n", newOrder.Id)

	// 7. 查询商品列表
	products, err := productService.ListProducts(ctx, product.ListProductsQuery{})
	if err != nil {
		log.Fatalf("Failed to list products: %v", err)
	}
	fmt.Printf("Found %d products\n", len(products))
	for _, p := range products {
		fmt.Printf("- %s: ¥%.2f (Stock: %d)\n", p.Name, p.Price.Amount(), p.Stock)
	}
}
