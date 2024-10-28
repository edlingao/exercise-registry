package routes

import (
	"exercise-app/controllers"

	"github.com/labstack/echo/v4"
)

type ExerciseRoutes string

const EXERCISE ExerciseRoutes = "/exercise"

func (ex ExerciseRoutes) Build(e *echo.Echo) {
	exercise_routes := e.Group(string(ex))
	exercise_controller := controllers.ExerciseController{}
  // GET
	exercise_routes.GET("", exercise_controller.Show)
  exercise_routes.GET("/edit/:name/:id", exercise_controller.ShowEditForm)
  exercise_routes.GET("/show/:name/:id", exercise_controller.ShowColumn)

  // POST
  exercise_routes.POST("/new", exercise_controller.Create)

  // PATCH
  exercise_routes.PATCH("/update/:name/:id", exercise_controller.UpdateColumn)
  exercise_routes.PATCH("/update/duration/:id", exercise_controller.UpdateDuration)

  // DELETE
  exercise_routes.DELETE("/delete/:id", exercise_controller.Delete)
}
