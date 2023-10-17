package handlers

import (
	"net/http"

	"estiam_golang_api_course_finalproject/services"
	"estiam_golang_api_course_finalproject/types"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{
		userService: service,
	}
}

/*REMOVE THIS ENDPOINT*/
// func (h *UserHandler) Get(ctx echo.Context) error {
// 	id := ctx.Param("id")

// 	user, err := h.userService.GetUser(ctx.Request().Context(), id)
// 	if err != nil {
// 		ctx.NoContent(http.StatusInternalServerError)
// 		return nil
// 	}

// 	if user == nil {
// 		ctx.JSON(http.StatusNotFound, map[string]interface{}{"message": "user not found"})
// 		return nil
// 	}

// 	ctx.JSON(http.StatusOK, user)
// 	return nil
// }

/*"POST /users"  endpoint*/
func (h *UserHandler) Post(ctx echo.Context) error {
	//Get the json
	newUser := new(types.User)
	if err := ctx.Bind(newUser); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}
	//method for create newuser
	if err := h.userService.CreateNewUser(ctx.Request().Context(), newUser); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create user"})
	}

	return ctx.JSON(http.StatusCreated, map[string]string{"message": "User created successfully"})
}

/*"POST /login"  endpoint*/
func (h *UserHandler) Login(ctx echo.Context) error {
	// Récupérez les données JSON du corps de la demande
	credentials := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	if err := ctx.Bind(&credentials); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}

	// Appelez la méthode de service pour gérer la création du JWT
	jwtToken, err := h.userService.GenerateJWT(ctx.Request().Context(), credentials.Username, credentials.Password)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"message": "Authentication failed"})
	}

	// Retournez le JWT dans la réponse
	return ctx.JSON(http.StatusOK, map[string]string{"token": jwtToken})
}

/*Get All Users*/
func (h *UserHandler) GetAllUsers(ctx echo.Context) error {
	users, err := h.userService.GetAllUsers(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to get users"})
	}
	return ctx.JSON(http.StatusOK, users)
}
