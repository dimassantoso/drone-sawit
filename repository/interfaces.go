package repository

import (
	"context"
)

type RepositoryInterface interface {
	CreateEstate(ctx context.Context, data *Estate) error
	FindEstate(ctx context.Context, filter *FilterEstate) (Estate, error)
	CreateEstateTree(ctx context.Context, data *EstateTree) error
	FindAllMapEstateTree(ctx context.Context, filter *FilterEstateTree) (map[CoordinatePoint]EstateTree, error)
	FindEstateTree(ctx context.Context, filter *FilterEstateTree) (EstateTree, error)
	CountEstateTree(ctx context.Context, filter *FilterEstateTree) int
	GetEstateTreeStats(ctx context.Context, filter *FilterEstateTree) (EstateTreeStats, error)
}
