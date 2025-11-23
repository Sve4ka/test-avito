package pullRequest

import (
	"context"
	"errors"
	"time"

	"avito/internal/cerr"
	"avito/internal/entity"
	"avito/internal/postgres"
	"avito/internal/repo"
	"github.com/jackc/pgx/v5"
)

type Repo struct {
	db *postgres.Pg
}

func InitPullRequestRepo(db *postgres.Pg) repo.PullRequest {
	return Repo{db: db}
}

func (r Repo) Create(ctx context.Context, pullRequestCreate *entity.PullRequestCreate) (*entity.PullRequest, error) {
	creatAt := time.Now().UTC()

	pullRequest := entity.PullRequest{
		AuthorId:        pullRequestCreate.AuthorId,
		PullRequestName: pullRequestCreate.PullRequestName,
		PullRequestId:   pullRequestCreate.PullRequestId,
		Status:          entity.PRStatusOPEN,
		CreatedAt:       &creatAt,
	}

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return nil, cerr.HandlePgErr(err)
	}

	choiceQuery := `SELECT u.id
FROM users AS u
    LEFT JOIN reviewers AS r ON u.id = r.reviewer_id
    LEFT JOIN pull_requests AS pr ON r.pull_request_id = pr.id AND pr.merged_at IS NULL
WHERE u.team_name = (SELECT team_name FROM users WHERE id = $1) AND u.id != $1 AND u.is_active = true
GROUP BY u.id
ORDER BY COUNT(r.reviewer_id)
LIMIT 2;`

	rows, err := r.db.Pool.Query(ctx, choiceQuery, pullRequestCreate.AuthorId)
	if err != nil {
		return nil, cerr.HandlePgErr(err)
	}
	defer rows.Close()

	assignQuery := `INSERT INTO reviewers (pull_request_id, reviewer_id) VALUES ($1, $2);`

	for rows.Next() {
		var userID string

		err = rows.Scan(&userID)
		if err != nil {
			if txErr := tx.Rollback(ctx); txErr != nil {
				return nil, cerr.HandlePgErr(txErr)
			}

			return nil, cerr.HandlePgErr(err)
		}

		_, err = tx.Exec(ctx, assignQuery, pullRequestCreate.PullRequestId, userID)
		if err != nil {
			if txErr := tx.Rollback(ctx); txErr != nil {
				return nil, cerr.HandlePgErr(txErr)
			}

			return nil, cerr.HandlePgErr(err)
		}

		pullRequest.AssignedReviewers = append(pullRequest.AssignedReviewers, userID)
	}

	createQuery := `INSERT INTO pull_requests (id, name, author_id, status_id, create_at) VALUES ($1, $2, $3, (SELECT id FROM statuses WHERE name = $4), $5);`

	_, err = tx.Exec(ctx, createQuery, pullRequestCreate.PullRequestId, pullRequestCreate.PullRequestName, pullRequestCreate.AuthorId, entity.PRStatusOPEN, creatAt)
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

	return &pullRequest, nil
}

func (r Repo) Merge(ctx context.Context, pullRequestID string) (*entity.PullRequest, error) {
	mergedAt := time.Now().UTC()
	pullRequest := entity.PullRequest{
		PullRequestId: pullRequestID,
		MergedAt:      &mergedAt,
		Status:        entity.PRStatusMERGED,
	}

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return nil, cerr.HandlePgErr(err)
	}

	updateQuery := `UPDATE pull_requests as pr SET status_id=(SELECT id FROM statuses WHERE name = $1), merged_at = $2 WHERE pr.id = $3
                                      returning name, author_id, create_at`

	err = tx.QueryRow(ctx, updateQuery, entity.PRStatusMERGED, mergedAt, pullRequestID).Scan(&pullRequest.PullRequestName, &pullRequest.AuthorId, &pullRequest.CreatedAt)
	if err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			return nil, cerr.HandlePgErr(txErr)
		}

		return nil, cerr.HandlePgErr(err)
	}

	var reviewer string

	query := `SELECT reviewer_id FROM reviewers WHERE pull_request_id = $1`

	rows, err := r.db.Pool.Query(ctx, query, pullRequestID)
	if err != nil {
		return nil, cerr.HandlePgErr(err)
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&reviewer)
		if err != nil {
			return nil, cerr.HandlePgErr(err)
		}

		pullRequest.AssignedReviewers = append(pullRequest.AssignedReviewers, reviewer)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, cerr.HandlePgErr(err)
	}

	return &pullRequest, nil
}

