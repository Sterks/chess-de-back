package delivery

import (
	"chess-backend/internal/config"
	"chess-backend/internal/service"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) Init(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
		corsMiddleware,
	)

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	router.GET("/allGetStep", h.getAllStep)
	router.MaxMultipartMemory = 10 << 20
	router.POST("/upload", h.saveUpload)
}

func (h *Handler) getAllStep(c *gin.Context) {
	arraySteps, err := h.services.StepsService.GetSteps()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "bad request")
	}
	c.JSON(http.StatusOK, arraySteps)
}

func (h *Handler) saveUpload(c *gin.Context) {
	name := c.Request.FormValue("name")

	_, file, err := c.Request.FormFile("file")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
		})
		return
	}

	h.services.ProcessingService.AddNamesBook(name)
	if err := c.SaveUploadedFile(file, "upload/"+file.Filename); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file exists",
		})
		return
	}

	steps, err := h.services.ProcessingService.ReadProcessing(file.Filename)
	sentence := strings.Join(steps, " ")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
	} else {
		c.String(http.StatusOK, fmt.Sprintf("Основные ходы: %v\n", sentence))
	}

	c.String(http.StatusOK, fmt.Sprintf(" ---- %s uploaded!", file.Filename))
}
