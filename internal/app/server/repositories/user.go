package repositories

import (
	"context"
	"errors"

	"github.com/IgorChicherin/gophkeeper/internal/pkg/authlib"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/db/models"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	HasLogin(login string) (bool, error)
	GetUser(login string) (models.User, error)
	CreateUser(login, password string) (models.User, error)
	Validate(hash string) (bool, error)
	DecodeToken(token string) (string, string, error)
}

type userRepo struct {
	DBConn      *pgx.Conn
	AuthService authlib.AuthService
}

func NewUserRepository(
	conn *pgx.Conn,
	service authlib.AuthService,
) UserRepository {
	return userRepo{DBConn: conn, AuthService: service}
}

func (ur userRepo) GetUser(login string) (models.User, error) {
	ctx := context.Background()
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Select().
		Columns("id", "login", "password", "created_at").
		From("users").
		Where(sq.Eq{"login": login}).
		ToSql()

	if err != nil {
		log.WithFields(log.Fields{"func": "GetUser"}).Errorln(err)
		return models.User{}, err
	}

	rows, err := ur.DBConn.Query(ctx, sql, args...)
	if err != nil {
		log.WithFields(log.Fields{"func": "GetUser"}).Errorln(err)
		return models.User{}, err
	}

	defer rows.Close()

	var u models.User

	for rows.Next() {
		err = rows.Scan(&u.UserID, &u.Login, &u.Password, &u.CreatedAt)
		if err != nil {
			log.WithFields(log.Fields{"func": "GetUser"}).Errorln(err)
			return models.User{}, err
		}
	}

	err = rows.Err()

	if err != nil {
		log.WithFields(log.Fields{"func": "GetUser"}).Errorln(err)
		return models.User{}, err
	}

	return u, nil
}

func (ur userRepo) CreateUser(login, password string) (models.User, error) {
	ctx := context.Background()

	pwdHash := ur.AuthService.GetHash(password)

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.
		Insert("users").
		Columns("login", "password").
		Values(login, pwdHash)
	sql, args, err := query.ToSql()

	if err != nil {
		log.WithFields(log.Fields{"func": "CreateUser"}).Errorln(err)
		return models.User{}, err
	}

	_, err = ur.DBConn.Exec(ctx, sql, args...)
	if err != nil {
		log.WithFields(log.Fields{"func": "CreateUser"}).Errorln(err)
		return models.User{}, err
	}

	return models.User{Login: login, Password: pwdHash}, nil
}

func (ur userRepo) HasLogin(login string) (bool, error) {
	ctx := context.Background()

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.
		Select("COALESCE(COUNT(*), 0) as count").
		From("users").
		Where(sq.Eq{"login": login}).
		ToSql()

	if err != nil {
		log.WithFields(log.Fields{"func": "HasLogin"}).Errorln(err)
		return false, err
	}

	rows, err := ur.DBConn.Query(ctx, sql, args...)

	if err != nil {
		log.WithFields(log.Fields{"func": "HasLogin"}).Errorln(err)
		return false, err
	}

	defer rows.Close()

	var count int
	rows.Next()
	err = rows.Scan(&count)

	if err != nil {
		log.WithFields(log.Fields{"func": "HasLogin"}).Errorln(err)
		return false, err
	}

	return count > 0, nil
}

func (ur userRepo) Validate(hash string) (bool, error) {
	login, hash, err := ur.AuthService.DecodeToken(hash)

	if err != nil {
		log.WithFields(log.Fields{"func": "Validate"}).Errorln(err)
		return false, err
	}

	user, err := ur.GetUser(login)
	if err != nil {
		return false, err
	}
	return user.Password == hash, nil
}

func (ur userRepo) DecodeToken(token string) (string, string, error) {
	return ur.AuthService.DecodeToken(token)
}
