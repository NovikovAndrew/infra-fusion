package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Hosts      []string
	DbName     string
	Login      string
	Password   string
	AuthSource string
}

type Database interface {
	Collection(name string, opts ...*options.CollectionOptions) Collection
	CloseConnection(ctx context.Context) error
}

type Collection interface {
	InsertOne(
		ctx context.Context,
		document interface{},
		opts ...*options.InsertOneOptions,
	) (*mongo.InsertOneResult, error)
	InsertMany(
		ctx context.Context,
		documents []interface{},
		opts ...*options.InsertManyOptions,
	) (*mongo.InsertManyResult, error)
	Find(
		ctx context.Context,
		filter interface{},
		opts ...*options.FindOptions,
	) (cur *mongo.Cursor, err error)
	FindOne(
		ctx context.Context,
		filter interface{},
		opts ...*options.FindOneOptions,
	) *mongo.SingleResult
	FindOneAndUpdate(
		ctx context.Context,
		filter interface{},
		update interface{},
		opts ...*options.FindOneAndUpdateOptions,
	) *mongo.SingleResult
	UpdateOne(
		ctx context.Context,
		filter interface{},
		update interface{},
		opts ...*options.UpdateOptions,
	) (*mongo.UpdateResult, error)
	UpdateMany(
		ctx context.Context,
		filter interface{},
		update interface{},
		opts ...*options.UpdateOptions,
	) (*mongo.UpdateResult, error)
	CountDocuments(
		ctx context.Context,
		filter interface{},
		opts ...*options.CountOptions,
	) (int64, error)
}

type DatabaseMongo struct {
	*mongo.Database
}

type CollectionMongo struct {
	*mongo.Collection
}

func GetDatabase(ctx context.Context, cfg Config) (Database, error) {
	ops := options.Client().SetHosts(cfg.Hosts)

	if cfg.Login != "" && cfg.Password != "" {
		ops.SetAuth(options.Credential{
			AuthSource: cfg.AuthSource,
			Username:   cfg.Login,
			Password:   cfg.Password,
		})
	}

	client, err := mongo.Connect(ctx, ops)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &DatabaseMongo{client.Database(cfg.DbName)}, nil
}

func (d *DatabaseMongo) CloseConnection(ctx context.Context) error {
	return d.Client().Disconnect(ctx)
}

func (d *DatabaseMongo) Collection(name string, opts ...*options.CollectionOptions) Collection {
	return &CollectionMongo{d.Database.Collection(name, opts...)}
}

func (c *CollectionMongo) InsertOne(
	ctx context.Context,
	document interface{},
	opts ...*options.InsertOneOptions,
) (*mongo.InsertOneResult, error) {
	return c.Collection.InsertOne(ctx, document, opts...)
}

func (c *CollectionMongo) InsertMany(
	ctx context.Context,
	documents []interface{},
	opts ...*options.InsertManyOptions,
) (*mongo.InsertManyResult, error) {
	return c.Collection.InsertMany(ctx, documents, opts...)
}

func (c *CollectionMongo) Find(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOptions,
) (cur *mongo.Cursor, err error) {
	return c.Collection.Find(ctx, filter, opts...)
}

func (c *CollectionMongo) FindOne(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) *mongo.SingleResult {
	return c.Collection.FindOne(ctx, filter, opts...)
}

func (c *CollectionMongo) FindOneAndUpdate(
	ctx context.Context,
	filter interface{},
	update interface{},
	opts ...*options.FindOneAndUpdateOptions,
) *mongo.SingleResult {
	return c.Collection.FindOneAndUpdate(ctx, filter, update, opts...)
}

func (c *CollectionMongo) UpdateOne(
	ctx context.Context,
	filter interface{},
	update interface{},
	opts ...*options.UpdateOptions,
) (*mongo.UpdateResult, error) {
	return c.Collection.UpdateOne(ctx, filter, update, opts...)
}

func (c *CollectionMongo) UpdateMany(
	ctx context.Context,
	filter interface{},
	update interface{},
	opts ...*options.UpdateOptions,
) (*mongo.UpdateResult, error) {
	return c.Collection.UpdateMany(ctx, filter, update, opts...)
}

func (c *CollectionMongo) CountDocuments(
	ctx context.Context,
	filter interface{},
	opts ...*options.CountOptions,
) (int64, error) {
	return c.Collection.CountDocuments(ctx, filter, opts...)
}
