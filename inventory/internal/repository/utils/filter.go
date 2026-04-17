package partUtils

import (
	"go.mongodb.org/mongo-driver/v2/bson"

	repoModel "github.com/danilfaer/golang/inventory/internal/repository/model"
)

// PartsFilterToMongoFilter строит фильтр MongoDB с той же семантикой, что у старого in-memory list.go:
// пустой/nil фильтр — вернуть nil (означает «без ограничений» для вызывающего);
// между непустыми группами условий — AND, внутри группы — OR.
func PartsFilterToMongoFilter(filter *repoModel.PartsFilter) bson.M {
	if filter == nil || isEmptyFilter(filter) {
		return nil
	}
	var ands []bson.M

	if len(filter.Uuids) > 0 {
		ors := make([]bson.M, 0, len(filter.Uuids))
		for _, u := range filter.Uuids {
			ors = append(ors, bson.M{"_id": u})
		}
		ands = append(ands, bson.M{"$or": ors})
	}
	if len(filter.Names) > 0 {
		ors := make([]bson.M, 0, len(filter.Names))
		for _, n := range filter.Names {
			ors = append(ors, bson.M{"name": n})
		}
		ands = append(ands, bson.M{"$or": ors})
	}
	if len(filter.Categories) > 0 {
		cats := make([]string, 0, len(filter.Categories))
		for _, c := range filter.Categories {
			cats = append(cats, string(c))
		}
		ands = append(ands, bson.M{"category": bson.M{"$in": cats}})
	}
	if len(filter.ManufacturerCountries) > 0 {
		ors := make([]bson.M, 0, len(filter.ManufacturerCountries))
		for _, country := range filter.ManufacturerCountries {
			ors = append(ors, bson.M{"manufacturer.country": country})
		}
		ands = append(ands, bson.M{"$or": ors})
	}
	if len(filter.Tags) > 0 {
		ors := make([]bson.M, 0, len(filter.Tags))
		for _, t := range filter.Tags {
			ors = append(ors, bson.M{"tags": t})
		}
		ands = append(ands, bson.M{"$or": ors})
	}

	if len(ands) == 1 {
		return ands[0]
	}
	return bson.M{"$and": ands}
}

func isEmptyFilter(filter *repoModel.PartsFilter) bool {
	return len(filter.Uuids) == 0 &&
		len(filter.Names) == 0 &&
		len(filter.Categories) == 0 &&
		len(filter.ManufacturerCountries) == 0 &&
		len(filter.Tags) == 0
}
