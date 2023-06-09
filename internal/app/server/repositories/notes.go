package repositories

import (
	"context"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/db/models"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

type NotesRepository interface {
	GetNote(userID, noteID int) (models.Note, error)
	CreateNote(note models.Note) (models.Note, error)
	GetUserNotesList(userID int) ([]models.Note, error)
}

type notesRepository struct {
	DBConn *pgx.Conn
}

func NewNotesRepository(conn *pgx.Conn) NotesRepository {
	return notesRepository{DBConn: conn}
}

func (nr notesRepository) GetNote(userID, noteID int) (models.Note, error) {
	ctx := context.Background()
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Select().
		Columns("id", "user_id", "data", "metadata", "data_type", "updated_at", "created_at").
		From("user_data").
		Where(sq.Eq{"id": noteID, "user_id": userID}).
		ToSql()

	if err != nil {
		log.WithFields(log.Fields{"func": "GetNote"}).Errorln(err)
		return models.Note{}, err
	}

	rows, err := nr.DBConn.Query(ctx, sql, args...)
	if err != nil {
		log.WithFields(log.Fields{"func": "GetNote"}).Errorln(err)
		return models.Note{}, err
	}

	defer rows.Close()

	var note models.Note

	for rows.Next() {
		err = rows.Scan(
			&note.ID,
			&note.UserID,
			&note.Data,
			&note.Metadata,
			&note.DataType,
			&note.UpdatedAt,
			&note.CreatedAt)

		if err != nil {
			log.WithFields(log.Fields{"func": "GetNote"}).Errorln(err)
			return models.Note{}, err
		}
	}

	err = rows.Err()
	if err != nil {
		log.WithFields(log.Fields{"func": "GetNote"}).Errorln(err)
		return models.Note{}, err
	}

	return note, nil
}

func (nr notesRepository) CreateNote(note models.Note) (models.Note, error) {
	ctx := context.Background()
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := psql.
		Insert("user_data").
		Columns("user_id", "data", "metadata", "data_type").
		Values(note.UserID, note.Data, note.Metadata, note.DataType).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()

	if err != nil {
		log.WithFields(log.Fields{"func": "CreateNote"}).Errorln(err)
		return models.Note{}, err
	}

	err = nr.DBConn.QueryRow(ctx, sql, args...).
		Scan(&note.ID)

	if err != nil {
		log.WithFields(log.Fields{"func": "CreateNote"}).Errorln(err)
		return models.Note{}, err
	}

	order, err := nr.GetNote(note.UserID, note.ID)

	if err != nil {
		log.WithFields(log.Fields{"func": "CreateOrder"}).Errorln(err)
		return models.Note{}, err
	}

	return order, nil
}

func (nr notesRepository) GetUserNotesList(userID int) ([]models.Note, error) {
	ctx := context.Background()

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.
		Select().
		Columns("id", "user_id", "data", "metadata", "data_type", "updated_at", "created_at").
		From("user_data").
		Where(sq.Eq{"user_id": userID}).
		ToSql()

	if err != nil {
		log.WithFields(log.Fields{"func": "GetUserNotesList"}).Errorln(err)
		return []models.Note{}, err
	}

	rows, err := nr.DBConn.Query(ctx, sql, args...)

	if err != nil {
		log.WithFields(log.Fields{"func": "GetUserNotesList"}).Errorln(err)
		return []models.Note{}, err
	}

	defer rows.Close()

	var notesList []models.Note

	for rows.Next() {
		var note models.Note

		err = rows.Scan(
			&note.ID,
			&note.UserID,
			&note.Data,
			&note.Metadata,
			&note.DataType,
			&note.UpdatedAt,
			&note.CreatedAt)

		if err != nil {
			log.WithFields(log.Fields{"func": "GetUserNotesList"}).Errorln(err)
			return []models.Note{}, err
		}

		notesList = append(notesList, note)
	}
	return notesList, nil
}
