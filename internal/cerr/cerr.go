package cerr

import (
	"errors"
	"fmt"
	"net/http"

	"avito/internal/gen"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type ErrorType struct {
	code    string
	message string
}

type CustomError struct {
	Err     error
	ErrType ErrorType
}

func (c CustomError) Error() string {
	return fmt.Sprintf("Type: %v, Error: %v", c.ErrType.message, c.Err)
}

const (
	M_TEAM_EXISTS  string = "team_name already exists"
	M_PR_EXISTS    string = "PR id already exists"
	M_PR_MERGED    string = "cannot reassign on merged PR"
	M_NOT_ASSIGNED string = "reviewer is not assigned to this PR"
	M_NO_CANDIDATE string = "no active replacement candidate in team"
	M_NOT_FOUND    string = "data not found"
	M_SERVER       string = "error in service work"
)

var (
	TEAM_EXISTS  = ErrorType{"TEAM_EXIST", M_TEAM_EXISTS}
	PR_EXISTS    = ErrorType{"PR_EXISTS", M_PR_EXISTS}
	PR_MERGED    = ErrorType{"PR_MERGED", M_PR_MERGED}
	NOT_ASSIGNED = ErrorType{"NOT_ASSIGNED", M_NOT_ASSIGNED}
	NO_CANDIDATE = ErrorType{"NO_CANDIDATE", M_NO_CANDIDATE}
	NOT_FOUND    = ErrorType{"NOT_FOUND", M_NOT_FOUND}
	SERVER       = ErrorType{"SERVER", M_SERVER}
)

var ErrServerTime = errors.New(M_SERVER)

func HandlePgErr(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return CustomError{
			Err:     err,
			ErrType: NOT_FOUND,
		}
	}

	var pgErr *pgconn.PgError

	if !errors.As(err, &pgErr) {
		return err
	}

	switch pgErr.Code {
	case "23505":
		switch pgErr.ConstraintName {
		case "pull_requests_pkey":
			return CustomError{
				Err:     err,
				ErrType: PR_EXISTS,
			}
		default:
			err = CustomError{
				Err:     err,
				ErrType: SERVER,
			}
		}
	case "23503":
		if pgErr.ConstraintName == "pull_requests_author_id_fkey" {
			return CustomError{
				Err:     err,
				ErrType: NOT_FOUND,
			}
		}

		return CustomError{
			Err:     err,
			ErrType: SERVER,
		}
	default:
		err = CustomError{
			Err:     err,
			ErrType: SERVER,
		}
	}

	return err
}

func HandleErrs(err error) (int, gen.ErrorResponse) {
	var Cerr CustomError
	if errors.As(err, &Cerr) {
		switch {
		case Cerr.ErrType == TEAM_EXISTS:
			return http.StatusBadRequest, gen.ErrorResponse{
				Error: struct {
					Code    gen.ErrorResponseErrorCode `json:"code"`
					Message string                     `json:"message"`
				}{Code: gen.TEAMEXISTS, Message: M_TEAM_EXISTS},
			}
		case Cerr.ErrType == NOT_FOUND:
			return http.StatusNotFound, gen.ErrorResponse{
				Error: struct {
					Code    gen.ErrorResponseErrorCode `json:"code"`
					Message string                     `json:"message"`
				}{Code: gen.NOTFOUND, Message: M_NOT_FOUND},
			}
		case Cerr.ErrType == PR_EXISTS:
			return http.StatusConflict, gen.ErrorResponse{
				Error: struct {
					Code    gen.ErrorResponseErrorCode `json:"code"`
					Message string                     `json:"message"`
				}{Code: gen.PREXISTS, Message: M_PR_EXISTS},
			}
		case Cerr.ErrType == PR_MERGED:
			return http.StatusConflict, gen.ErrorResponse{
				Error: struct {
					Code    gen.ErrorResponseErrorCode `json:"code"`
					Message string                     `json:"message"`
				}{Code: gen.PRMERGED, Message: M_PR_MERGED},
			}
		case Cerr.ErrType == NOT_ASSIGNED:
			return http.StatusConflict, gen.ErrorResponse{
				Error: struct {
					Code    gen.ErrorResponseErrorCode `json:"code"`
					Message string                     `json:"message"`
				}{Code: gen.NOTASSIGNED, Message: M_NOT_ASSIGNED},
			}
		case Cerr.ErrType == NO_CANDIDATE:
			return http.StatusConflict, gen.ErrorResponse{
				Error: struct {
					Code    gen.ErrorResponseErrorCode `json:"code"`
					Message string                     `json:"message"`
				}{Code: gen.NOCANDIDATE, Message: M_NO_CANDIDATE},
			}
		default:
			return http.StatusTeapot, gen.ErrorResponse{
				Error: struct {
					Code    gen.ErrorResponseErrorCode `json:"code"`
					Message string                     `json:"message"`
				}{Code: "SERVER", Message: M_SERVER},
			}
		}
	}

	return http.StatusTeapot, gen.ErrorResponse{
		Error: struct {
			Code    gen.ErrorResponseErrorCode `json:"code"`
			Message string                     `json:"message"`
		}{Code: "SERVER", Message: M_SERVER},
	}
}
