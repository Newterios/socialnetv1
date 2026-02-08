package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoDB(uri string) (*MongoDB, error) {
	if uri == "" {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB")

	db := client.Database("socialnet")

	mongodb := &MongoDB{
		Client:   client,
		Database: db,
	}

	if err := mongodb.ensureIndexes(); err != nil {
		log.Printf("Warning: failed to create MongoDB indexes: %v", err)
	}

	return mongodb, nil
}

func (m *MongoDB) ensureIndexes() error {
	ctx := context.Background()

	_, err := m.Database.Collection("emoji_stats").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "emoji", Value: 1}},
	})
	if err != nil {
		return err
	}

	_, err = m.Database.Collection("user_status").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: 1}},
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoDB) Close() error {
	if m == nil || m.Client == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.Client.Disconnect(ctx)
}

func (m *MongoDB) EmojiStats() *mongo.Collection {
	return m.Database.Collection("emoji_stats")
}

func (m *MongoDB) UserStatus() *mongo.Collection {
	return m.Database.Collection("user_status")
}

func (m *MongoDB) SiteAnalytics() *mongo.Collection {
	return m.Database.Collection("site_analytics")
}
