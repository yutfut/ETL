package generator

const (
	insert = `
	insert into client.client
	(
		name,
		settlement,
		margin_algorithm,
		gateway,
		vendor,
		is_active,
		is_pro,
		is_interbank
	) values (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8
	)
	`

	update = `
	update client.client
	set
		name = $1,
		update_at = now()
	where id = $2
	`
)