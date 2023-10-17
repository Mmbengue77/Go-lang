package services

/*
TODO: Add service that ENCRYPTS PASSWORD and STORES IN DB

Tips
- Use bcrypt package!
- Save bcrypt secret on .env and load it in App configuration!
- Inject app configuration (bcrypt secret) into here (user service)
*/

import (
	"context"
	"fmt"
	"os"
	"time"

	"estiam_golang_api_course_finalproject/repos"
	"estiam_golang_api_course_finalproject/types"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateNewUser(context.Context, *types.User) error
	GetUser(context.Context, string) (*types.User, error)
	//JWT
	GenerateJWT(context.Context, string, string) (string, error)
	//Get all Users
	GetAllUsers(context.Context) ([]*types.User, error)
}

type userServiceImpl struct {
	repo       repos.UserRepository
	bcryptCost int // cost of bcrypt
	jwtSecret  []byte
}

func NewUserService(repo repos.UserRepository, bcryptCost int) UserService {
	// Load the JWT secret from the .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Couldn't load .env file")
	}

	jwtSecret := []byte(getEnv("BCRYPT_SECRET", "your_default_secret"))

	return &userServiceImpl{
		repo:       repo,
		bcryptCost: bcryptCost,
		jwtSecret:  jwtSecret,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func (u *userServiceImpl) CreateNewUser(c context.Context, user *types.User) error {
	//TODO: Logic to create user - hash the password !!!
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), u.bcryptCost)
	if err != nil {
		return fmt.Errorf("erreur de hachage du mot de passe : %v", err)
	}

	//Store the hashed password in the user object
	user.Password = string(hashedPassword)

	//Call the repository method to create the user
	if err := u.repo.CreateUser(c, user); err != nil {
		return fmt.Errorf("erreur lors de la création de l'utilisateur : %v", err)
	}

	return nil
}

/*REMOVE THIS SERVICE*/
func (u *userServiceImpl) GetUser(c context.Context, id string) (*types.User, error) {
	return u.repo.GetUser(c, id)
}

func (u *userServiceImpl) GetAllUsers(c context.Context) ([]*types.User, error) {
	// Utilisez le UserRepository pour récupérer tous les utilisateurs.
	users, err := u.repo.GetAllUsers(c)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// JWT implementation
func (u *userServiceImpl) GenerateJWT(c context.Context, username, password string) (string, error) {
	//Create a token with custom claims (in this example just a username)
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), //expire after 24hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//Sign the token with your secret key
	tokenString, err := token.SignedString(u.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
