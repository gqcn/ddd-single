package mongodb

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"main/internal/domain/order/entity"
	"main/internal/domain/order/repository"
	"main/internal/domain/order/valueobject"
	"main/utility/mongodb"
)

// OrderPO 订单持久化对象
type OrderPO struct {
	Id          string        `bson:"_id"`
	UserId      string        `bson:"user_id"`
	Items       []OrderItemPO `bson:"items"`
	TotalAmount MoneyPO       `bson:"total_amount"`
	Status      string        `bson:"status"`
	Remark      string        `bson:"remark"`
	CreatedAt   int64         `bson:"created_at"`
	UpdatedAt   int64         `bson:"updated_at"`
}

// OrderItemPO 订单项持久化对象
type OrderItemPO struct {
	ProductId string  `bson:"product_id"`
	Quantity  int     `bson:"quantity"`
	Price     MoneyPO `bson:"price"`
}

// MoneyPO Money值对象的持久化对象
type MoneyPO struct {
	Amount   float64 `bson:"amount"`
	Currency string  `bson:"currency"`
}

// impOrderRepository MongoDB订单持久化实现
type impOrderRepository struct {
	mongoDb         *mongo.Database
	orderCollection *mongo.Collection
}

// NewOrderRepository 创建MongoDB订单持久化实例
func NewOrderRepository(ctx context.Context, cfg mongodb.Config) (repository.OrderRepository, error) {
	client, err := mongodb.NewMongoClient(ctx, cfg)
	if err != nil {
		return nil, err
	}
	mongoDb := client.Database(cfg.Database)
	return &impOrderRepository{
		mongoDb:         mongoDb,
		orderCollection: mongoDb.Collection("order"),
	}, nil
}

// Save 保存订单
func (imp *impOrderRepository) Save(ctx context.Context, order *entity.Order) error {
	po := imp.toOrderPO(order)
	opts := options.Update().SetUpsert(true)
	_, err := imp.orderCollection.UpdateOne(
		ctx,
		bson.M{"_id": po.Id},
		bson.M{"$set": po},
		opts,
	)
	return err
}

// FindById 根据Id查找订单
func (imp *impOrderRepository) FindById(ctx context.Context, id string) (*entity.Order, error) {
	var po OrderPO
	err := imp.orderCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&po)
	if err != nil {
		if gerror.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return imp.toEntity(&po), nil
}

// FindByUserId 查找用户的所有订单
func (imp *impOrderRepository) FindByUserId(ctx context.Context, userId string) ([]*entity.Order, error) {
	cursor, err := imp.orderCollection.Find(ctx, bson.M{"user_id": userId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pos []OrderPO
	if err = cursor.All(ctx, &pos); err != nil {
		return nil, err
	}

	orders := make([]*entity.Order, len(pos))
	for index, po := range pos {
		orders[index] = imp.toEntity(&po)
	}

	return orders, nil
}

// toOrderPO 将领域实体转换为订单持久化对象
func (imp *impOrderRepository) toOrderPO(order *entity.Order) *OrderPO {
	items := make([]OrderItemPO, len(order.Items))
	for i, item := range order.Items {
		items[i] = OrderItemPO{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
			Price: MoneyPO{
				Amount:   item.Price.Amount(),
				Currency: item.Price.Currency(),
			},
		}
	}

	return &OrderPO{
		Id:     order.Id,
		UserId: order.UserId,
		Items:  items,
		TotalAmount: MoneyPO{
			Amount:   order.TotalAmount.Amount(),
			Currency: order.TotalAmount.Currency(),
		},
		Status:    string(order.Status),
		Remark:    order.Remark,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}
}

// toEntity 将持久化对象转换为领域实体
func (imp *impOrderRepository) toEntity(po *OrderPO) *entity.Order {
	items := make([]*entity.OrderItem, len(po.Items))
	for i, item := range po.Items {
		items[i] = &entity.OrderItem{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
			Price:     valueobject.NewMoney(item.Price.Amount, item.Price.Currency),
		}
	}

	order := &entity.Order{
		Id:          po.Id,
		UserId:      po.UserId,
		Items:       items,
		TotalAmount: valueobject.NewMoney(po.TotalAmount.Amount, po.TotalAmount.Currency),
		Status:      valueobject.OrderStatus(po.Status),
		Remark:      po.Remark,
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
	}

	return order
}

// Update updates an existing order
func (imp *impOrderRepository) Update(ctx context.Context, order *entity.Order) error {
	po := imp.toOrderPO(order)
	_, err := imp.orderCollection.ReplaceOne(
		ctx,
		bson.M{"_id": po.Id},
		po,
	)
	return err
}

// Delete removes an order from storage
func (imp *impOrderRepository) Delete(ctx context.Context, id string) error {
	_, err := imp.orderCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
