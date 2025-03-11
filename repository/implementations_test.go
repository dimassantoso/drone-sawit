package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRepository_CreateEstate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := &Repository{Db: db}

		id := uuid.NewString()
		mock.ExpectExec("INSERT INTO estates").WithArgs(id, 100, 200).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.CreateEstate(context.Background(), &Estate{
			BaseModel: BaseModel{
				ID: id,
			},
			Width:  100,
			Length: 200,
		})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := &Repository{Db: db}

		mock.ExpectExec("INSERT INTO estates").
			WithArgs(sqlmock.AnyArg(), 100, 200).
			WillReturnError(assert.AnError)

		err = repo.CreateEstate(context.Background(), &Estate{
			BaseModel: BaseModel{
				ID: uuid.NewString(),
			},
			Width:  100,
			Length: 200,
		})

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRepository_FindEstate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := &Repository{Db: db}

		id := uuid.NewString()
		expectedEstate := Estate{
			BaseModel: BaseModel{
				ID:        id,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Width:  100,
			Length: 200,
		}

		mock.ExpectQuery("SELECT .* FROM estates").WithArgs(id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "width", "length"}).
				AddRow(expectedEstate.ID, expectedEstate.CreatedAt, expectedEstate.UpdatedAt, nil, expectedEstate.Width, expectedEstate.Length))

		result, err := repo.FindEstate(context.Background(), &FilterEstate{ID: id})
		assert.NoError(t, err)
		assert.Equal(t, expectedEstate, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed: No rows found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := &Repository{Db: db}

		id := uuid.NewString()
		mock.ExpectQuery("SELECT .* FROM estates").WithArgs(id).WillReturnError(sql.ErrNoRows)

		result, err := repo.FindEstate(context.Background(), &FilterEstate{ID: id})
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.Equal(t, Estate{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed: Query error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := &Repository{Db: db}

		id := "#######"
		mock.ExpectQuery("SELECT .* FROM estates").
			WithArgs(id).
			WillReturnError(errors.New("database error"))

		result, err := repo.FindEstate(context.Background(), &FilterEstate{ID: id})
		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		assert.Equal(t, Estate{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRepository_CreateEstateTree(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := &Repository{Db: db}

		id := uuid.NewString()
		estateID := uuid.NewString()
		mock.ExpectExec("INSERT INTO estate_trees").WithArgs(id, estateID, 1, 2, 30).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.CreateEstateTree(context.Background(), &EstateTree{
			BaseModel: BaseModel{
				ID: id,
			},
			EstateID: estateID,
			X:        1,
			Y:        2,
			Height:   30,
		})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := &Repository{Db: db}

		mock.ExpectExec("INSERT INTO estate_trees").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), 1, 2, 30).
			WillReturnError(assert.AnError)

		err = repo.CreateEstateTree(context.Background(), &EstateTree{
			BaseModel: BaseModel{
				ID: uuid.NewString(),
			},
			EstateID: uuid.NewString(),
			X:        1,
			Y:        2,
			Height:   30,
		})
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRepository_FindAllMapEstateTree(t *testing.T) {
	t.Run("Success : filter estate_id", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := &Repository{Db: db}

		estateID := uuid.NewString()
		filter := &FilterEstateTree{
			Filter: Filter{
				Page:    1,
				Limit:   10,
				OrderBy: "id",
				Sort:    "ASC",
			},
			EstateID: estateID,
		}

		mock.ExpectQuery("SELECT .* FROM estate_trees WHERE estate_id = \\$1 AND deleted_at IS NULL ORDER BY id ASC LIMIT \\$2 OFFSET \\$3").
			WithArgs(estateID, 10, 0).
			WillReturnRows(sqlmock.NewRows([]string{"id", "estate_id", "created_at", "updated_at", "deleted_at", "x", "y", "height"}).
				AddRow(uuid.NewString(), estateID, time.Now(), time.Now(), nil, 1, 2, 30).
				AddRow(uuid.NewString(), estateID, time.Now(), time.Now(), nil, 3, 4, 25))

		result, err := repo.FindAllMapEstateTree(context.Background(), filter)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, 30, result[CoordinatePoint{X: 1, Y: 2}].Height)
		assert.Equal(t, 25, result[CoordinatePoint{X: 3, Y: 4}].Height)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed : Query", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := &Repository{Db: db}

		estateID := uuid.NewString()
		filter := &FilterEstateTree{
			Filter: Filter{
				Page:    1,
				Limit:   10,
				OrderBy: "id",
				Sort:    "ASC",
			},
			EstateID: estateID,
		}

		mock.ExpectQuery("SELECT .* FROM estate_trees WHERE estate_id = \\$1 AND deleted_at IS NULL ORDER BY id ASC LIMIT \\$2 OFFSET \\$3").
			WithArgs(estateID, 10, 0).WillReturnError(errors.New("query failed"))

		result, err := repo.FindAllMapEstateTree(context.Background(), filter)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "query failed", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed : Row Scan", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := &Repository{Db: db}

		estateID := uuid.NewString()
		filter := &FilterEstateTree{
			Filter: Filter{
				Page:    1,
				Limit:   10,
				OrderBy: "id",
				Sort:    "ASC",
			},
			EstateID: estateID,
		}

		mock.ExpectQuery("SELECT .* FROM estate_trees WHERE estate_id = \\$1 AND deleted_at IS NULL ORDER BY id ASC LIMIT \\$2 OFFSET \\$3").
			WithArgs(estateID, 10, 0).
			WillReturnRows(sqlmock.NewRows([]string{"id", "estate_id", "created_at", "updated_at", "deleted_at", "x", "y", "height"}).
				AddRow("", estateID, time.Now(), time.Now(), nil, 1, "xxxxxxxx", 30))

		result, err := repo.FindAllMapEstateTree(context.Background(), filter)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "sql: Scan error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed : error row", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := &Repository{Db: db}

		estateID := uuid.NewString()
		filter := &FilterEstateTree{
			Filter: Filter{
				Page:  1,
				Limit: 10,
			},
			EstateID: estateID,
		}

		rows := sqlmock.NewRows([]string{
			"id", "estate_id", "created_at", "updated_at", "deleted_at", "x", "y", "height",
		}).
			AddRow(uuid.NewString(), estateID, time.Now(), time.Now(), nil, 1, 2, 30).
			AddRow(uuid.NewString(), estateID, time.Now(), time.Now(), nil, 3, 4, 21).
			RowError(1, fmt.Errorf("iteration error"))

		mock.ExpectQuery("SELECT .* FROM estate_trees WHERE estate_id = \\$1 AND deleted_at IS NULL LIMIT \\$2 OFFSET \\$3").
			WithArgs(estateID, 10, 0).
			WillReturnRows(rows)

		result, err := repo.FindAllMapEstateTree(context.Background(), filter)

		assert.Error(t, err)
		assert.EqualError(t, err, "iteration error")
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRepository_FindEstateTree(t *testing.T) {
	t.Run("Success : filter by id", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := &Repository{Db: db}

		id := uuid.NewString()
		estateID := uuid.NewString()
		expectedEstateTree := EstateTree{
			BaseModel: BaseModel{
				ID:        id,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			EstateID: estateID,
			X:        1,
			Y:        2,
			Height:   20,
		}

		mock.ExpectQuery("SELECT .* FROM estate_trees").WithArgs(id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "estate_id", "created_at", "updated_at", "deleted_at", "x", "y", "height"}).
				AddRow(expectedEstateTree.ID, expectedEstateTree.EstateID, expectedEstateTree.CreatedAt, expectedEstateTree.UpdatedAt, nil, expectedEstateTree.X, expectedEstateTree.Y, expectedEstateTree.Height))

		result, err := repo.FindEstateTree(context.Background(), &FilterEstateTree{ID: id})
		assert.NoError(t, err)
		assert.Equal(t, expectedEstateTree, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success : filter by x and y", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := &Repository{Db: db}

		id := uuid.NewString()
		estateID := uuid.NewString()
		expectedEstateTree := EstateTree{
			BaseModel: BaseModel{
				ID:        id,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			EstateID: estateID,
			X:        3,
			Y:        4,
			Height:   10,
		}

		mock.ExpectQuery("SELECT .* FROM estate_trees WHERE estate_id = \\$1 AND x = \\$2 AND y = \\$3 AND deleted_at IS NULL").
			WithArgs(estateID, 3, 4).
			WillReturnRows(sqlmock.NewRows([]string{"id", "estate_id", "created_at", "updated_at", "deleted_at", "x", "y", "height"}).
				AddRow(expectedEstateTree.ID, expectedEstateTree.EstateID, expectedEstateTree.CreatedAt, expectedEstateTree.UpdatedAt, nil, expectedEstateTree.X, expectedEstateTree.Y, expectedEstateTree.Height))

		result, err := repo.FindEstateTree(context.Background(), &FilterEstateTree{X: 3, Y: 4, EstateID: estateID})
		assert.NoError(t, err)
		assert.Equal(t, expectedEstateTree, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed: No rows found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := &Repository{Db: db}

		id := uuid.NewString()
		mock.ExpectQuery("SELECT .* FROM estate_trees").WithArgs(id).WillReturnError(sql.ErrNoRows)

		result, err := repo.FindEstateTree(context.Background(), &FilterEstateTree{ID: id})
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.Equal(t, EstateTree{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed: Query error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := &Repository{Db: db}

		id := "#######"
		mock.ExpectQuery("SELECT .* FROM estate_trees").WithArgs(id).
			WillReturnError(errors.New("database error"))

		result, err := repo.FindEstateTree(context.Background(), &FilterEstateTree{ID: id})
		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		assert.Equal(t, EstateTree{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRepository_CountEstateTree(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := &Repository{Db: db}

		filter := &FilterEstateTree{
			EstateID: uuid.NewString(),
		}

		expectedQuery := "SELECT COUNT\\(1\\) FROM estate_trees WHERE estate_id = \\$1 AND deleted_at IS NULL"
		mock.ExpectQuery(expectedQuery).WithArgs(filter.EstateID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

		count := repo.CountEstateTree(context.Background(), filter)
		assert.Equal(t, 5, count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ErrorDuringQuery", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := &Repository{Db: db}
		filter := &FilterEstateTree{
			EstateID: uuid.NewString(),
		}

		expectedQuery := "SELECT COUNT\\(1\\) FROM estate_trees WHERE estate_id = \\$1 AND deleted_at IS NULL"
		mock.ExpectQuery(expectedQuery).WithArgs(filter.EstateID).WillReturnError(fmt.Errorf("query error"))

		count := repo.CountEstateTree(context.Background(), filter)
		assert.Equal(t, 0, count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMockRepositoryInterface_GetEstateTreeStats(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := &Repository{Db: db}
		filter := &FilterEstateTree{
			EstateID: uuid.NewString(),
		}

		expectedQuery := `SELECT MAX\(height\) as max, MIN\(height\) as min, PERCENTILE_CONT\(0.5\) WITHIN GROUP \(ORDER BY COALESCE\(height, 0\)\) AS median FROM estate_trees`
		mock.ExpectQuery(expectedQuery).
			WithArgs(filter.EstateID).WillReturnRows(sqlmock.NewRows([]string{"max", "min", "median"}).
			AddRow(30, 3, 15))

		stats, err := repo.GetEstateTreeStats(context.Background(), filter)
		assert.NoError(t, err)
		assert.Equal(t, 30, stats.Max)
		assert.Equal(t, 3, stats.Min)
		assert.Equal(t, float32(15), stats.Median)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ErrorDuringQuery", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := &Repository{Db: db}

		filter := &FilterEstateTree{
			Filter: Filter{
				Page:  1,
				Limit: 10,
			},
			EstateID: uuid.NewString(),
		}

		expectedQuery := `SELECT MAX\(height\) as max, MIN\(height\) as min, PERCENTILE_CONT\(0.5\) WITHIN GROUP \(ORDER BY COALESCE\(height, 0\)\) AS median FROM estate_trees`
		mock.ExpectQuery(expectedQuery).WithArgs(filter.EstateID).WillReturnError(fmt.Errorf("query error"))

		stats, err := repo.GetEstateTreeStats(context.Background(), filter)
		assert.Error(t, err)
		assert.Equal(t, EstateTreeStats{}, stats)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
