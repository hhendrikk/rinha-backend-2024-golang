package httpserver

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"rinha.backend.2024/src/adapters/config"
	"rinha.backend.2024/src/adapters/data/repositories"
	"rinha.backend.2024/src/internal/domains/dtos/requests"
	"rinha.backend.2024/src/internal/usercases"
	"rinha.backend.2024/src/pkg/database"
	"rinha.backend.2024/src/pkg/validation"
)

func NewClientGroup(app *echo.Echo, dataBase *database.PostgresDatabase, rinhaConfig *config.RinhaBackendConfig) {
	g := app.Group("/clientes")

	g.GET("/:id/extrato", handleGetLatestStatementBalance(dataBase, rinhaConfig))
	g.POST("/:id/transacoes", handleRegisterTransaction(dataBase, rinhaConfig))
}

func handleGetLatestStatementBalance(dataBase *database.PostgresDatabase, rinhaConfig *config.RinhaBackendConfig) func(c echo.Context) error {
	return func(c echo.Context) error {
		var payload requests.GetLatestStatementBalanceRequestDto
		if err := c.Bind(&payload); err != nil {
			return c.NoContent(http.StatusUnprocessableEntity)
		}

		clientRepo := repositories.NewClientRepository(dataBase.DB)
		transactionRepo := repositories.NewTransactionRepository(dataBase.DB)
		statementUserCase := usercases.NewStatementUsercase(clientRepo, transactionRepo, time.Duration(rinhaConfig.TimeZone))

		result, err := statementUserCase.Execute(c.Request().Context(), payload)
		return buildResponse(c, result, err)
	}
}

func handleRegisterTransaction(dataBase *database.PostgresDatabase, rinhaConfig *config.RinhaBackendConfig) func(c echo.Context) error {
	return func(c echo.Context) error {
		var payload requests.RegisterTransactionRequestDto
		if err := c.Bind(&payload); err != nil {
			return c.NoContent(http.StatusUnprocessableEntity)
		}

		clientRepo := repositories.NewClientRepository(dataBase.DB)
		transactionRepo := repositories.NewTransactionRepository(dataBase.DB)
		uow := database.NewUnitOfWork(dataBase.DB)
		transactionUserCase := usercases.NewTransactionUsercase(clientRepo, transactionRepo, uow, func() time.Time { return time.Now().UTC() })

		result, err := transactionUserCase.Execute(c.Request().Context(), payload)
		return buildResponse(c, result, err)
	}
}

func buildResponse[TResult any](c echo.Context, result TResult, err error) error {
	if err != nil {
		if errors.Is(err, validation.ErrValidation) {
			return c.JSON(http.StatusUnprocessableEntity, err.Error())
		}

		if errors.Is(err, validation.ErrNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, result)
}
