package mongodb

import (
	"github.com/alonyb/url_shortener/shortener"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type mongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

func (r mongoRepository) Find(code string) (*shortener.Redirect, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	redirect := &shortener.Redirect{}
	collection := r.client.Database(r.database).Collection("redirects")
	filter := bson.M{"code":code}
	err := collection.FindOne(ctx, filter).Decode(&redirect)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(shortener.ErrRedirectNotFound, "repository.mongodb")
		}
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}
	return redirect, nil
}

func (r mongoRepository) Store(redirect *shortener.Redirect) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	//tal vez se puede sacar en otro metodo
	collection := r.client.Database(r.database).Collection("redirects")
	_, err := collection.InsertOne(
		ctx,
		bson.M{
			"code":	redirect.Code,
			"url":	redirect.URL,
			"created_at":	redirect.CreateAt,
		})
	if err != nil {
		log.Println(err)
		return errors.Wrap(err, "repository.Redirect.Store")
	}
	return nil
}

func newMongoClient(mongourl string, mongoTimeOut int) (*mongo.Client, error){


	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeOut))
	defer cancel()
	clientOptions := options.Client().ApplyURI(mongourl)
	client, err := mongo.Connect(ctx, clientOptions)
	//client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongourl))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	return client, err
}

func NewMongoRepository(mongoUrl, mongodb string, mongoTimeout int) (shortener.RedirectRepository, error){
	repo := &mongoRepository{
		timeout: time.Duration(mongoTimeout) * time.Second,
		database: mongodb,
	}
	client, err := newMongoClient(mongoUrl, mongoTimeout)
	if err != nil {
		log.Println("no database connection", err)
		return nil, errors.Wrap(err, "repository.NewMongoRepository")
	}
	repo.client = client
	return repo, nil
}