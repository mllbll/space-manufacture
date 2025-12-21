package v1

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/mllbll/space-manufacture/inventory/internal/converter"
	"github.com/mllbll/space-manufacture/inventory/internal/model"
	inventoryV1 "github.com/mllbll/space-manufacture/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *APISuite) TestGetPartSuccess() {
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

		getPartRequest = &inventoryV1.GetPartRequest{
			Uuid: uuid,
		}

		manufacture = &inventoryV1.Manufacturer{
			Name:    manufacturer_name,
			Country: manufacturer_country,
			Website: manufacturer_website,
		}

		demensions = &inventoryV1.Dimensions{
			Length: dimensions_length,
			Width:  dimensions_width,
			Height: dimensions_height,
			Weight: dimensions_weight,
		}

		metadata = map[string]*inventoryV1.Value{
			gofakeit.Word(): {
				Types: &inventoryV1.Value_StringValue{
					StringValue: gofakeit.Word(),
				},
			},
			gofakeit.Word(): {
				Types: &inventoryV1.Value_BoolValue{
					BoolValue: gofakeit.Bool(),
				},
			},
			gofakeit.Word(): {
				Types: &inventoryV1.Value_DoubleValue{
					DoubleValue: gofakeit.Float64(),
				},
			},
			gofakeit.Word(): {
				Types: &inventoryV1.Value_Int64Value{
					Int64Value: gofakeit.Int64(),
				},
			},
		}

		part = &inventoryV1.Part{
			Uuid:          uuid,
			Name:          name,
			Description:   description,
			Price:         price,
			StockQuantity: stock_quantity,
			Category:      inventoryV1.Category(category),
			Dimensions:    demensions,
			Manufacturer:  manufacture,
			Tags:          tags,
			Metadata:      metadata,
			CreatedAt:     timestamppb.New(createdAt),
			UpdatedAt:     timestamppb.New(updatedAt),
		}

		expectedGetPartRequest = converter.GetPartRequestToModel(getPartRequest)
	)

	// Конвертируем proto Part в model Part для мока
	expectedModelPart := converter.PartToModel(part)
	expectedModelResponse := model.GetPartResponse{
		Parts: expectedModelPart,
	}

	// Мок должен ожидать model.GetPartRequest и возвращать model.GetPartResponse
	s.inventoryService.On("GetPart", s.ctx, expectedGetPartRequest).Return(expectedModelResponse, nil)

	res, err := s.api.GetPart(s.ctx, getPartRequest)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(part, res.GetParts())
}

func (s *APISuite) TestGetPartNotFound() {
	var (
		uuid = gofakeit.UUID()

		getPartRequest = &inventoryV1.GetPartRequest{
			Uuid: uuid,
		}

		expectedGetPartRequest = converter.GetPartRequestToModel(getPartRequest)
	)
	s.inventoryService.On("GetPart", s.ctx, expectedGetPartRequest).Return(model.GetPartResponse{}, model.ErrPartNotFound)

	res, err := s.api.GetPart(s.ctx, getPartRequest)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())
}

func (s *APISuite) TestGetPartError() {
	var (
		serviceErr = gofakeit.Error()
		uuid       = gofakeit.UUID()

		getPartRequest = &inventoryV1.GetPartRequest{
			Uuid: uuid,
		}

		expectedGetPartRequest = converter.GetPartRequestToModel(getPartRequest)
	)

	s.inventoryService.On("GetPart", s.ctx, expectedGetPartRequest).Return(model.GetPartResponse{}, serviceErr)

	res, err := s.api.GetPart(s.ctx, getPartRequest)
	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().Nil(res)

}
