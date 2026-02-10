package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// информация о заказе
type Part struct {
	ID primitive.ObjectID		`bson:"_id,omitempty"`
	//Уникальный идентификатор заказа
	UUID string `bson:"uuid"`

	// Название детали
	Name string  `bson:"name"`

	// Описание детали
	Description string  `bson:"description"`

	//Цена за единицу
	Price float64  `bson:"price"`

	// Количество на складе
	Stock_quantity int64  `bson:"stock_quantity"`

	//Категория
	Category int32		`bson:"category"`

	// Размер детали
	Dimensions Dimensions  `bson:"dimensions"`

	// Информация о производителе
	Manufacturer Manufacturer		`bson:"manufacturer"`

	// Теги для быстрого поика (слайс строк)
	Tags []string		`bson:"tags"`

	// Гибкие метаданные
	Metadata map[string]interface{}		`bson:"metadata,omitempty"`

	// Дата создания
	CreatedAt *time.Time		`bson:"created_at,omitempty"`

	// Дата обновления
	UpdatedAt *time.Time		`bson:"updated_at,omitempty"`
}

// размеры детали
type Dimensions struct {
	// Длина в см
	Length float64

	//Ширина в см
	Width float64

	// Высота в см
	Height float64

	// Вес в кг
	Weight float64
}

// Производитель
type Manufacturer struct {
	//Название
	Name string

	//Старна производителя
	Country string

	//Сайт производителя
	Website string
}

// enum структура для перечисления категорий
type Category int32

const (
	CATEGORY_UNKNOWN_UNSPECIFIED Category = 0
	// Двигатель
	CATEGORY_ENGINE Category = 1
	// Топливо
	CATEGORY_FUEL Category = 2
	// Иллюминатор
	CATEGORY_PORTHOLE Category = 3
	// Крыло
	CATEGORY_WING Category = 4
)

// Я вообще хз как сделать тут one of так что пока что через интерфейс
type Value interface {
	isValue()
}

// Строка
type StringValue struct {
	Value string
}

func (StringValue) isValue() {}

// Числа int
type Int64Value struct {
	Value int64
}

func (Int64Value) isValue() {}

// числа с плавающей точкой (float32)
type DoubleValue struct {
	Value float64
}

func (DoubleValue) isValue() {}

// Логические значения
type BoolValue struct {
	Value bool
}

func (BoolValue) isValue() {}

// структура фильтра по деталям, все поля опциональный (пустой слайс)
type PatrsFilter struct {
	// Список UUID'ов. Пусто - не фильтруем по UUID
	UUIDs []string

	// Список имен. Пусто - не фильтруем по имени
	Names []string

	// Список категорий. Пусто - не фильтруем по категории
	Categories []Category

	// Список стран производителей. Пусто - не фильтруем по стране
	Manufacture_contries []string

	// Список тегов. Пусто - не фильтруем по тегам
	Tags []string
}

// GetPartRequest запрос на получение информации о детали по ее UUID
type GetPartRequest struct {
	//Идентификатор детали
	UUID string
}

// GetPartResponse ответ с инфорацией о детали
type GetPartResponse struct {
	// Вся информация о детали
	Parts Part
}

// ListPartsRequest запрос на получение списка деталей с возможностью фильтрации
type ListPartsRequest struct {
	// Фильтр по деталям (все поля опциональны)
	Filter PatrsFilter
}

// ListPartsResponse ответ со списком найденных деталей
type ListPartsResponse struct {
	// Список найденных деталей
	Parts []Part
}
