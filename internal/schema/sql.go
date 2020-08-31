package schema

const CREATE_PRODUCTS_TABLE = `
	CREATE TABLE products (
		product_id   UUID,
		name         TEXT,
		cost         INT,
		quantity     INT,
		date_created TIMESTAMP,
		date_updated TIMESTAMP,
	
		PRIMARY KEY (product_id)
	);`

const CREATE_SALES_TABLE = `CREATE TABLE sales (
	sale_id      UUID,
	product_id   UUID,
	quantity     INT,
	paid         INT,
	date_created TIMESTAMP,

	PRIMARY KEY (sale_id),
	FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE
);`

const CREATE_USERS_TABLE = `
CREATE TABLE users (
	user_id       UUID,
	name          TEXT,
	email         TEXT UNIQUE,
	roles         TEXT[],
	password_hash TEXT,

	date_created TIMESTAMP,
	date_updated TIMESTAMP,

	PRIMARY KEY (user_id)
);`
