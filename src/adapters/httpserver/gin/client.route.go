package httpserver_gin

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"rinha.backend.2024/src/adapters/config"
	"rinha.backend.2024/src/adapters/data/repositories"
	"rinha.backend.2024/src/internal/domains"
	"rinha.backend.2024/src/internal/domains/dtos/requests"
	"rinha.backend.2024/src/internal/usercases"
	"rinha.backend.2024/src/pkg/database"
	"rinha.backend.2024/src/pkg/validation"
)

func NewClientGroup(app *gin.Engine, dataBase *database.PostgresDatabase, rinhaConfig *config.RinhaBackendConfig) {
	g := app.Group("/clientes")
	{
		g.GET("/:id/extrato", handleGetLatestStatementBalance(dataBase, rinhaConfig))
		g.POST("/:id/transacoes", handleRegisterTransaction(dataBase, rinhaConfig))
	}
}

func handleGetLatestStatementBalance(dataBase *database.PostgresDatabase, rinhaConfig *config.RinhaBackendConfig) func(c *gin.Context) {
	return func(c *gin.Context) {
		var payload requests.GetLatestStatementBalanceRequestDto
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.Status(http.StatusUnprocessableEntity)
			return
		}

		payload.ClientID = domains.ID(id)

		clientRepo := repositories.NewClientRepository(dataBase.DB)
		transactionRepo := repositories.NewTransactionRepository(dataBase.DB)
		statementUserCase := usercases.NewStatementUsercase(clientRepo, transactionRepo, time.Duration(rinhaConfig.TimeZone))

		result, err := statementUserCase.Execute(c.Request.Context(), payload)
		buildResponse(c, result, err)
	}
}

func handleRegisterTransaction(dataBase *database.PostgresDatabase, rinhaConfig *config.RinhaBackendConfig) func(c *gin.Context) {
	return func(c *gin.Context) {
		var payload requests.RegisterTransactionRequestDto
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.Status(http.StatusUnprocessableEntity)
			return
		}

		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.Status(http.StatusUnprocessableEntity)
			return
		}

		payload.ClientID = domains.ID(id)

		clientRepo := repositories.NewClientRepository(dataBase.DB)
		transactionRepo := repositories.NewTransactionRepository(dataBase.DB)
		uow := database.NewUnitOfWork(dataBase.DB)
		transactionUserCase := usercases.NewTransactionUsercase(clientRepo, transactionRepo, uow, func() time.Time { return time.Now().UTC() })

		result, err := transactionUserCase.Execute(c.Request.Context(), payload)
		buildResponse(c, result, err)
	}
}

func buildResponse[TResult any](c *gin.Context, result TResult, err error) {
	if err != nil {
		if errors.Is(err, validation.ErrValidation) {
			c.Status(http.StatusUnprocessableEntity)
			return
		}

		if errors.Is(err, validation.ErrNotFound) {
			c.Status(http.StatusNotFound)
			return
		}

		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, result)
}
