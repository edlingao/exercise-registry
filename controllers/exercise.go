package controllers

import (
	"errors"
	"exercise-app/components"
	"exercise-app/models"
	"exercise-app/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type ExerciseController struct{}
type ExerciseParamType[ParamType string | int] struct {}

// @Router /exercise [GET]
func (e ExerciseController) Show(c echo.Context) error {
	user_req := c.Get("user")
	user, ok := user_req.(models.User)

	if !ok {
		return c.Redirect(http.StatusTemporaryRedirect, "/auth/signin")
	}

	excercises, err := models.EXERCISE.GetByUserID(user.ID)

	if err != nil {
		return c.HTML(
			http.StatusBadGateway,
			err.Error(),
		)
	}
	return utils.Render(
		c,
		components.ExerciseTable(
			components.ExerciseTableViewModel{
				Exercises: excercises,
			},
		),
		http.StatusOK,
	)
}

// @Router /exercise/new [POST]
func (e ExerciseController) Create(c echo.Context) error {
	user_req := c.Get("user")
	user, ok := user_req.(models.User)

	if !ok {
		return c.Redirect(http.StatusTemporaryRedirect, "/auth/signin")
	}

	date := time.Now()
	feeling := c.FormValue("feeling")
  hours := c.FormValue("hours")
  minutes := c.FormValue("minutes")
	calories := c.FormValue("calories")
	calories_int, _ := strconv.ParseInt(calories, 10, 64)
  hours_int, _ := strconv.ParseInt(hours, 10, 64)
  minutes_int, _ := strconv.ParseInt(minutes, 10, 64)

	exercise := models.Exercise{
		User_ID:  user.ID,
		Date:     date,
		Feeling:  feeling,
    Hours: int(hours_int),
    Minutes: int(minutes_int),
		Calories: int(calories_int),
	}

	_, err := models.EXERCISE.Create(exercise)

	if err != nil {
		return utils.Render(
			c,
			components.ExerciseForm(
				components.ExerciseFormViewModel{
					Error: map[string]string{
						"general": err.Error(),
					}},
			),
			http.StatusOK,
		)
	}
	_, error := models.EXERCISE.GetByUserID(user.ID)

	if error != nil {
		return utils.Render(
			c,
			components.ExerciseForm(
				components.ExerciseFormViewModel{
					Error: map[string]string{
						"general": err.Error(),
					}},
			),
			http.StatusOK,
		)
	}

  c.Response().Header().Set("HX-Trigger", "{ \"exercises:loaded\": \"\" }")

  return utils.Render(
    c,
    components.ExerciseForm(
      components.ExerciseFormViewModel{},
    ),
    http.StatusOK,
  )
}

// @Router /exercise/delete/{id} [DELETE]
func (e ExerciseController) Delete(c echo.Context) error {
  id := c.Param("id")
  id_int, _ := strconv.ParseInt(id, 10, 64)

  _, err := models.EXERCISE.Delete(int(id_int))

  if err != nil {
    return c.HTML(
      http.StatusBadGateway,
      err.Error(),
    )
  }
  
  c.Response().Header().Set("HX-Trigger", "{ \"exercises:loaded\": \"\" }")

  return c.JSON(http.StatusOK, map[string]string{
    "message": "Exercise deleted",
  })
}

// @Router /exercise/edit/{name}/{id} [GET]
func (e ExerciseController) ShowEditForm(c echo.Context) error {
  name := c.Param("name")
  id := c.Param("id")
  type_ := c.QueryParam("type")

  exercise, err := models.EXERCISE.GetByID(id)
  value, ok := exercise.GetParam(name)

  if err != nil || !ok {
    return utils.Render(
      c,
      components.EditTableColumn(components.ExerciseTableColumnViewModel{
        Name: name,
        Value: "",
        Error: errors.New("Exercise not found"),
      }),
      http.StatusNotFound,
    )
  }
  
  return utils.Render(
    c,
    components.EditTableColumn(components.ExerciseTableColumnViewModel{
      ID: id,
      Name: name,
      Value: value,
      Hours: exercise.GetHours(),
      Minutes: exercise.GetMinutes(),
      Type: type_,
    }),
    http.StatusOK,
  )
}

// @Router /exercise/show/{name}/{id} [GET]
func (e ExerciseController) ShowColumn(c echo.Context) error {
  name := c.Param("name")
  id := c.Param("id")
  type_ := c.QueryParam("type")

  exercise, err := models.EXERCISE.GetByID(id)
  value, ok := exercise.GetParam(name)

  if !ok || err != nil {
    return utils.Render(
      c,
      components.TableColumn(components.ExerciseTableColumnViewModel{
        Name: name,
        Value: "",
        Error: errors.New("Exercise not found"),
      }),
      http.StatusNotFound,
    )
  }
  
  return utils.Render(
    c,
    components.TableColumn(components.ExerciseTableColumnViewModel{
      ID: id,
      Name: name,
      Value: string(value),
      Type: type_,
    }),
    http.StatusOK,
  )
}

// @Router /exercise/update/{name}/{id} [PATCH]
func (e ExerciseController) UpdateColumn(c echo.Context) error {
  name := c.Param("name")
  id := c.Param("id")
  value := c.FormValue(name)
  type_ := c.FormValue("type")

  saved, err := models.EXERCISE.Update(name, value, id)

  if err != nil || !saved {
    return utils.Render(
      c,
      components.TableColumn(components.ExerciseTableColumnViewModel{
        Name: name,
        Value: "",
        Error: errors.New("Exercise not found"),
      }),
      http.StatusNotFound,
    )
  }

  exercise, _ := models.EXERCISE.GetByID(id)
  value, _ = exercise.GetParam(name)

  return utils.Render(
    c,
    components.TableColumn(components.ExerciseTableColumnViewModel{
      ID: id,
      Name: name,
      Value: value,
      Type: type_,
    }),
    http.StatusOK,
  )
}

func (e ExerciseController) UpdateDuration(c echo.Context) error {
  id := c.Param("id")
  hours := c.FormValue("hours")
  minutes := c.FormValue("minutes")

  saved_hours, err := models.EXERCISE.Update("hours", hours, id)
  saved_minutes, err := models.EXERCISE.Update("minutes", minutes, id)

  if err != nil || !saved_hours || !saved_minutes {
    return c.HTML(
      http.StatusBadGateway,
      err.Error(),
    )
  }

  exercise, _ := models.EXERCISE.GetByID(id)

  return utils.Render(
    c,
    components.TableColumn(components.ExerciseTableColumnViewModel{
      ID: id,
      Name: "duration",
      Value: exercise.GetDuration(),
    }),
    http.StatusOK,
  )
}
