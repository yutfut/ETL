package repository

const (
	insertBatch = `
	INSERT INTO client.client
	`

	update = `
	ALTER TABLE client.client
	update
	    name = $3,
	    settlement = $4,
		margin_algorithm = $5,
		gateway = $6,
		vendor = $7,
		is_active = $8,
		is_pro = $9,
		is_interbank = $10,
		update_at = $11,
	where postgresql_id = $1 and id = $2;
	`
)
