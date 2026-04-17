package partUtils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"

	repoModel "github.com/danilfaer/golang/inventory/internal/repository/model"
)

func TestPartsFilterToMongoFilter_NilOrEmpty(t *testing.T) {
	assert.Nil(t, PartsFilterToMongoFilter(nil))
	assert.Nil(t, PartsFilterToMongoFilter(&repoModel.PartsFilter{}))
}

func TestPartsFilterToMongoFilter_UUIDs(t *testing.T) {
	f := &repoModel.PartsFilter{
		Uuids: []string{"a", "b"},
	}
	got := PartsFilterToMongoFilter(f)
	orsAny, ok := got["$or"]
	assert.True(t, ok)
	ors, ok := orsAny.(bson.A)
	if !ok {
		ors2, ok2 := orsAny.([]bson.M)
		assert.True(t, ok2)
		assert.Len(t, ors2, 2)
		return
	}
	assert.Len(t, ors, 2)
}

func TestPartsFilterToMongoFilter_CategoriesORIn(t *testing.T) {
	f := &repoModel.PartsFilter{
		Categories: []repoModel.Category{repoModel.CategoryEngine, repoModel.CategoryFuel},
	}
	got := PartsFilterToMongoFilter(f)
	in, ok := got["category"].(bson.M)
	assert.True(t, ok)
	assert.Equal(t, []string{"ENGINE", "FUEL"}, in["$in"])
}

func TestPartsFilterToMongoFilter_ANDBetweenGroups(t *testing.T) {
	f := &repoModel.PartsFilter{
		Uuids:      []string{"u1"},
		Categories: []repoModel.Category{repoModel.CategoryEngine},
	}
	got := PartsFilterToMongoFilter(f)
	and, ok := got["$and"]
	assert.True(t, ok)
	arr, ok := and.([]bson.M)
	assert.True(t, ok)
	assert.Len(t, arr, 2)
}
