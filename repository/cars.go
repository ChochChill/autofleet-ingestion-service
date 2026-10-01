package repository

import (
	"carsAPI/model"
	"database/sql"
)

type CarRepository struct {
	db *sql.DB
}

func NewCarRepository(db *sql.DB) *CarRepository {
	return &CarRepository{
		db: db,
	}
}

func (repo *CarRepository) GetAll() ([]*model.Car, error) {
	rows, err := repo.db.Query("SELECT id, make, model, year FROM cars")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cars := make([]*model.Car, 0)
	for rows.Next() {
		car := &model.Car{}
		if err := rows.Scan(
			&car.ID,
			&car.Make,
			&car.Model,
			&car.Year,
		); err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return cars, nil
}

func (repo *CarRepository) GetByID(id int) (*model.Car, bool, error) {
	row := repo.db.QueryRow(
		"SELECT id, make, model, year FROM cars WHERE id = ?",
		id,
	)
	car := &model.Car{}

	err := row.Scan(&car.ID, &car.Make, &car.Model, &car.Year)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return car, true, nil
}

func (repo *CarRepository) Update(id int, car *model.Car) (*model.Car, bool, error) {
	row := repo.db.QueryRow(`
		UPDATE cars 
		SET make = ?, model = ?, year = ? 
		WHERE id = ?
		RETURNING id, make, model, year 
		`,
		car.Make,
		car.Model,
		car.Year,
		id,
	)

	updatedCar := &model.Car{}

	err := row.Scan(
		&updatedCar.ID,
		&updatedCar.Make,
		&updatedCar.Model,
		&updatedCar.Year,
	)

	if err == sql.ErrNoRows {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, err
	}

	return updatedCar, true, nil
}

func (repo *CarRepository) Create(car *model.Car) (*model.Car, error) {
	result, err := repo.db.Exec(
		"INSERT INTO cars (make,model,year) VALUES (?, ?, ?)",
		car.Make,
		car.Model,
		car.Year,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	car.ID = int(id)
	return car, nil
}

func (repo *CarRepository) Delete(id int) (bool, error) {
	result, err := repo.db.Exec(
		"DELETE FROM cars WHERE id = ?",
		id,
	)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}
