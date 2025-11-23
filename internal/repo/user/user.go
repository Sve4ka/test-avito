package user

import (
	"context"

	"avito/internal/cerr"
	"avito/internal/entity"
	"avito/internal/postgres"
	"avito/internal/repo"
)

type Repo struct {
	db *postgres.Pg
}

func InitUserRepo(db *postgres.Pg) repo.User {
	return Repo{db: db}
}

func (r Repo) SetIsActive(ctx context.Context, userID string, isActive bool) (*entity.User, error) {
	var user entity.User

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return nil, cerr.HandlePgErr(err)
	}

	updateQuery := `UPDATE users SET is_active = $1 WHERE id = $2`

	_, err = tx.Exec(ctx, updateQuery, isActive, userID)
	if err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			return nil, cerr.HandlePgErr(txErr)
		}

		return nil, cerr.HandlePgErr(err)
	}

	selectQuery := `SELECT id, username, team_name, is_active FROM users WHERE id = $1`

	err = tx.QueryRow(ctx, selectQuery, userID).Scan(&user.UserId, &user.Username, &user.TeamName, &user.IsActive)
	if err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			return nil, cerr.HandlePgErr(txErr)
		}

		return nil, cerr.HandlePgErr(err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, cerr.HandlePgErr(err)
	}

	return &user, nil
}

func (r Repo) GetReview(ctx context.Context, userID string) ([]entity.PullRequestShort, error) {
	var prs []entity.PullRequestShort

	var pr entity.PullRequestShort

	var count int

	checkQuery := `SELECT COUNT(*) FROM users WHERE id = $1`

	err := r.db.Pool.QueryRow(ctx, checkQuery, userID).Scan(&count)
	if err != nil {
		return nil, cerr.HandlePgErr(err)
	}

	if count == 0 {
		return nil, cerr.CustomError{
			Err:     err,
			ErrType: cerr.NOT_FOUND,
		}
	}

	query := `SELECT pr.author_id, pr.id, pr.name, s.name FROM pull_requests as pr
    INNER JOIN reviewers AS r on pr.id = r.pull_request_id
    INNER JOIN statuses s on s.id = pr.status_id
    WHERE r.reviewer_id = $1`

	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, cerr.HandlePgErr(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&pr.AuthorId, &pr.PullRequestId, &pr.PullRequestName, &pr.Status)
		if err != nil {
			return nil, cerr.HandlePgErr(err)
		}

		prs = append(prs, pr)
	}

	return prs, nil
}
