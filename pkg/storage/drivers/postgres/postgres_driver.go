package pgdriver

import _ "github.com/lib/pq"

const Name = "postgres"

type Postgres struct{}

func (Postgres) Name() string {
	return Name
}

func (Postgres) MustUseWithImportedSQLDriver() {}
