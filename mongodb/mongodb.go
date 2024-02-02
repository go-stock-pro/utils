package mongodb

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rohanraj7316/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	collOptions *options.CollectionOptions

	DB *mongo.Database
}

func New(ctx context.Context, config ...Config) (*Storage, error) {
	err := logger.Configure()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cfg, err := configDefault(config...)
	if err != nil {
		return nil, err
	}

	client, err := mongo.Connect(ctx, cfg.cOptions)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	db := client.Database(cfg.DbName, cfg.dOptions)

	return &Storage{collOptions: cfg.collOptions, DB: db}, nil
}

func (s *Storage) Find(ctx context.Context, filter, dst any,
	opts ...*options.FindOptions) error {

	coll, err := s.getCollection(filter)
	if err != nil {
		return errors.WithStack(err)
	}

	cur, err := coll.Find(ctx, filter, opts...)
	if err != nil {
		return errors.WithStack(err)
	}
	defer cur.Close(ctx)

	err = cur.All(ctx, dst)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Storage) Create(ctx context.Context, doc any,
	opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {

	coll, err := s.getCollection(doc)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = s.IsValidStruct(doc)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	dst, err := coll.InsertOne(ctx, doc, opts...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return dst, nil
}

func (s *Storage) Update(ctx context.Context, filter, update any,
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {

	coll, err := s.getCollection(filter)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	opts = append(opts, (&options.UpdateOptions{}).SetUpsert(true))

	dst, err := coll.UpdateMany(ctx, filter, update, opts...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return dst, nil
}

func (s *Storage) Delete(ctx context.Context, filter any,
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {

	coll, err := s.getCollection(filter)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	dst, err := coll.DeleteMany(ctx, filter, opts...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return dst, nil
}
