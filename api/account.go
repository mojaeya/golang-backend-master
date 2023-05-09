package api

import (
	"database/sql"
	db "golang-backend-master/db/sqlc"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	createAccountRequest struct {
		Owner    string `json:"owner" validate:"required"`
		Currency string `json:"currency" validate:"required,oneof=USD EUR"`
	}
)

func (server *Server) createAccount(ctx echo.Context) error {
	var req createAccountRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	if err := ctx.Validate(req); err != nil {
		return err
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx.Request().Context(), arg)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return ctx.JSON(http.StatusOK, account)
}

// 기본적으로 echo.Context의 유효성 검증기는 POST나 PUT과 같은 요청에서 요청 본문의 데이터를 바인딩할 때, 필드 유효성 검증을 수행합니다. Param이나 Query의 경우, 데이터를 바인딩할 때 해당하는 필드 유효성 검증을 수행하지 않습니다.

type getAccountRequest struct {
	ID int64 `param:"id" validate:"required,min=1"`
}

func (server *Server) getAccount(ctx echo.Context) error {
	var req getAccountRequest

	// 먼저, echo.Context.Bind() 함수를 호출하여 매개 변수를 구조체 필드에 할당한 다음, 유효성 검증을 수행합니다
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	// go-playground/validator 라이브러리는 구조체 필드에 대해서만 유효성 검증이 가능하므로, param 태그 값을 무시하게 됩니다. 따라서, ID 필드에 대한 유효성 검증을 수행하려면, 요청 매개변수 배열에서 값을 뽑아서 직접 검사해야 합니다.
	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	// err := (&echo.DefaultBinder{}).BindPathParams(ctx, &req)
	// if err != nil {
	// 	return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	// }

	// id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	// if err != nil {
	// 	return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	// }

	// if req.ID <= 0 {
	// 	return ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("Invalid account ID")))
	// }

	account, err := server.store.GetAccount(ctx.Request().Context(), req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, errorResponse(err))
		}
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	PageID   int32 `query:"page_id" validate:"required,min=1"`
	PageSize int32 `query:"page_size" validate:"required,min=5,max=10"`
}

func (server *Server) listAccount(ctx echo.Context) error {
	var req listAccountRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	account, err := server.store.ListAccounts(ctx.Request().Context(), arg)

	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, errorResponse(err))
		}
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return ctx.JSON(http.StatusOK, account)
}
