package stat

import (
	"context"
	"time"

	"avito/internal/cerr"
	"avito/internal/entity"
	"avito/internal/postgres"
	"avito/internal/repo"
)

type Repo struct {
	db *postgres.Pg
}

func InitStatRepo(db *postgres.Pg) repo.Stat {
	return Repo{db: db}
}

func (r Repo) User(ctx context.Context, userID string) (*entity.UserStat, error) {
	var user entity.UserStat

	var cntMerged, duration float64

	var creatAt, mergedAt *time.Time

	query := `SELECT pr.create_at, pr.merged_at, u.is_active, u.id FROM pull_requests as pr
    INNER JOIN reviewers AS r on pr.id = r.pull_request_id
    INNER JOIN users AS u on u.id = r.reviewer_id
    WHERE r.reviewer_id = $1`

	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, cerr.HandlePgErr(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&creatAt, &mergedAt, &user.IsActive, &user.UserId)
		if err != nil {
			return nil, cerr.HandlePgErr(err)
		}

		user.CountPr++

		if mergedAt != nil {
			duration += mergedAt.Sub(*creatAt).Hours()
			cntMerged++
		}
	}

	if user.UserId == "" {
		query = `SELECT is_active, id from users where id = $1`

		err = r.db.Pool.QueryRow(ctx, query, userID).Scan(&user.IsActive, &user.UserId)
		if err != nil {
			return nil, cerr.HandlePgErr(err)
		}
	}

	if cntMerged != 0 {
		avg := duration / cntMerged
		user.AvgDuration = &avg
	}

	return &user, nil
}

func (r Repo) Team(ctx context.Context, teamName string) (*entity.TeamStat, error) {
	users := []entity.UserStat{}

	var userID string

	var cntMerged, duration, avgDuration float64

	query := `SELECT id FROM users WHERE team_name = $1`

	rows, err := r.db.Pool.Query(ctx, query, teamName)
	if err != nil {
		return nil, cerr.HandlePgErr(err)
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			return nil, cerr.HandlePgErr(err)
		}

		user, err := r.User(ctx, userID)
		if err != nil {
			return nil, cerr.HandlePgErr(err)
		}

		if user.AvgDuration != nil {
			duration += *user.AvgDuration
			cntMerged++
		}

		users = append(users, *user)
	}

	if cntMerged != 0 {
		avgDuration = duration / cntMerged
	} else {
		avgDuration = -1
	}

	return &entity.TeamStat{
		TeamName:    teamName,
		UsersStat:   users,
		AvgDuration: avgDuration,
	}, nil
}
