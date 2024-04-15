package repo

import (
	"cars-service/internal/model"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repoImpl struct {
	*pgxpool.Pool
}

func (r *repoImpl) GetCarById(ctx context.Context, id uint64) (model.Car, error) {
	row := r.QueryRow(ctx, getCarByIdQuery, id)
	var car model.Car
	if err := row.Scan(
		&car.Id,
		&car.RegNum,
		&car.Mark,
		&car.Model,
		&car.Year,
		&car.Owner.Name,
		&car.Owner.Surname,
		&car.Owner.Patronymic,
	); errors.Is(err, pgx.ErrNoRows) {
		return model.Car{}, model.ErrCarNotFound
	} else if err != nil {
		return model.Car{}, errors.Join(model.ErrDatabaseError, err)
	} else {
		return car, nil
	}
}

func (r *repoImpl) GetCars(ctx context.Context, filter model.Filter) ([]model.Car, error) {
	rows, err := r.Query(ctx, getCarsQuery,
		filter.RegNum,
		filter.ByRegNum,
		filter.Mark,
		filter.ByMark,
		filter.Model,
		filter.ByModel,
		filter.Year,
		filter.ByYear,
		filter.OwnerName,
		filter.ByOwnerName,
		filter.OwnerSurname,
		filter.ByOwnerSurname,
		filter.OwnerPatronymic,
		filter.ByOwnerPatronymic,
		filter.Limit,
		filter.Offset,
	)
	if err != nil {
		return []model.Car{}, errors.Join(model.ErrDatabaseError, err)
	}
	defer rows.Close()

	cars := make([]model.Car, 0, filter.Limit)
	for rows.Next() {
		var car model.Car
		_ = rows.Scan(
			&car.Id,
			&car.RegNum,
			&car.Mark,
			&car.Model,
			&car.Year,
			&car.Owner.Name,
			&car.Owner.Surname,
			&car.Owner.Patronymic,
		)
		cars = append(cars, car)
	}
	return cars, nil
}

func (r *repoImpl) AddCar(ctx context.Context, car model.Car) (model.Car, error) {
	ownerId, err := r.getOwnerId(ctx, car.Owner)
	if err != nil {
		return model.Car{}, err
	}
	markId, err := r.getMarkId(ctx, car.Mark)
	if err != nil {
		return model.Car{}, err
	}
	modelId, err := r.getModelId(ctx, car.Model, markId)
	if err != nil {
		return model.Car{}, err
	}

	// TODO добавить обработку ошибки добавления дубликата номера
	if err = r.QueryRow(ctx, insertCarQuery,
		car.RegNum,
		modelId,
		car.Year,
		ownerId,
	).Scan(&car.Id); err != nil {
		return model.Car{}, errors.Join(model.ErrDatabaseError, err)
	}
	return car, nil
}

func (r *repoImpl) UpdateCar(ctx context.Context, id uint64, car model.Car) (model.Car, error) {
	ownerId, err := r.getOwnerId(ctx, car.Owner)
	if err != nil {
		return model.Car{}, err
	}
	markId, err := r.getMarkId(ctx, car.Mark)
	if err != nil {
		return model.Car{}, err
	}
	modelId, err := r.getModelId(ctx, car.Model, markId)
	if err != nil {
		return model.Car{}, err
	}

	// TODO сюда тоже добавить
	if e, err := r.Exec(ctx, updateCarQuery,
		id,
		car.RegNum,
		modelId,
		car.Year,
		ownerId,
	); err != nil {
		return model.Car{}, errors.Join(model.ErrDatabaseError, err)
	} else if e.RowsAffected() == 0 {
		return model.Car{}, model.ErrCarNotFound
	} else {
		car.Id = id
		return car, nil
	}
}

func (r *repoImpl) DeleteCar(ctx context.Context, id uint64) error {
	if e, err := r.Exec(ctx, deleteCarQuery, id); err != nil {
		return errors.Join(model.ErrDatabaseError, err)
	} else if e.RowsAffected() == 0 {
		return model.ErrCarNotFound
	}
	return nil
}

func (r *repoImpl) getOwnerId(ctx context.Context, owner model.Owner) (uint64, error) {
	var id uint64
	if err := r.QueryRow(ctx, getOwnerIdQuery,
		owner.Name,
		owner.Surname,
		owner.Patronymic,
	).Scan(&id); err != nil {
		return 0, errors.Join(model.ErrDatabaseError, err)
	}
	return id, nil
}

func (r *repoImpl) getMarkId(ctx context.Context, mark string) (uint64, error) {
	var id uint64
	if err := r.QueryRow(ctx, getMarkIdQuery,
		mark,
	).Scan(&id); err != nil {
		return 0, errors.Join(model.ErrDatabaseError, err)
	}
	return id, nil
}

func (r *repoImpl) getModelId(ctx context.Context, model string, markId uint64) (uint64, error) {
	var id uint64
	if err := r.QueryRow(ctx, getModelIdQuery,
		model,
		markId,
	).Scan(&id); err != nil {
		return 0, errors.Join(model.ErrDatabaseError, err)
	}
	return id, nil
}

const (
	getCarByIdQuery = `
		SELECT "cars"."id", "reg_num", "marks"."name", "models"."name", "year", "owners"."name", "owners"."surname", "owners"."patronymic"
		FROM "cars"
			INNER JOIN "models" ON "cars"."model_id" = "models"."id"
			INNER JOIN "marks" ON "models"."mark_id" = "marks"."id"
			INNER JOIN "owners" ON "cars"."owner_id" = "owners"."id"
		WHERE "cars"."id" = $1;`

	getCarsQuery = `
		SELECT "cars"."id", "reg_num", "marks"."name", "models"."name", "year", "owners"."name", "owners"."surname", "owners"."patronymic"
		FROM "cars"
			INNER JOIN "models" ON "cars"."model_id" = "models"."id"
			INNER JOIN "marks" ON "models"."mark_id" = "marks"."id"
			INNER JOIN "owners" ON "cars"."owner_id" = "owners"."id"
		WHERE ("reg_num" = $1 OR (NOT $2))
			AND ("marks"."name" = $3 OR (NOT $4))
			AND ("models"."name" = $5 OR (NOT $6))
			AND ("year" = $7 OR (NOT $8))
			AND ("owners"."name" = $9 OR (NOT $10))
			AND ("owners"."surname" = $11 OR (NOT $12))
			AND ("owners"."patronymic" = $13 OR (NOT $14))
		LIMIT $15 OFFSET $16;`

	insertCarQuery = `
		INSERT INTO "cars" (reg_num, model_id, year, owner_id) 
		VALUES ($1, $2, $3, $4)
		RETURNING "id";`

	updateCarQuery = `
		UPDATE "cars"
		SET "reg_num" = $2,
		    "model_id" = $3,
		    "year" = $4,
		    "owner_id" = $5
		WHERE "id" = $1;`

	deleteCarQuery = `
		DELETE FROM "cars"
		WHERE "id" = $1;`

	getOwnerIdQuery = `
		WITH inserted_or_existing AS (
			INSERT INTO "owners" ("name", "surname", "patronymic")
			SELECT $1, $2, $3
			WHERE NOT EXISTS (
				SELECT 1 FROM "owners"
				WHERE "name" = $1 AND "surname" = $2 AND "patronymic" = $3
			)
			RETURNING "id"
		)
		SELECT "id" FROM inserted_or_existing
		UNION ALL
		SELECT "id" FROM "owners"
		WHERE "name" = $1 AND "surname" = $2 AND "patronymic" = $3;`

	getMarkIdQuery = `
		WITH inserted_or_existing AS (
			INSERT INTO "marks" ("name")
			SELECT $1
			WHERE NOT EXISTS (
				SELECT 1 FROM "marks"
				WHERE "name" = $1
			)
			RETURNING "id"
		)
		SELECT "id" FROM inserted_or_existing
		UNION ALL
		SELECT "id" FROM "marks"
		WHERE "name" = $1;`

	getModelIdQuery = `
		WITH inserted_or_existing AS (
			INSERT INTO "models" ("name", "mark_id")
			SELECT $1, $2
			WHERE NOT EXISTS (
				SELECT 1 FROM "models"
				WHERE "name" = $1 AND "mark_id" = $2
			)
			RETURNING "id"
		)
		SELECT "id" FROM inserted_or_existing
		UNION ALL
		SELECT "id" FROM "models"
		WHERE "name" = $1 AND "mark_id" = $2;`
)
