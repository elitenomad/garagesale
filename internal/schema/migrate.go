package schema

import (
	"github.com/GuiaBolso/darwin"
	"github.com/jmoiron/sqlx"
)

var migrations = []darwin.Migration{
	{
		Version:     1,
		Description: "Add products",
		Script:      CREATE_PRODUCTS_TABLE,
	},
	{
		Version:     2,
		Description: "Add sales",
		Script:      CREATE_SALES_TABLE,
	},
}

func Migrate(db *sqlx.DB) error {
	driver := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})
	d := darwin.New(driver, migrations, nil)
	return d.Migrate()
}
