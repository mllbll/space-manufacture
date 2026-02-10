package part

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/mllbll/space-manufacture/inventory/internal/model"
	"github.com/samber/lo"
)

func (s *ServiceSuite) TestListPartsSuccess() {
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

		uuids = []string{uuid}
		names = []string{name}
		categories = []model.Category{model.Category(category)}
		manufacturer_countries = []string{manufacturer_country}


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

		partsFilter = model.PatrsFilter{
			UUIDs: uuids,
			Names: names,
			Categories: categories,
			Manufacture_contries: manufacturer_countries,
			Tags: tags,
		}

		listPartsResponse = model.ListPartsResponse{Parts: []model.Part{part}}

		listPartsRequest = model.ListPartsRequest{
			Filter: partsFilter,
		}
	)

	s.inventoryRepository.On("ListParts", s.ctx, listPartsRequest).Return(listPartsResponse, nil)

	res, err := s.service.ListParts(s.ctx, listPartsRequest)

	s.Require().NoError(err)
	s.Require().Equal(listPartsResponse, res)
}

func (s *ServiceSuite) TestListPartsWithEmptyFilter() {
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

		part1 = model.Part{
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

		part2 = model.Part{
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


		listPartsResponse = model.ListPartsResponse{Parts: []model.Part{part1, part2}}

		listPartsRequest = model.ListPartsRequest{
			Filter: model.PatrsFilter{},
		}
	)

	s.inventoryRepository.On("ListParts", s.ctx, listPartsRequest).Return(listPartsResponse, nil)

	res, err := s.service.ListParts(s.ctx, listPartsRequest)

	s.Require().NoError(err)
	s.Require().Equal(res, listPartsResponse)
}

func (s *ServiceSuite) TestListPartsError() {
	var (
		uuid           = gofakeit.UUID()
		name           = gofakeit.Word()
		category       = gofakeit.Number(0, 4)
		tags           = []string{gofakeit.Word(), gofakeit.Word()}
		manufacturer_country = gofakeit.Country()
		repoErr = gofakeit.Error()

		uuids = []string{uuid}
		names = []string{name}
		categories = []model.Category{model.Category(category)}
		manufacturer_countries = []string{manufacturer_country}

		partsFilter = model.PatrsFilter{
			UUIDs: uuids,
			Names: names,
			Categories: categories,
			Manufacture_contries: manufacturer_countries,
			Tags: tags,
		}

		listPartsRequest = model.ListPartsRequest{
			Filter: partsFilter,
		}
	)

	s.inventoryRepository.On("ListParts", s.ctx, listPartsRequest).Return(model.ListPartsResponse{}, repoErr)

	res, err := s.service.ListParts(s.ctx, listPartsRequest)

	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(res)
}

func (s *ServiceSuite) TestListPartsNotFound() {
	var (
		uuid           = gofakeit.UUID()
		name           = gofakeit.Word()
		category       = gofakeit.Number(0, 4)
		tags           = []string{gofakeit.Word(), gofakeit.Word()}
		manufacturer_country = gofakeit.Country()

		uuids = []string{uuid}
		names = []string{name}
		categories = []model.Category{model.Category(category)}
		manufacturer_countries = []string{manufacturer_country}

		partsFilter = model.PatrsFilter{
			UUIDs: uuids,
			Names: names,
			Categories: categories,
			Manufacture_contries: manufacturer_countries,
			Tags: tags,
		}

		listPartsRequest = model.ListPartsRequest{
			Filter: partsFilter,
		}
	)

	s.inventoryRepository.On("ListParts", s.ctx, listPartsRequest).Return(model.ListPartsResponse{}, nil)

	res, err := s.service.ListParts(s.ctx, listPartsRequest)

	s.Require().NoError(err)
	s.Require().Empty(res)
	s.Require().Equal(res, model.ListPartsResponse{})
}

func (s *ServiceSuite) TestListPartsExcludesByUUID() {
	var (
		uuid           = gofakeit.UUID()
		filterUUID = gofakeit.UUID()
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

		uuids = []string{filterUUID}
		names = []string{name}
		categories = []model.Category{model.Category(category)}
		manufacturer_countries = []string{manufacturer_country}


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

		partsFilter = model.PatrsFilter{
			UUIDs: uuids,
			Names: names,
			Categories: categories,
			Manufacture_contries: manufacturer_countries,
			Tags: tags,
		}

		listPartsResponse = model.ListPartsResponse{Parts: []model.Part{part}}

		listPartsRequest = model.ListPartsRequest{
			Filter: partsFilter,
		}
	)

	s.inventoryRepository.On("ListParts", s.ctx, listPartsRequest).Return(listPartsResponse, nil)

	res, err := s.service.ListParts(s.ctx, listPartsRequest)

	s.Require().NoError(err)
	//s.Require().Equal(model.ListPartsResponse{}, res)
	s.Require().Equal(listPartsResponse, res)
}

func (s *ServiceSuite) TestListPartsExcludesByTags() {
	var (
		uuid           = gofakeit.UUID()
		name           = gofakeit.Word()
		description    = gofakeit.Paragraph(1, 3, 3)
		price          = gofakeit.Float64()
		stock_quantity = gofakeit.Int64()
		category       = gofakeit.Number(0, 4)
		tags           = []string{gofakeit.Word(), gofakeit.Word()}
		filterTags     = []string{gofakeit.Word(), gofakeit.Word()}
		createdAt      = gofakeit.Date()
		updatedAt      = gofakeit.Date()

		manufacturer_name    = gofakeit.BeerName()
		manufacturer_country = gofakeit.Country()
		manufacturer_website = gofakeit.URL()

		uuids = []string{uuid}
		names = []string{name}
		categories = []model.Category{model.Category(category)}
		manufacturer_countries = []string{manufacturer_country}


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

		partsFilter = model.PatrsFilter{
			UUIDs: uuids,
			Names: names,
			Categories: categories,
			Manufacture_contries: manufacturer_countries,
			Tags: filterTags,
		}

		listPartsResponse = model.ListPartsResponse{Parts: []model.Part{part}}

		listPartsRequest = model.ListPartsRequest{
			Filter: partsFilter,
		}
	)

	s.inventoryRepository.On("ListParts", s.ctx, listPartsRequest).Return(listPartsResponse, nil)

	res, err := s.service.ListParts(s.ctx, listPartsRequest)

	s.Require().NoError(err)
	//s.Require().Equal(model.ListPartsResponse{}, res)
	s.Require().Equal(listPartsResponse, res)
}

func (s *ServiceSuite) TestListPartsExcludesByManufactureCounries() {
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
		filterManufacture_country = gofakeit.Country()
		manufacturer_website = gofakeit.URL()

		uuids = []string{uuid}
		names = []string{name}
		categories = []model.Category{model.Category(category)}
		manufacturer_countries = []string{filterManufacture_country}


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

		partsFilter = model.PatrsFilter{
			UUIDs: uuids,
			Names: names,
			Categories: categories,
			Manufacture_contries: manufacturer_countries,
			Tags: tags,
		}

		listPartsResponse = model.ListPartsResponse{Parts: []model.Part{part}}

		listPartsRequest = model.ListPartsRequest{
			Filter: partsFilter,
		}
	)

	s.inventoryRepository.On("ListParts", s.ctx, listPartsRequest).Return(listPartsResponse, nil)

	res, err := s.service.ListParts(s.ctx, listPartsRequest)

	s.Require().NoError(err)
	//s.Require().Equal(model.ListPartsResponse{}, res)
	s.Require().Equal(listPartsResponse, res)
}

