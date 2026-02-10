package part

import (
	"context"
	"errors"

	"github.com/mllbll/space-manufacture/inventory/internal/model"
	repoConverter "github.com/mllbll/space-manufacture/inventory/internal/repository/converter"
	repoModel "github.com/mllbll/space-manufacture/inventory/internal/repository/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *mongoRepository) GetPart(ctx context.Context, req model.GetPartRequest) (model.GetPartResponse, error) {
	var res repoModel.Part

	err := r.collection.FindOne(ctx, bson.M{"uuid":req.UUID}).Decode(&res)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.GetPartResponse{}, model.ErrPartNotFound
		}
		return model.GetPartResponse{}, err
	}

	return repoConverter.GetPartResponseToModel(repoModel.GetPartResponse{Parts: res}), nil
}
