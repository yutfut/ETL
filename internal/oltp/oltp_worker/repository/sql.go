package repository

const (
	selectMeta = `
		select
			last_insert_id,
    		last_update_at
		from client.etl
		where id = 1;
	`

	updateMeta = `
		update client.etl
		set
			last_insert_id = $1,
			last_update_at = $2
		where id = 1
		returning last_insert_id, last_update_at;
	`

	u1 = `
		update client.etl
		set
			last_insert_id = $1
		where id = 1
		returning last_insert_id;
	`

	u2 = `
		update client.etl
		set
    		last_update_at = $1
		where id = 1
		returning last_update_at;
	`

	selectByID = `
		select
			id,
			name,
			settlement,
			margin_algorithm,
			gateway,
			vendor,
			is_active,
			is_pro,
			is_interbank,
			create_at,
			update_at
		from client.client
		where id > $1;
	`

	selectByUpdateAT = `
		select
			id,
			name,
			settlement,
			margin_algorithm,
			gateway,
			vendor,
			is_active,
			is_pro,
			is_interbank,
			create_at,
			update_at
		from client.client
		where update_at > $1;
	`
)
