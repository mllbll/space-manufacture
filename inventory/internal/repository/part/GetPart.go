package part

import (
	"context"

	repoConverter "github.com/mllbll/space-manufacture/inventory/internal/repository/converter"
	repoModel "github.com/mllbll/space-manufacture/inventory/internal/repository/model"
	"github.com/mllbll/space-manufacture/inventory/internal/model"
	
)

func (r *repository) GetPart(_ context.Context, req model.GetPartRequest) (model.GetPartResponse, error) {
	r.mu.Lock()
	defer r.mu.Unlock()


	part, ok := r.data[req.UUID]
	if !ok {
		return model.GetPartResponse{}, model.ErrPartNotFound
	}

	return repoConverter.GetPartResponseToModel(repoModel.GetPartResponse{
		Parts: part,
	}), nil

}


