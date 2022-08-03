package xmysql

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func ErrorNumber(err error) uint16 {
	var merr *mysql.MySQLError
	if errors.As(err, &merr) {
		return merr.Number
	}
	return 0
}

func IsForeignKeyError(err error, constraint string) bool {
	var merr *mysql.MySQLError
	if errors.As(err, &merr) {
		if merr.Number == 1452 && strings.Contains(merr.Message, fmt.Sprintf("`%s`", constraint)) {
			return true
		}
	}
	return false
}

func IsDuplicateError(err error, field string) bool {
	var merr *mysql.MySQLError
	if errors.As(err, &merr) {
		if merr.Number == 1062 && strings.Contains(merr.Message, fmt.Sprintf("'%s'", field)) {
			return true
		}
	}
	return false
}
