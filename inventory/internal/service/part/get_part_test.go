package part

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/mllbll/space-manufacture/inventory/internal/model"
	"github.com/samber/lo"
)

func (s *ServiceSuite) TestGetPartSuccess() {
	var (
		uuid           = gofakeit.UUID()
		name           = gofakeit.Word()
		description    = gofakeit.Paragraph(1, 3, 3)
		price          = gofakeit.Float64()
		stock_quantity = gofakeit.Int64()
		category       = gofakeit.Number(0, 4)
		tags           = []string{gofakeit.Word(), gofakeit.Word()}
		createdAt      = gofakeit.Date()
		updatedAt      = gofakeit.Date()

		manufacturer_name    = gofakeit.BeerName()
		manufacturer_country = gofakeit.Country()
		manufacturer_website = gofakeit.URL()

		dimensions_length = gofakeit.Float64()
		dimensions_width  = gofakeit.Float64()
		dimensions_height = gofakeit.Float64()
		dimensions_weight = gofakeit.Float64()

		getPartRequest = model.GetPartRequest{
			UUID: uuid,
		}

		manufacture = model.Manufacturer{
			Name:    manufacturer_name,
			Country: manufacturer_country,
			Website: manufacturer_website,
		}

		demensions = model.Dimensions{
			Length: dimensions_length,
			Width:  dimensions_width,
			Height: dimensions_height,
			Weight: dimensions_weight,
		}

		metadata = map[string]model.Value{
			gofakeit.Word(): model.StringValue{Value: gofakeit.Word()},
			gofakeit.Word(): model.BoolValue{Value: gofakeit.Bool()},
			gofakeit.Word(): model.DoubleValue{Value: gofakeit.Float64()},
			gofakeit.Word(): model.Int64Value{Value: gofakeit.Int64()},
		}

		part = model.Part{
			UUID:           uuid,
			Name:           name,
			Description:    description,
			Price:          price,
			Stock_quantity: stock_quantity,
			Category:       model.Category(category),
			Dimensions:     demensions,
			Manufacturer:   manufacture,
			Tags:           tags,
			Metadata:       metadata,
			CreatedAt:      lo.ToPtr(createdAt),
			UpdatedAt:      lo.ToPtr(updatedAt),
		}
	)

	s.inventoryRepository.On("GetPart", s.ctx, getPartRequest).Return(model.GetPartResponse{Parts: part}, nil)

	res, err := s.service.GetPart(s.ctx, getPartRequest)
	s.NoError(err)
	s.Equal(part, res.Parts)
}

func (s *ServiceSuite) TestGetPartError() {
	var (
		uuid    = gofakeit.UUID()
		repoErr = gofakeit.Error()

		getPartRequest = model.GetPartRequest{
			UUID: uuid,
		}
	)

	s.inventoryRepository.On("GetPart", s.ctx, getPartRequest).Return(model.GetPartResponse{}, repoErr)

	res, err := s.service.GetPart(s.ctx, getPartRequest)

	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(res)

}
