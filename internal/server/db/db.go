package db

import (
	"context"
	"database/sql"
	"embed"
	"errors"

	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"gophkeeper/internal/server/dto"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type databaseService struct {
	db        *sql.DB
	pSQL      sq.StatementBuilderType
	logger    *slog.Logger
	txTimeout time.Duration
}

func NewService(db *sql.DB, logger slog.Handler) (*databaseService, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(db)
	svc := &databaseService{
		logger:    slog.New(logger),
		db:        db,
		pSQL:      psql,
		txTimeout: time.Second * 8,
	}
	err := svc.init()
	return svc, err
}

func (s *databaseService) GetUser(ctx context.Context, name string) (dto.User, error) {
	var (
		userId       int64
		userLogin    string
		userPassword string
	)
	query := s.pSQL.Select(
		"id",
		"login",
		"password",
	).From(
		"users",
	).Where(sq.Eq{"login": name})
	row := query.QueryRowContext(ctx)
	if err := row.Scan(&userId, &userLogin, &userPassword); err != nil {
		return dto.User{}, err
	}
	return dto.User{Id: userId, Name: userLogin, Password: userPassword}, nil
}

func (s *databaseService) CreateUser(ctx context.Context, user dto.User) error {
	query := s.pSQL.Insert(
		"users",
	).Columns(
		"login",
		"password",
	).Values(
		user.Name,
		user.Password,
	).Suffix("ON CONFLICT (login) DO NOTHING")

	res, err := query.ExecContext(ctx)
	if err != nil {
		return err
	}
	if affected, err := res.RowsAffected(); err != nil {
		return err
	} else if affected == 0 {
		return errors.New("login is already taken")
	}
	return nil
}

func (s *databaseService) GetFile(ctx context.Context, userId int64, name string) (dto.File, error) {
	var (
		fId     int64
		fName   string
		fStatus dto.FileStatus
	)
	query := s.pSQL.Select(
		"id",
		"file_name",
		"file_status",
	).From(
		"files",
	).Where(sq.And{sq.Eq{"file_name": name}, sq.Eq{"owner_id": userId}})
	row := query.QueryRowContext(ctx)
	if err := row.Scan(&fId, &fName, &fStatus); err != nil {
		return dto.File{}, err
	}
	return dto.File{Id: fId, Name: fName, Status: fStatus}, nil
}

func (s *databaseService) SetFileStatus(ctx context.Context, id int64, status dto.FileStatus) error {
	query := s.pSQL.Update("files").Set("file_status", status).Where(sq.Eq{"id": id})
	res, err := query.ExecContext(ctx)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	return err
}

func (s *databaseService) GetUserFiles(ctx context.Context, userId int64) ([]dto.File, error) {
	query := s.pSQL.Select(
		"id",
		"file_name",
		"file_status",
		"file_size",
	).From("files").Where(
		sq.Eq{"owner_id": userId},
	)
	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	var files []dto.File
	for rows.Next() {
		var (
			id      int64
			fName   string
			fStatus dto.FileStatus
			fSize   int
		)
		err := rows.Scan(&id, &fName, &fStatus, &fSize)
		if err != nil {
			return nil, err
		}
		files = append(files, dto.File{Id: id, Name: fName, Status: fStatus, Size: fSize})
	}
	return files, nil
}

func (s *databaseService) CreateFile(ctx context.Context, file dto.File, userId int64) (int64, error) {
	var fileId int64
	query := s.pSQL.Insert(
		"files",
	).Columns(
		"file_name",
		"file_status",
		"file_size",
		"owner_id",
	).Values(
		file.Name,
		file.Status,
		file.Size,
		userId,
	).Suffix("RETURNING id")
	row := query.QueryRowContext(ctx)
	if err := row.Scan(&fileId); err != nil {
		return 0, err
	}
	return fileId, nil
}

func (s *databaseService) init() error {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(s.db, "migrations")
}
