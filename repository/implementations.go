package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

const (
	InsertEstateQuery     = `INSERT INTO estates (id, width, length) VALUES ($1, $2, $3) RETURNING id`
	GetEstateQuery        = `SELECT id, created_at, updated_at, deleted_at, width, length FROM estates`
	InsertEstateTreeQuery = `INSERT INTO estate_trees (id, estate_id, x, y, height) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	GetEstateTreeQuery    = `SELECT id, estate_id, created_at, updated_at, deleted_at, x, y, height FROM estate_trees`
	EstateTreeCountQuery  = `SELECT COUNT(1) FROM estate_trees`
	EstateTreeStatsQuery  = `SELECT MAX(height) as max, MIN(height) as min, PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY COALESCE(height, 0)) AS median FROM estate_trees`
)

func (r *Repository) CreateEstate(ctx context.Context, data *Estate) error {
	_, err := r.Db.ExecContext(
		ctx,
		InsertEstateQuery,
		data.ID,
		data.Width,
		data.Length,
	)
	return err
}

func (r *Repository) FindEstate(ctx context.Context, filter *FilterEstate) (Estate, error) {
	finalQuery, paramValue := r.setFilterEstate(GetEstateQuery, filter)
	var estate Estate
	err := r.Db.QueryRowContext(ctx, finalQuery, paramValue...).Scan(&estate.ID, &estate.CreatedAt, &estate.UpdatedAt, &estate.DeletedAt, &estate.Width, &estate.Length)
	if err != nil {
		return Estate{}, err
	}

	return estate, nil
}

func (r *Repository) setFilterEstate(baseQuery string, filter *FilterEstate) (string, []interface{}) {
	var (
		where      []string
		paramValue []interface{}
	)
	if filter.ID != "" {
		where = append(where, "id = $"+strconv.Itoa(len(paramValue)+1))
		paramValue = append(paramValue, filter.ID)
	}

	where = append(where, "deleted_at IS NULL")
	clauseWhere := strings.Join(where, " AND ")
	if clauseWhere != "" {
		baseQuery += " WHERE " + clauseWhere
	}

	return baseQuery, paramValue
}

func (r *Repository) CreateEstateTree(ctx context.Context, data *EstateTree) error {
	_, err := r.Db.ExecContext(
		ctx,
		InsertEstateTreeQuery,
		data.ID,
		data.EstateID,
		data.X,
		data.Y,
		data.Height,
	)
	return err
}

func (r *Repository) FindAllMapEstateTree(ctx context.Context, filter *FilterEstateTree) (map[CoordinatePoint]EstateTree, error) {
	result := make(map[CoordinatePoint]EstateTree)

	finalQuery, paramValue := r.setFilterEstateTree(GetEstateTreeQuery, filter)
	if filter.OrderBy != "" {
		finalQuery += " ORDER BY " + filter.OrderBy + " " + filter.Sort
	}
	if filter.Limit > 0 || !filter.ShowAll {
		limit := "LIMIT $" + strconv.Itoa(len(paramValue)+1) + " OFFSET $" + strconv.Itoa(len(paramValue)+2)
		finalQuery = fmt.Sprintf("%s %s", finalQuery, limit)
		paramValue = append(paramValue, filter.Limit, filter.CalculateOffset())
	}
	rows, err := r.Db.QueryContext(ctx, finalQuery, paramValue...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var estateTree EstateTree
		if err = rows.Scan(&estateTree.ID, &estateTree.EstateID, &estateTree.CreatedAt,
			&estateTree.UpdatedAt, &estateTree.DeletedAt,
			&estateTree.X, &estateTree.Y,
			&estateTree.Height); err != nil {
			return nil, err
		}
		result[CoordinatePoint{X: estateTree.X, Y: estateTree.Y}] = estateTree
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Repository) FindEstateTree(ctx context.Context, filter *FilterEstateTree) (EstateTree, error) {
	finalQuery, paramValue := r.setFilterEstateTree(GetEstateTreeQuery, filter)
	var estateTree EstateTree
	err := r.Db.QueryRowContext(ctx, finalQuery, paramValue...).
		Scan(&estateTree.ID, &estateTree.EstateID, &estateTree.CreatedAt,
			&estateTree.UpdatedAt, &estateTree.DeletedAt,
			&estateTree.X, &estateTree.Y,
			&estateTree.Height)
	if err != nil {
		return EstateTree{}, err
	}

	return estateTree, nil
}

func (r *Repository) CountEstateTree(ctx context.Context, filter *FilterEstateTree) int {
	finalQuery, paramValue := r.setFilterEstateTree(EstateTreeCountQuery, filter)

	var count int64
	err := r.Db.QueryRowContext(ctx, finalQuery, paramValue...).Scan(&count)
	if err != nil {
		return 0
	}

	return int(count)
}

func (r *Repository) setFilterEstateTree(baseQuery string, filter *FilterEstateTree) (string, []interface{}) {
	var (
		where      []string
		paramValue []interface{}
	)
	if filter.ID != "" {
		where = append(where, "id = $"+strconv.Itoa(len(paramValue)+1))
		paramValue = append(paramValue, filter.ID)
	}
	if filter.EstateID != "" {
		where = append(where, "estate_id = $"+strconv.Itoa(len(paramValue)+1))
		paramValue = append(paramValue, filter.EstateID)
	}
	if filter.X != 0 {
		where = append(where, "x = $"+strconv.Itoa(len(paramValue)+1))
		paramValue = append(paramValue, filter.X)
	}
	if filter.Y != 0 {
		where = append(where, "y = $"+strconv.Itoa(len(paramValue)+1))
		paramValue = append(paramValue, filter.Y)
	}

	where = append(where, "deleted_at IS NULL")
	clauseWhere := strings.Join(where, " AND ")
	if clauseWhere != "" {
		baseQuery += " WHERE " + clauseWhere
	}

	return baseQuery, paramValue
}

func (r *Repository) GetEstateTreeStats(ctx context.Context, filter *FilterEstateTree) (EstateTreeStats, error) {
	finalQuery, paramValue := r.setFilterEstateTree(EstateTreeStatsQuery, filter)
	var estateTreeStats EstateTreeStats
	err := r.Db.QueryRowContext(ctx, finalQuery, paramValue...).
		Scan(&estateTreeStats.Max, &estateTreeStats.Min, &estateTreeStats.Median)
	if err != nil {
		return EstateTreeStats{}, err
	}

	return estateTreeStats, nil
}
