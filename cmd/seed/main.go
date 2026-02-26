package main

import (
	"context"
	"log"
	"time"

	"api-quest/config"
	"api-quest/internal/infrastructure/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type bookDoc struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title"`
	Author    string             `bson:"author"`
	ISBN      string             `bson:"isbn"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

var seeds = []bookDoc{
	{
		ID:        primitive.NewObjectID(),
		Title:     "The Go Programming Language",
		Author:    "Alan Donovan",
		ISBN:      "978-0134190440",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        primitive.NewObjectID(),
		Title:     "Clean Architecture",
		Author:    "Robert Martin",
		ISBN:      "978-0134494166",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        primitive.NewObjectID(),
		Title:     "Domain-Driven Design",
		Author:    "Eric Evans",
		ISBN:      "978-0321125217",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        primitive.NewObjectID(),
		Title:     "Designing Data-Intensive Applications",
		Author:    "Martin Kleppmann",
		ISBN:      "978-1449373320",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        primitive.NewObjectID(),
		Title:     "The Pragmatic Programmer",
		Author:    "David Thomas",
		ISBN:      "978-0135957059",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        primitive.NewObjectID(),
		Title:     "Concurrency in Go",
		Author:    "Katherine Cox-Buday",
		ISBN:      "978-1491941195",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        primitive.NewObjectID(),
		Title:     "Building Microservices",
		Author:    "Sam Newman",
		ISBN:      "978-1492034025",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        primitive.NewObjectID(),
		Title:     "Refactoring",
		Author:    "Martin Fowler",
		ISBN:      "978-0134757599",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}

func main() {
	cfg := config.Load()

	client, err := mongodb.NewClient(cfg.MongoURI)
	if err != nil {
		log.Fatalf("MongoDB connection failed: %v", err)
	}
	defer client.Disconnect(context.Background())

	col := client.Database(cfg.DBName).Collection("books")
	ctx := context.Background()

	// Drop existing data
	if err := col.Drop(ctx); err != nil {
		log.Fatalf("Failed to drop collection: %v", err)
	}
	log.Println("Dropped existing books collection.")

	// Insert seed data
	docs := make([]interface{}, len(seeds))
	for i, s := range seeds {
		docs[i] = s
	}

	result, err := col.InsertMany(ctx, docs)
	if err != nil {
		log.Fatalf("Failed to insert seed data: %v", err)
	}

	log.Printf("Seeded %d books successfully.\n", len(result.InsertedIDs))

	// Print IDs for easy reference
	for i, id := range result.InsertedIDs {
		log.Printf("  [%d] %-42s — %s", i+1, seeds[i].Title, id.(primitive.ObjectID).Hex())
	}

	// Ensure index on author
	_, _ = col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "author", Value: 1}},
	})
	log.Println("Index on 'author' ensured.")
}
