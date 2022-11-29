package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// add new "GET /hello" route to the app router (echo)
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/emi_mobile/emprendimientos/:id", //variable
			Handler: func(c echo.Context) error {
				emprendimiento_record, err := app.Dao().FindFirstRecordByData("emprendimientos", "id", c.PathParam("id"))
				if err != nil {
					return apis.NewNotFoundError("The articlee does not exist.", err)
				}

				promotor_fk := emprendimiento_record.Get("id_promotor_fk").(string)
				prioridad_fk := emprendimiento_record.Get("id_prioridad_fk").(string)
				status_fk := emprendimiento_record.Get("id_status_sync_fk").(string)
				fase_fk := emprendimiento_record.Get("id_fase_emp_fk").(string)
				emprededor_fk := emprendimiento_record.Get("id_emprendedor_fk").(string)
				cat_proyecto_fk := emprendimiento_record.Get("id_nombre_proyecto_fk").(string)
				jornada_fk := emprendimiento_record.Get("id").(string)

				// var tareajornada = map[string]interface{}{}
				// tareajornada = make(map[string]interface{})

				// var jornadalist = map[int]interface{}{}
				// jornadalist = make(map[int]interface{})

				jornada, err := app.Dao().FindRecordsByExpr("jornadas", dbx.HashExp{"id_emprendimiento_fk": jornada_fk})
				if err != nil {
					return apis.NewNotFoundError("The articlee does not exist.", err)
				}

				tarea, err := app.Dao().FindRecordsByExpr("tareas", dbx.HashExp{"id": jornada[0].Get("id_tarea_fk")})
				if err != nil {
					return apis.NewNotFoundError("The articlee does not exist.", err)
				}
				// for i := 0; i < len(jornada); i++ {
				// 	jornada[i].Id = jornadalist[i].(map[string]interface{})["id"].(string)

				// 	/*
				// 		id_jornada, err := app.Dao().FindFirstRecordByData("jornadas", "id_emprendimiento_fk", jornadalist[i].(string))
				// 		if err != nil {
				// 			return apis.NewNotFoundError("The articlee does not exist.", err)
				// 		}
				// 		id_jornada_fk := id_jornada.Get("id_tarea_fk").(string)
				// 	*/

				// 	tareas, err := app.Dao().FindRecordsByExpr("tareas", dbx.HashExp{"id_tarea_fk": jornada[i].Id})
				// 	if err != nil {
				// 		return apis.NewNotFoundError("The articlee does not exist.", err)
				// 	}
				// 	tareajornada = make(map[string]interface{})

				// 	//Aqui lo que tenemos que hacer que agarre el id de la tarea de la tabla jornada y usarlo.

				// }

				promotor, err := app.Dao().FindRecordsByExpr("emi_users", dbx.HashExp{"id": promotor_fk})
				if err != nil {
					return apis.NewNotFoundError("The articlee does not exist.", err)
				}
				cat_proyecto, err := app.Dao().FindRecordsByExpr("cat_proyecto", dbx.HashExp{"id": cat_proyecto_fk})
				if err != nil {
					return apis.NewNotFoundError("The articlee does not exist.", err)
				}

				prioridad, err := app.Dao().FindRecordsByExpr("prioridades_emp", dbx.HashExp{"id": prioridad_fk})
				if err != nil {
					return apis.NewNotFoundError("The articlee does not exist.", err)
				}

				fase, err := app.Dao().FindRecordsByExpr("fases_emp", dbx.HashExp{"id": fase_fk})
				if err != nil {
					return apis.NewNotFoundError("The articlee does not exist.", err)
				}

				status, err := app.Dao().FindRecordsByExpr("status_sync", dbx.HashExp{"id": status_fk})
				if err != nil {
					return apis.NewNotFoundError("The articlee does not exist.", err)
				}

				emprendedor, err := app.Dao().FindFirstRecordByData("emprendedores", "id", emprededor_fk)
				if err != nil {
					return apis.NewNotFoundError("The articlee does not exist.", err)
				}

				comunidad_fk := emprendedor.Get("id_comunidad_fk").(string)

				comunidad, err := app.Dao().FindRecordsByExpr("comunidades", dbx.HashExp{"id": comunidad_fk})
				if err != nil {
					return apis.NewNotFoundError("The articlee does not exist.", err)
				}

				infoEmprendimiento := map[string]interface{}{
					"emprendimiento": emprendimiento_record,
					"fase":           fase,
					"proyecto":       cat_proyecto,
				}

				infoEmprendedor := map[string]interface{}{
					"emprendedor": emprendedor,
					"comunidad":   comunidad,
				}
				infoJornada := map[string]interface{}{
					"jornadas": jornada,
					"tareas":   tarea,
				}

				InfoTotal := map[string]interface{}{
					"info_emprendimiento": infoEmprendimiento,
					"promotor":            promotor,
					"prioridad":           prioridad,
					"status":              status,
					"info_emprendedor":    infoEmprendedor,
					"info_jornada":        infoJornada,
				}

				return c.JSON(http.StatusOK, InfoTotal)
				//Todo ok
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
			},
		})
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
