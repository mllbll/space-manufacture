package part

import (
	def "github.com/mllbll/space-manufacture/inventory/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ def.InventoryRepository = (*mongoRepository)(nil)

const partsCollection = "parts"

type mongoRepository struct {
	collection *mongo.Collection
}

func NewMongoCollection (db *mongo.Database) *mongoRepository {
	return &mongoRepository{
		collection: db.Collection(partsCollection),
	}
}
