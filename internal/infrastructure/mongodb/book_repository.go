package mongodb

import (
	"context"
	"errors"
	"time"

	"api-quest/internal/domain/book"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type bookDoc struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title"`
	Author    string             `bson:"author"`
	ISBN      string             `bson:"isbn"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func toDoc(b *book.Book) *bookDoc {
	doc := &bookDoc{
		Title:     b.Title,
		Author:    b.Author,
		ISBN:      b.ISBN,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
	if b.ID != "" {
		if oid, err := primitive.ObjectIDFromHex(b.ID); err == nil {
			doc.ID = oid
		}
	}
	return doc
}

func fromDoc(doc *bookDoc) *book.Book {
	return &book.Book{
		ID:        doc.ID.Hex(),
		Title:     doc.Title,
		Author:    doc.Author,
		ISBN:      doc.ISBN,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}
}

type BookRepository struct {
	col *mongo.Collection
}

func NewBookRepository(db *mongo.Database) *BookRepository {
	col := db.Collection("books")

	// Index on author for search performance (Level 6 + 8)
	_, _ = col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{Key: "author", Value: 1}},
	})

	return &BookRepository{col: col}
}

func (r *BookRepository) Create(ctx context.Context, b *book.Book) (*book.Book, error) {
	doc := toDoc(b)
	doc.ID = primitive.NewObjectID()

	_, err := r.col.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}
	return fromDoc(doc), nil
}

func (r *BookRepository) FindAll(ctx context.Context, filter book.Filter) ([]*book.Book, int64, error) {
	f := bson.D{}
	if filter.Author != "" {
		f = append(f, bson.E{Key: "author", Value: primitive.Regex{Pattern: filter.Author, Options: "i"}})
	}

	total, err := r.col.CountDocuments(ctx, f)
	if err != nil {
		return nil, 0, err
	}

	skip := int64((filter.Page - 1) * filter.Limit)
	opts := options.Find().SetSkip(skip).SetLimit(int64(filter.Limit))

	cursor, err := r.col.Find(ctx, f, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var books []*book.Book
	for cursor.Next(ctx) {
		var doc bookDoc
		if err := cursor.Decode(&doc); err != nil {
			return nil, 0, err
		}
		books = append(books, fromDoc(&doc))
	}

	if books == nil {
		books = []*book.Book{}
	}
	return books, total, nil
}

func (r *BookRepository) FindByID(ctx context.Context, id string) (*book.Book, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, book.ErrNotFound
	}

	var doc bookDoc
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, book.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return fromDoc(&doc), nil
}

func (r *BookRepository) Update(ctx context.Context, id string, b *book.Book) (*book.Book, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, book.ErrNotFound
	}

	update := bson.M{
		"$set": bson.M{
			"title":      b.Title,
			"author":     b.Author,
			"isbn":       b.ISBN,
			"updated_at": b.UpdatedAt,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var doc bookDoc
	err = r.col.FindOneAndUpdate(ctx, bson.M{"_id": oid}, update, opts).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, book.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return fromDoc(&doc), nil
}

func (r *BookRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return book.ErrNotFound
	}

	result, err := r.col.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return book.ErrNotFound
	}
	return nil
}
