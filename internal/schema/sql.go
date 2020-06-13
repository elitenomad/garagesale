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
