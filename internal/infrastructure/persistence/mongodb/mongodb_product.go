package mongodb

import (
	"context"

	"main/internal/domain/product/entity"
	"main/internal/domain/product/repository"
	"main/internal/domain/product/valueobject"
	"main/utility/mongodb"

	"github.com/gogf/gf/v2/errors/gerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ProductPO 商品持久化对象
type ProductPO struct {
	Id          string  `bson:"_id"`
	Name        string  `bson:"name"`
	Description string  `bson:"description"`
	Price       MoneyPO `bson:"price"`
	Stock       int     `bson:"stock"`
	Status      string  `bson:"status"`
	CreatedAt   int64   `bson:"created_at"`
	UpdatedAt   int64   `bson:"updated_at"`
}

// impProductRepository MongoDB商品持久化实现
type impProductRepository struct {
	mongoDb           *mongo.Database
	productCollection *mongo.Collection
}

// NewProductRepository 创建MongoDB商品持久化实例
func NewProductRepository(ctx context.Context, cfg mongodb.Config) (repository.ProductRepository, error) {
	client, err := mongodb.NewMongoClient(ctx, cfg)
	if err != nil {
		return nil, err
	}
	mongoDb := client.Database(cfg.Database)
	return &impProductRepository{
		mongoDb:           mongoDb,
		productCollection: mongoDb.Collection("product"),
	}, nil
}

// Save 保存商品
func (imp *impProductRepository) Save(ctx context.Context, product *entity.Product) error {
	po := imp.toProductPO(product)
	opts := options.Update().SetUpsert(true)
	_, err := imp.productCollection.UpdateOne(
		ctx,
		bson.M{"_id": po.Id},
		bson.M{"$set": po},
		opts,
	)
	return err
}

// FindById 根据Id查找商品
func (imp *impProductRepository) FindById(ctx context.Context, id string) (*entity.Product, error) {
	var po ProductPO
	err := imp.productCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&po)
	if err != nil {
		if gerror.Is(err, mongo.ErrNoDocuments) {
			return nil, valueobject.ErrProductNotFound
		}
		return nil, err
	}
	return imp.toEntity(&po), nil
}

// FindAll 查找所有商品
func (imp *impProductRepository) FindAll(ctx context.Context) ([]*entity.Product, error) {
	cursor, err := imp.productCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pos []ProductPO
	if err = cursor.All(ctx, &pos); err != nil {
		return nil, err
	}

	products := make([]*entity.Product, len(pos))
	for i, po := range pos {
		products[i] = imp.toEntity(&po)
	}
	return products, nil
}

// Update 更新商品
func (imp *impProductRepository) Update(ctx context.Context, product *entity.Product) error {
	po := imp.toProductPO(product)
	result, err := imp.productCollection.ReplaceOne(
		ctx,
		bson.M{"_id": po.Id},
		po,
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return gerror.New("product not found")
	}
	return nil
}

// Delete 删除商品
func (imp *impProductRepository) Delete(ctx context.Context, id string) error {
	result, err := imp.productCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return gerror.New("product not found")
	}
	return nil
}

// toProductPO 将领域实体转换为商品持久化对象
func (imp *impProductRepository) toProductPO(product *entity.Product) *ProductPO {
	return &ProductPO{
		Id:          product.Id,
		Name:        product.Name,
		Description: product.Description,
		Price: MoneyPO{
			Amount:   product.Price.Amount(),
			Currency: product.Price.Currency(),
		},
		Stock:     product.Stock,
		Status:    string(product.Status),
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}

// toEntity 将持久化对象转换为领域实体
func (imp *impProductRepository) toEntity(po *ProductPO) *entity.Product {
	return entity.NewProduct(
		po.Id,
		po.Name,
		po.Description,
		valueobject.NewMoney(po.Price.Amount, po.Price.Currency),
		po.Stock,
		valueobject.ProductStatus(po.Status),
		po.CreatedAt,
		po.UpdatedAt,
	)
}
