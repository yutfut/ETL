package repository

const(
	insertBatch = `
	INSERT INTO client.client
	`

	update = `
	ALTER TABLE client.client
	update name='hello'
	where name='name';
	`
)