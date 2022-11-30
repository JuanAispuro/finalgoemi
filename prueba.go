package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
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
					return apis.NewNotFoundError(" No hay emprendimientos con ese ID.", err)
				}

				promotor_fk := emprendimiento_record.Get("id_promotor_fk").(string)
				prioridad_fk := emprendimiento_record.Get("id_prioridad_fk").(string)
				status_fk := emprendimiento_record.Get("id_status_sync_fk").(string)
				fase_fk := emprendimiento_record.Get("id_fase_emp_fk").(string)
				emprededor_fk := emprendimiento_record.Get("id_emprendedor_fk").(string)
				cat_proyecto_fk := emprendimiento_record.Get("id_nombre_proyecto_fk").(string)
				jornada_fk := emprendimiento_record.Get("id").(string)

				// Sacamos los datos de la jornada con su id_emprendimientos
				jornada, err := app.Dao().FindRecordsByExpr("jornadas", dbx.HashExp{"id_emprendimiento_fk": jornada_fk})
				if err != nil {
					return apis.NewNotFoundError(" No hay jornadas con ese ID.", err)
				}
				// Creamos un arreglo de tipo []*model.Record{}
				tareajornada := []*models.Record{}

				// Creamos un for para que nos almacene los diferentes ID de jornadas
				for i := 0; i < len(jornada); i++ {

					// Creamos un arreglo donde seleccione de la tabla tareas los id que agarramos con el for.
					tareas, err := app.Dao().FindRecordsByExpr("tareas", dbx.HashExp{"id": jornada[i].Get("id_tarea_fk")})
					if err != nil {
						return apis.NewNotFoundError(" No hay tareas con este ID.", err)
					}
					tareajornada = append(tareas, tareajornada...)
				}

				// Sacamos la información de los promotores usando el id.
				promotor, err := app.Dao().FindRecordsByExpr("emi_users", dbx.HashExp{"id": promotor_fk})
				if err != nil {
					return apis.NewNotFoundError(" No hay usuarios con este ID.", err)
				}
				//Sacas la información de la tabla de catalogo_proyecto
				cat_proyecto, err := app.Dao().FindRecordsByExpr("cat_proyecto", dbx.HashExp{"id": cat_proyecto_fk})
				if err != nil {
					return apis.NewNotFoundError(" No hay ningun catalogo para el proyecto.", err)
				}

				// Scas la prioridad de la tabla prioridades
				prioridad, err := app.Dao().FindRecordsByExpr("prioridades_emp", dbx.HashExp{"id": prioridad_fk})
				if err != nil {
					return apis.NewNotFoundError(" No hay prioridades con este ID.", err)
				}

				// Sacas la fase de la tabla de fase_emprendedores.
				fase, err := app.Dao().FindRecordsByExpr("fases_emp", dbx.HashExp{"id": fase_fk})
				if err != nil {
					return apis.NewNotFoundError(" No hay fases con este ID.", err)
				}

				// Sacas el status de la tabla status_sync
				status, err := app.Dao().FindRecordsByExpr("status_sync", dbx.HashExp{"id": status_fk})
				if err != nil {
					return apis.NewNotFoundError(" No hay ningun status de sincronización con este ID.", err)
				}

				// Sacas la información de la tabla emprendedor
				emprendedor, err := app.Dao().FindFirstRecordByData("emprendedores", "id", emprededor_fk)
				if err != nil {
					return apis.NewNotFoundError(" No hay ningun emprededor con este ID.", err)
				}

				// Sacas de la tabla de los emprededores el id de la comunidad.
				comunidad_fk := emprendedor.Get("id_comunidad_fk").(string)

				// Sacas la información de la tabla comunidades con su id.
				comunidad, err := app.Dao().FindRecordsByExpr("comunidades", dbx.HashExp{"id": comunidad_fk})
				if err != nil {
					return apis.NewNotFoundError(" No hay ninguna comunidad con este ID.", err)
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
					"jornadas, ": jornada,
					"tareas":     tareajornada,
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
