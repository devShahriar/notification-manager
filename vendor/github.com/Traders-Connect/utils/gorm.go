package utils

import (
	"encoding/json"
	"errors"

	"github.com/go-sql-driver/mysql"
)

// GormErr is a structure representing a Gorm error that can be used for parsing
type GormErr struct {
	Number  int    `json:"Number"`
	Message string `json:"Message"`
}

// ErrGormDuplicateKey error for when there are duplicated keys
type ErrGormDuplicateKey struct{}

// Error returns a string representing the ErrGormDuplicateKey error
func (e *ErrGormDuplicateKey) Error() string {
	return "duplicate key"
}

// ErrGormUnknown ErrGormDuplicateKey error for when there are duplicated keys
type ErrGormUnknown struct{}

// Error returns a string representing the ErrGormUnknown error
func (e *ErrGormUnknown) Error() string {
	return "unknown"
}

// ParseGormErr tries to parse an error to a GormErr
func ParseGormErr(gErr error) error {
	var newError GormErr
	err := json.Unmarshal(([]byte(gErr.Error())), &newError)
	if err != nil {
		return &ErrGormUnknown{}
	}

	switch newError.Number {
	case 1062:
		return &ErrGormDuplicateKey{}
	default:
		return &ErrGormUnknown{}
	}
}

// IsDuplicateKeyError returns if the db error is a duplicate key error
func IsDuplicateKeyError(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}