func (s *ServiceSuite) TestListPartsExcludesByNames() {
	var (
		uuid           = gofakeit.UUID()
		name           = gofakeit.Word()
		filterName = gofakeit.Word()
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

		uuids = []string{uuid}
		names = []string{filterName}
		categories = []model.Category{model.Category(category)}
		manufacturer_countries = []string{manufacturer_country}


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

		partsFilter = model.PatrsFilter{
			UUIDs: uuids,
			Names: names,
			Categories: categories,
			Manufacture_contries: manufacturer_countries,
			Tags: tags,
		}

		listPartsResponse = model.ListPartsResponse{Parts: []model.Part{part}}

		listPartsRequest = model.ListPartsRequest{
			Filter: partsFilter,
		}
	)

	s.inventoryRepository.On("ListParts", s.ctx, listPartsRequest).Return(listPartsResponse, nil)

	res, err := s.service.ListParts(s.ctx, listPartsRequest)

	s.Require().NoError(err)
	//s.Require().Equal(model.ListPartsResponse{}, res)
	s.Require().Equal(listPartsResponse, res)
}

func (s *ServiceSuite) TestListPartsExcludesByCategories() {
	var (
		uuid           = gofakeit.UUID()
		name           = gofakeit.Word()
		description    = gofakeit.Paragraph(1, 3, 3)
		price          = gofakeit.Float64()
		stock_quantity = gofakeit.Int64()
		category       = gofakeit.Number(0, 4)
		filterCategory = gofakeit.Number(5, 8)
		tags           = []string{gofakeit.Word(), gofakeit.Word()}
		createdAt      = gofakeit.Date()
		updatedAt      = gofakeit.Date()

		manufacturer_name    = gofakeit.BeerName()
		manufacturer_country = gofakeit.Country()
		manufacturer_website = gofakeit.URL()

		uuids = []string{uuid}
		names = []string{name}
		categories = []model.Category{model.Category(filterCategory)}
		manufacturer_countries = []string{manufacturer_country}


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

		partsFilter = model.PatrsFilter{
			UUIDs: uuids,
			Names: names,
			Categories: categories,
			Manufacture_contries: manufacturer_countries,
			Tags: tags,
		}

		listPartsResponse = model.ListPartsResponse{Parts: []model.Part{part}}

		listPartsRequest = model.ListPartsRequest{
			Filter: partsFilter,
		}
	)

	s.inventoryRepository.On("ListParts", s.ctx, listPartsRequest).Return(listPartsResponse, nil)

	res, err := s.service.ListParts(s.ctx, listPartsRequest)

	s.Require().NoError(err)
	//s.Require().Equal(model.ListPartsResponse{}, res)
	s.Require().Equal(listPartsResponse, res)
}

func (s *ServiceSuite) TestListPartsOnlyOnePart() {
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

		uuids = []string{uuid}
		names = []string{name}
		categories = []model.Category{model.Category(category)}
		manufacturer_countries = []string{manufacturer_country}


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

		part1 = model.Part{
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
		part2 = model.Part{
			UUID: gofakeit.UUID(),
			Name: gofakeit.Word(),
			Description: description,
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

		partsFilter = model.PatrsFilter{
			UUIDs: uuids,
			Names: names,
			Categories: categories,
			Manufacture_contries: manufacturer_countries,
			Tags: tags,
		}

		listPartsResponse = model.ListPartsResponse{Parts: []model.Part{part1, part2}}

		listPartsRequest = model.ListPartsRequest{
			Filter: partsFilter,
		}
	)

	s.inventoryRepository.On("ListParts", s.ctx, listPartsRequest).Return(listPartsResponse, nil)

	res, err := s.service.ListParts(s.ctx, listPartsRequest)

	s.Require().NoError(err)
	//s.Require().Equal(model.ListPartsResponse{Parts: []model.Part{part1}}, res)
	s.Require().Equal(listPartsResponse, res)
}


