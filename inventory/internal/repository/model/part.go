package model

import "time"

// информация о заказе
type Part struct {
	//Уникальный идентификатор заказа
	UUID string

	// Название детали
	Name string

	// Описание детали
	Description string

	//Цена за единицу
	Price float64

	// Количество на складе
	Stock_quantity int64

	//Категория
	Category Category

	// Размер детали
	Dimentions Dimentions

	// Информация о производителе
	Manufacturer Manufacturer

	// Теги для быстрого поика (слайс строк)
	Tags []string

	// Гибкие метаданные
	Metadata map[string]Value

	// Дата создания
	created_at *time.Time

	// Дата обновления
	updated_at *time.Time
}

// размеры детали
type Dimentions struct {
	// Длина в см
	Lenght float64

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
