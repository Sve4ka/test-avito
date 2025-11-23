package team

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

func InitTeamRepo(db *postgres.Pg) repo.Team {
	return Repo{db: db}
}

func (r Repo) CheckTeamName(ctx context.Context, teamName string) (bool, error) {
	var count int

	query := `SELECT COUNT(*) FROM users WHERE team_name = $1`

	err := r.db.Pool.QueryRow(ctx, query, teamName).Scan(&count)
	if err != nil {
		return false, cerr.HandlePgErr(err)
	}

	return count == 0, nil
}

func (r Repo) Create(ctx context.Context, team *entity.Team) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return cerr.HandlePgErr(err)
	}

	userQuery := `INSERT INTO users (id, username, team_name, is_active) VALUES ($1, $2, $3, $4)`
	for _, user := range team.Members {
		_, err = tx.Exec(ctx, userQuery, user.UserId, user.Username, team.TeamName, user.IsActive)
		if err != nil {
			if txErr := tx.Rollback(ctx); txErr != nil {
				return cerr.HandlePgErr(txErr)
			}

			return cerr.HandlePgErr(err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return cerr.HandlePgErr(err)
	}

	return nil
}

func (r Repo) Get(ctx context.Context, teamName string) (*entity.Team, error) {
	var team entity.Team
	team.TeamName = teamName

	var member entity.TeamMember

	query := `SELECT id, username, is_active FROM users WHERE team_name = $1`

	rows, err := r.db.Pool.Query(ctx, query, teamName)
	if err != nil {
		return nil, cerr.HandlePgErr(err)
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&member.UserId, &member.Username, &member.IsActive)
		if err != nil {
			return nil, cerr.HandlePgErr(err)
		}

		team.Members = append(team.Members, member)
	}

	if err = rows.Err(); err != nil {
		return nil, cerr.HandlePgErr(err)
	}

	return &team, nil
}
