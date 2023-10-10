package repository

import (
	"chess-backend/internal/domain"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RepoSteps struct {
	MongoDB *mongo.Client
}

func NewRepoStep(mongoDB *mongo.Client) *RepoSteps {
	return &RepoSteps{
		MongoDB: mongoDB,
	}
}

func (r *RepoSteps) StepsSave(book domain.InfoStep) error {

	var findBook domain.InfoStep
	coll := r.MongoDB.Database("admin").Collection("Books")

	filter := bson.D{{"Name", book.Name}, {"Party", book.Party}, {"NumberParty", book.NumberParty}}
	err := coll.FindOneAndReplace(context.TODO(), filter, book, nil).Decode(&findBook)
	if err != nil {
		log.Println(err)
	}

	// coll.FindOneAndUpdate(context.TODO(), filter, book, nil).Decode(&findBook)
	ll := len(findBook.Name)
	if ll == 0 {
		_, err := coll.InsertOne(context.TODO(), book, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RepoSteps) GetAllSteps() ([]domain.InfoStep, error) {
	// filter := bson.D{{"Main", true}}
	filter := bson.D{{}}
	coll := r.MongoDB.Database("ChessBoard").Collection("Parts")
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var infoSTtt []domain.InfoStep
	for cursor.Next(context.Background()) {

		var iifo domain.InfoStep
		if err = cursor.Decode(&iifo); err != nil {
			log.Fatal(err)
		}
		infoSTtt = append(infoSTtt, iifo)
	}
	return infoSTtt, nil
}
