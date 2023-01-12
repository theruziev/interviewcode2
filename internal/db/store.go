package db

import (
	"github.com/Masterminds/squirrel"
)

var pgsql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
