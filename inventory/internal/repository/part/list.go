package part

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"

	repoModel "github.com/danilfaer/golang/inventory/internal/repository/model"
	partUtils "github.com/danilfaer/golang/inventory/internal/repository/utils"
)	

func (r *repository) ListParts(ctx context.Context, filter *repoModel.PartsFilter) ([]*repoModel.Part, error) {
	q := partUtils.PartsFilterToMongoFilter(filter)
	if q == nil {
		q = bson.M{}
	}
	cur, err := r.coll.Find(ctx, q)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []*repoModel.Part
	for cur.Next(ctx) {
		var doc repoModel.Part
		if err := cur.Decode(&doc); err != nil {
			return nil, err
		}
		// копия на куче: иначе все указатели укажут на одну и ту же переменную цикла
		row := doc
		out = append(out, &row)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	if out == nil {
		out = []*repoModel.Part{}
	}
	return out, nil
}