func (r Repo) Reassign(ctx context.Context, pullRequestID string, oldUserID string) (*entity.PullRequest, string, error) {
	pullRequest := entity.PullRequest{
		PullRequestId: pullRequestID,
	}

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return nil, "", cerr.HandlePgErr(err)
	}

	var cnt int

	checkUserQuery := `SELECT COUNT(*) from users where id=$1`

	err = tx.QueryRow(ctx, checkUserQuery, oldUserID).Scan(&cnt)
	if err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			return nil, "", cerr.HandlePgErr(txErr)
		}

		return nil, "", cerr.HandlePgErr(err)
	}

	if cnt == 0 {
		return nil, "", cerr.CustomError{Err: err, ErrType: cerr.NOT_FOUND}
	}

	PRQuery := `SELECT merged_at, name, author_id from pull_requests where id = $1`

	err = r.db.Pool.QueryRow(ctx, PRQuery, pullRequestID).Scan(&pullRequest.MergedAt, &pullRequest.PullRequestName, &pullRequest.AuthorId)
	if err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			return nil, "", cerr.HandlePgErr(txErr)
		}

		return nil, "", cerr.HandlePgErr(err)
	}

	if pullRequest.MergedAt != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			return nil, "", cerr.HandlePgErr(txErr)
		}

		return nil, "", cerr.CustomError{
			Err:     err,
			ErrType: cerr.PR_MERGED,
		}
	}

	pullRequest.Status = entity.PRStatusOPEN

	OldReviewersQuery := `SELECT reviewer_id from reviewers where pull_request_id = $1`
	haveOldReviewer := false

	var OldReviewers []string

	var secondUser string

	var newUser string

	rows, err := r.db.Pool.Query(ctx, OldReviewersQuery, pullRequestID)
	if err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			return nil, "", cerr.HandlePgErr(txErr)
		}

		return nil, "", cerr.HandlePgErr(err)
	}

	defer rows.Close()

	for rows.Next() {
		var userID string

		err = rows.Scan(&userID)
		if err != nil {
			if txErr := tx.Rollback(ctx); txErr != nil {
				return nil, "", cerr.HandlePgErr(txErr)
			}

			return nil, "", cerr.HandlePgErr(err)
		}

		if userID == oldUserID {
			haveOldReviewer = true
		} else {
			secondUser = userID
		}

		OldReviewers = append(OldReviewers, userID)
	}

	if !haveOldReviewer {
		if txErr := tx.Rollback(ctx); txErr != nil {
			return nil, "", cerr.HandlePgErr(txErr)
		}

		return nil, "", cerr.CustomError{
			Err:     err,
			ErrType: cerr.NOT_ASSIGNED,
		}
	}

	if len(OldReviewers) == 1 {
		OldReviewers = append(OldReviewers, OldReviewers[0])
	}

	choiceQuery := `SELECT u.id
FROM users AS u
    LEFT JOIN reviewers AS r ON u.id = r.reviewer_id
    LEFT JOIN pull_requests AS pr ON r.pull_request_id = pr.id AND pr.merged_at IS NULL
WHERE u.team_name = (SELECT team_name FROM users WHERE id = $1) AND u.id != $1 And u.id != $2 and u.id != $3 AND u.is_active = true
GROUP BY u.id
ORDER BY COUNT(r.reviewer_id)
LIMIT 1;`

	err = r.db.Pool.QueryRow(ctx, choiceQuery, pullRequest.AuthorId, OldReviewers[0], OldReviewers[1]).Scan(&newUser)
	if err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			return nil, "", cerr.HandlePgErr(txErr)
		}

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, "", cerr.CustomError{
				Err:     err,
				ErrType: cerr.NO_CANDIDATE,
			}
		}

		return nil, "", cerr.HandlePgErr(err)
	}

	defer rows.Close()

	assignQuery := `UPDATE reviewers SET reviewer_id = $1 WHERE pull_request_id = $2 AND reviewer_id=$3;`

	_, err = tx.Exec(ctx, assignQuery, newUser, pullRequestID, oldUserID)
	if err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			return nil, "", cerr.HandlePgErr(txErr)
		}

		return nil, "", cerr.HandlePgErr(err)
	}

	pullRequest.AssignedReviewers = []string{newUser, secondUser}

	err = tx.Commit(ctx)
	if err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			return nil, "", cerr.HandlePgErr(txErr)
		}

		return nil, "", cerr.HandlePgErr(err)
	}

	return &pullRequest, newUser, nil
}
