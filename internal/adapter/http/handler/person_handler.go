package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	personusecase "github.com/mashurimansur/goCMS/internal/usecase/person"
)

// PersonHandler exposes HTTP endpoints for person-related use cases.
type PersonHandler struct {
	personUseCase personusecase.UseCase
}

func NewPersonHandler(personUseCase personusecase.UseCase) *PersonHandler {
	return &PersonHandler{
		personUseCase: personUseCase,
	}
}

// Register wires the handler routes under the provided router group.
func (h *PersonHandler) Register(router *gin.RouterGroup) {
	router.GET("/person", h.GetDefaultPerson)
}

// GetDefaultPerson responds with the default person payload.
func (h *PersonHandler) GetDefaultPerson(c *gin.Context) {
	person, err := h.personUseCase.GetDefaultPerson(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, person)
}
