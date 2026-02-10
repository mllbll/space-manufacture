package part

import (
	"context"
	"log"

	"github.com/mllbll/space-manufacture/inventory/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	repoConverter "github.com/mllbll/space-manufacture/inventory/internal/repository/converter"
	repoModel "github.com/mllbll/space-manufacture/inventory/internal/repository/model"
)

func (r *mongoRepository) ListParts(ctx context.Context, req model.ListPartsRequest) (model.ListPartsResponse, error) {
	filter := buildFilter(req.Filter)

	cursor, err := r.collection.Find(ctx, filter, options.Find())
	if err != nil {
		log.Printf("Ошибка при чтении заметок %v\n", err)
		return model.ListPartsResponse{}, err
	}

	defer func () {
		err := cursor.Close(ctx);
		if err != nil {
			log.Printf("Ошибка при закрытии курсора: %v\n", err)
		}
	}()

	var parts []repoModel.Part
	if err = cursor.All(ctx, &parts); err != nil {
		log.Printf("Ошибка при декодировании заметок: %v\n", err)
		return model.ListPartsResponse{}, err
	}

	return repoConverter.ListPartsResponseToModel(repoModel.ListPartsResponse{Parts: parts}), nil
}

func buildFilter(filter model.PatrsFilter) bson.M {
	var conditions []bson.M
	if len(filter.UUIDs) > 0 {
		conditions = append(conditions, bson.M{"uuid": bson.M{"$in": filter.UUIDs}})
	}
	if len(filter.Names) > 0 {
		conditions = append(conditions, bson.M{"name": bson.M{"$in": filter.Names}})
	}
	if len(filter.Categories) > 0 {
		categoriesToInt := make([]int32, len(filter.Categories))
		for id, val := range filter.Categories {
			categoriesToInt[id] = int32(val)
		}
		conditions = append(conditions, bson.M{"category": bson.M{"$in": categoriesToInt}})
	}
	if len(filter.Manufacture_contries) > 0 {
		conditions = append(conditions, bson.M{"manufacturer.contry": bson.M{"$in": filter.Manufacture_contries}})
	}
	if len(filter.Tags) > 0 {
		conditions = append(conditions, bson.M{"tags": bson.M{"$in": filter.Tags}})
	}
	if len(conditions) == 0 {
		return bson.M{}
	}
	return bson.M{"$and":conditions}
}
