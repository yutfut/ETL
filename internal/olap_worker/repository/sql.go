package repository

const (
	insertBatch = `
	INSERT INTO client.client
	`

	update = `
	ALTER TABLE client.client
	update
	    name = $2,
	    settlement = $3,
		margin_algorithm = $4,
		gateway = $5,
		vendor = $6,
		is_active = $7,
		is_pro = $8,
		is_interbank = $9,
		update_at = $10,
	where id = $1;
	`
)
