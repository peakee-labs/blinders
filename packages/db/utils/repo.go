package dbutils

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IRawModel interface {
	SetID(primitive.ObjectID)
	SetInitTimeByNow()
	SetUpdatedAtByNow()
}

type RawModel struct {
	ID        primitive.ObjectID `bson:"_id"       json:"id"`
	CreatedAt primitive.DateTime `bson:"createdAt" json:"createdAt"`
	UpdatedAt primitive.DateTime `bson:"updatedAt" json:"updatedAt"`
}

func (m *RawModel) SetID(id primitive.ObjectID) {
	m.ID = id
}

func (m *RawModel) SetInitTimeByNow() {
	now := primitive.NewDateTimeFromTime(time.Now())
	m.CreatedAt = now
	m.UpdatedAt = now
}

func (m *RawModel) SetUpdatedAtByNow() {
	now := primitive.NewDateTimeFromTime(time.Now())
	m.UpdatedAt = now
}

type IRepo[M IRawModel] interface {
	Insert(obj M) (M, error)
	InsertRaw(obj M) (M, error)
	GetByID(ID primitive.ObjectID) (M, error)
	UpdateByID(ID primitive.ObjectID, obj M) (M, error)
}

type SingleCollectionRepo[M IRawModel] struct {
	*mongo.Collection
}

func (r *SingleCollectionRepo[M]) Insert(obj M) (M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := r.InsertOne(ctx, obj)
	return obj, err
}

func (r *SingleCollectionRepo[M]) InsertRaw(obj M) (M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	obj.SetID(primitive.NewObjectID())
	obj.SetInitTimeByNow()

	_, err := r.InsertOne(ctx, obj)
	return obj, err
}

func (r *SingleCollectionRepo[M]) GetByID(ID primitive.ObjectID) (M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var obj M
	err := r.FindOne(ctx, bson.M{"_id": ID}).Decode(&obj)

	return obj, err
}
