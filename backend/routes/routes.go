package routes

import (
	"github.com/gin-gonic/gin"
	"portal/controllers"
	"portal/middleware"
)

func RegisterRoutes(r *gin.Engine) {
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	receptionist := r.Group("/receptionist")
	receptionist.Use(middleware.AuthMiddleware("receptionist"))
	{
		receptionist.POST("/patients", controllers.CreatePatient)
		receptionist.GET("/patients", controllers.GetPatients)
		receptionist.GET("/patients/:id", controllers.FetchPatient)
		receptionist.PUT("/patients/:id", controllers.UpdatePatient)
		receptionist.DELETE("/patients/:id", controllers.DeletePatient)
	}

	doctor := r.Group("/doctor")
	doctor.Use(middleware.AuthMiddleware("doctor"))
	{
		doctor.GET("/patients", controllers.GetPatients)
		doctor.GET("/patients/:id", controllers.FetchPatient)
		doctor.PUT("/patients/:id", controllers.UpdatePatient)
	}
}
