package part

import (
	"context"

	"github.com/mllbll/space-manufacture/inventory/internal/model"
	repoConverter "github.com/mllbll/space-manufacture/inventory/internal/repository/converter"
	repoModel "github.com/mllbll/space-manufacture/inventory/internal/repository/model"
)

func (r *repository) ListParts(_ context.Context, req model.ListPartsRequest) (model.ListPartsResponse, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var parts []repoModel.Part
	for _, part := range r.data {
		p := part
		parts = append(parts, p)
	}

	resp := repoConverter.ListPartsResponseToModel(repoModel.ListPartsResponse{
		Parts: parts,
	})

	return resp, nil

}
