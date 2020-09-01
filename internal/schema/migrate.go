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
	{
		Version:     3,
		Description: "Add users",
		Script:      CREATE_USERS_TABLE,
	},
	{
		Version:     4,
		Description: "Add user column to products",
		Script: `
ALTER TABLE products
	ADD COLUMN user_id UUID DEFAULT '00000000-0000-0000-0000-000000000000'
`,
	},
}

func Migrate(db *sqlx.DB) error {
	driver := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})
	d := darwin.New(driver, migrations, nil)
	return d.Migrate()
}
