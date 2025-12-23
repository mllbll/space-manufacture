package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/mllbll/space-manufacture/order/internal/model"
	"github.com/samber/lo"
	"github.com/stretchr/testify/mock"
)

func (s *ServiceSuite) TestCreateSuccess() {
	var (
		userUUID = gofakeit.UUID()
		partUUIDs = []string{gofakeit.UUID()}

		createOrderRequest = model.CreateOrderRequest{
			UserUUID: userUUID,
			PartUUIDs: partUUIDs,
		}

		getPriceMessage = model.PatrsFilter{
			UUIDs: partUUIDs,
			Names: nil,
			Categories: nil,
			Manufacture_contries: nil,
			Tags: nil,
		}

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

	s.inventoryClient.On("ListParts", s.ctx, getPriceMessage).Return(model.ListPartsResponse{Parts: []model.Part{part}}, nil)

	//s.orderRepository.On("Create", s.ctx, createOrderRequest).Return(nil)
	s.orderRepository.On("Create", s.ctx, mock.MatchedBy(func(order model.Order) bool {
		return order.UserUUID == userUUID &&
			len(order.PartUUIDs) == len(partUUIDs) &&
			order.PartUUIDs[0] == partUUIDs[0] &&
			order.TotalPrice == float32(price) &&
			order.Status == "PENDING_PAYMENT" &&
			order.OrderUUID != ""
	})).Return(nil)

	res, err := s.service.Create(s.ctx, createOrderRequest)

	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().NotEmpty(res.OrderUUID)
	s.Require().Equal(float32(price), res.TotalPrice) 
}

func (s *ServiceSuite) TestCreateClientError() {
	var (
		userUUID = gofakeit.UUID()
		partUUIDs = []string{gofakeit.UUID()}


		createOrderRequest = model.CreateOrderRequest{
			UserUUID: userUUID,
			PartUUIDs: partUUIDs,
		}

		getPriceMessage = model.PatrsFilter{
			UUIDs: partUUIDs,
			Names: nil,
			Categories: nil,
			Manufacture_contries: nil,
			Tags: nil,
		}
	)

	s.inventoryClient.On("ListParts", s.ctx, getPriceMessage).Return(model.ListPartsResponse{}, model.ErrOrderNotFound)

	res, err := s.service.Create(s.ctx, createOrderRequest)

	s.Require().Error(err)
	s.Require().Empty(res)
	s.Require().Equal(err, model.ErrOrderNotFound)
}

