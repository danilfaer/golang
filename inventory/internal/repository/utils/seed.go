package partUtils

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	repoModel "github.com/danilfaer/golang/inventory/internal/repository/model"
)

// sampleParts возвращает эталонный набор деталей (как бывший createSampleParts).
func sampleParts() []repoModel.Part {
	now := time.Now()
	return []repoModel.Part{
		{
			UUID: "550e8400-e29b-41d4-a716-446655440001", Name: "Ионный двигатель X-2000",
			Description: "Высокоэффективный ионный двигатель для межпланетных полетов",
			Price: 150000.0, StockQuantity: 200, Category: repoModel.CategoryEngine,
			Dimensions: &repoModel.Dimensions{Length: 120, Width: 80, Height: 60, Weight: 250},
			Manufacturer: &repoModel.Manufacturer{Name: "КосмоТех", Country: "Россия", Website: "https://cosmotech.ru"},
			Tags: []string{"ионный", "двигатель", "межпланетный", "высокоэффективный"},
			CreatedAt: now.Add(-30 * 24 * time.Hour), UpdatedAt: now,
		},
		{
			UUID: "550e8400-e29b-41d4-a716-446655440002", Name: "Плазменный двигатель P-500",
			Description: "Мощный плазменный двигатель для тяжелых грузов",
			Price: 200000.0, StockQuantity: 100, Category: repoModel.CategoryEngine,
			Dimensions: &repoModel.Dimensions{Length: 150, Width: 100, Height: 80, Weight: 400},
			Manufacturer: &repoModel.Manufacturer{Name: "StarTech Industries", Country: "США", Website: "https://startech.com"},
			Tags: []string{"плазменный", "двигатель", "тяжелый", "грузовой"},
			CreatedAt: now.Add(-45 * 24 * time.Hour), UpdatedAt: now,
		},
		{
			UUID: "550e8400-e29b-41d4-a716-446655440003", Name: "Криогенное топливо H2-O2",
			Description: "Высокоэнергетическое криогенное топливо для ракетных двигателей",
			Price: 50000.0, StockQuantity: 999, Category: repoModel.CategoryFuel,
			Dimensions: &repoModel.Dimensions{Length: 200, Width: 100, Height: 100, Weight: 1500},
			Manufacturer: &repoModel.Manufacturer{Name: "КриоТопливо", Country: "Россия", Website: "https://cryofuel.ru"},
			Tags: []string{"криогенное", "топливо", "водород", "кислород"},
			CreatedAt: now.Add(-60 * 24 * time.Hour), UpdatedAt: now,
		},
		{
			UUID: "550e8400-e29b-41d4-a716-446655440004", Name: "Ядерное топливо U-235",
			Description: "Обогащенный уран для ядерных реакторов",
			Price: 300000.0, StockQuantity: 2, Category: repoModel.CategoryFuel,
			Dimensions: &repoModel.Dimensions{Length: 50, Width: 30, Height: 30, Weight: 100},
			Manufacturer: &repoModel.Manufacturer{Name: "АтомЭнерго", Country: "Россия", Website: "https://atomenergo.ru"},
			Tags: []string{"ядерное", "топливо", "уран", "реактор"},
			CreatedAt: now.Add(-90 * 24 * time.Hour), UpdatedAt: now,
		},
		{
			UUID: "550e8400-e29b-41d4-a716-446655440005", Name: "Кварцевое окно QW-100",
			Description: "Прозрачное кварцевое окно для космических кораблей",
			Price: 25000.0, StockQuantity: 15, Category: repoModel.CategoryPorthole,
			Dimensions: &repoModel.Dimensions{Length: 100, Width: 100, Height: 10, Weight: 50},
			Manufacturer: &repoModel.Manufacturer{Name: "КварцТех", Country: "Россия", Website: "https://quartztech.ru"},
			Tags: []string{"кварцевое", "окно", "прозрачное", "космическое"},
			CreatedAt: now.Add(-20 * 24 * time.Hour), UpdatedAt: now,
		},
		{
			UUID: "550e8400-e29b-41d4-a716-446655440006", Name: "Бронированное окно BW-200",
			Description: "Защищенное окно с многослойным покрытием",
			Price: 40000.0, StockQuantity: 8, Category: repoModel.CategoryPorthole,
			Dimensions: &repoModel.Dimensions{Length: 120, Width: 120, Height: 15, Weight: 80},
			Manufacturer: &repoModel.Manufacturer{Name: "ArmorGlass", Country: "Германия", Website: "https://armorglass.de"},
			Tags: []string{"бронированное", "окно", "защищенное", "многослойное"},
			CreatedAt: now.Add(-15 * 24 * time.Hour), UpdatedAt: now,
		},
		{
			UUID: "550e8400-e29b-41d4-a716-446655440007", Name: "Солнечная панель SP-500",
			Description: "Высокоэффективная солнечная панель для космических станций",
			Price: 75000.0, StockQuantity: 12, Category: repoModel.CategoryWing,
			Dimensions: &repoModel.Dimensions{Length: 500, Width: 200, Height: 5, Weight: 300},
			Manufacturer: &repoModel.Manufacturer{Name: "СолнТех", Country: "Россия", Website: "https://solntech.ru"},
			Tags: []string{"солнечная", "панель", "энергия", "космическая"},
			CreatedAt: now.Add(-25 * 24 * time.Hour), UpdatedAt: now,
		},
		{
			UUID: "550e8400-e29b-41d4-a716-446655440008", Name: "Аэродинамическое крыло AW-300",
			Description: "Легкое аэродинамическое крыло для атмосферных полетов",
			Price: 60000.0, StockQuantity: 10, Category: repoModel.CategoryWing,
			Dimensions: &repoModel.Dimensions{Length: 300, Width: 150, Height: 20, Weight: 200},
			Manufacturer: &repoModel.Manufacturer{Name: "AeroDynamics", Country: "Франция", Website: "https://aerodynamics.fr"},
			Tags: []string{"аэродинамическое", "крыло", "легкое", "атмосферное"},
			CreatedAt: now.Add(-35 * 24 * time.Hour), UpdatedAt: now,
		},
	}
}

// SeedPartsIfEmpty вставляет эталонные детали, если коллекция пуста.
func SeedPartsIfEmpty(ctx context.Context, coll *mongo.Collection) error {
	n, err := coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	return insertSampleParts(ctx, coll)
}

// ReplaceAllWithSampleParts удаляет все документы и вставляет эталонный набор (для интеграционных тестов).
func ReplaceAllWithSampleParts(ctx context.Context, coll *mongo.Collection) error {
	if _, err := coll.DeleteMany(ctx, bson.M{}); err != nil {
		return err
	}
	return insertSampleParts(ctx, coll)
}

func insertSampleParts(ctx context.Context, coll *mongo.Collection) error {
	docs := sampleParts()
	raw := make([]any, len(docs))
	for i := range docs {
		raw[i] = &docs[i]
	}
	_, err := coll.InsertMany(ctx, raw)
	return err
}
