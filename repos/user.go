package repos

import (
	"context"
	"fmt"
	"strconv"

	"estiam_golang_api_course_finalproject/types"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetUser(context.Context, string) (*types.User, error)
	CreateUser(context.Context, *types.User) error
	GetAllUsers(context.Context) ([]*types.User, error)
}

type userRepositoryImpl struct {
	dbConn *pgxpool.Pool
}

func NewUserRepository(conn *pgxpool.Pool) UserRepository {
	return &userRepositoryImpl{
		dbConn: conn,
	}
}

// init the .env file
func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Couldn't load .env file")
	}
}

/*REMOVE THIS METHOD*/
const SQL_GET_USER = `
		select
			u.id,
			u.username,
			u.pass
		from
			"user" as u
		where u.id = $1;`

func (repo *userRepositoryImpl) GetUser(c context.Context, userId string) (*types.User, error) {
	rows, err := repo.dbConn.Query(c, SQL_GET_USER, userId)
	if err != nil {
		return nil, fmt.Errorf("error during query to get user: %v", err)
	}

	if rows.Next() {
		user := &types.User{}
		err = rows.Scan(
			&user.Id,
			&user.Username,
			&user.Password,
		)

		if err != nil {
			return nil, err
		}

		return user, nil
	}

	return nil, nil
}

/*get all users methods*/
func (repo *userRepositoryImpl) GetAllUsers(c context.Context) ([]*types.User, error) {
	const SQL_GET_ALL_USERS = `SELECT id, username, password FROM "user";`
	rows, err := repo.dbConn.Query(c, SQL_GET_ALL_USERS)
	if err != nil {
		return nil, fmt.Errorf("error during query to get all users: %v", err)
	}

	var users []*types.User
	for rows.Next() {
		user := &types.User{}
		err = rows.Scan(
			&user.Id,
			&user.Username,
			&user.Password,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

/*IMPLEMENT THIS METHOD*/
// const SQL_INSERT_USER = `?`

func (repo *userRepositoryImpl) CreateUser(c context.Context, user *types.User) error {
	// Vérifiez que l'ID est une chaîne de caractères représentant un nombre entier valide
	userID, err := strconv.Atoi(user.Id)
	if err != nil {
		return fmt.Errorf("l'ID de l'utilisateur n'est pas valide : %v", err)
	}
	//hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}

	//insert new user
	const SQL_INSERT_USER = `
		INSERT INTO "user" (id, username, password) VALUES ($1, $2, $3) RETURNING id;`

	_, err = repo.dbConn.Exec(c, SQL_INSERT_USER, userID, user.Username, hashedPassword)
	if err != nil {
		return fmt.Errorf("error inserting user: %v", err)
	}

	return nil

}
