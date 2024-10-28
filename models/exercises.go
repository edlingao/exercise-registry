package models

import (
	"exercise-app/libs"
)

type ExerciseModel string

const EXERCISE ExerciseModel = "EXCERSICE_MODEL"

func (c ExerciseModel) GetByID(exercise_id string) (Exercise, error) {
	var exercise Exercise
	libs.ConnectDB()
	defer libs.DisconnectDB()

	err := libs.DB.QueryRowx(`
    SELECT id, user_id, date, feeling, hours, minutes, calories
    FROM exercises
    WHERE id = ?
  `, exercise_id).StructScan(&exercise)

	if err != nil {
		return exercise, err
	}

	return exercise, nil
}

func (e ExerciseModel) GetByUserID(user_id string) ([]Exercise, error) {
	var exercises []Exercise
	libs.ConnectDB()
	defer libs.DisconnectDB()

	err := libs.DB.Select(&exercises, `
    SELECT id, user_id, date, feeling, hours, minutes, calories
    FROM exercises
    WHERE user_id = ?
  `, user_id)

	if err != nil {
		return exercises, err
	}

	return exercises, nil
}

func (e ExerciseModel) Create(exercise Exercise) (Exercise, error) {
	libs.ConnectDB()
	defer libs.DisconnectDB()

	row, err := libs.DB.NamedExec(`
    INSERT INTO exercises (user_id, date, feeling, hours, minutes, calories)
    VALUES (:user_id, :date, :feeling, :hours, :minutes, :calories)
  `, exercise)

	if err != nil {
		return Exercise{}, err
	}

	id, err := row.LastInsertId()

	if err != nil {
		return Exercise{}, err
	}
	var returendRow Exercise

	error := libs.DB.QueryRowx(`
    SELECT id, user_id, date, feeling, hours, minutes, calories
    FROM exercises
    WHERE id = ?
  `, id).StructScan(&returendRow)

	return returendRow, error
}

func (e ExerciseModel) Delete(id int) (int, error) {
	libs.ConnectDB()
	defer libs.DisconnectDB()

	_, err := libs.DB.Exec(`
    DELETE FROM exercises
    WHERE id = ?
  `, id)

	return id, err
}

func (e ExerciseModel) Update(field string, value string, id string) (bool, error) {
  libs.ConnectDB()
  defer libs.DisconnectDB()

  _, err := libs.DB.Exec(`
    UPDATE exercises
    SET `+ field +` = ?
    WHERE id = ?
  `, value, id)

  if err != nil {
    return false, err
  }

  return true, nil
}
