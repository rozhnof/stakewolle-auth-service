package pgrepo

const userQueryCreate = `
	INSERT INTO users (
		username,
		referrer_id,
		hash_password
	) VALUES (
	 	$1, $2, $3
	)
	RETURNING 
		id,
		username,
		referrer_id,
		hash_password
`

const userQueryDelete = `
	UPDATE 
		users
	SET 
		deleted_at = COALESCE(deleted_at, NOW())
	WHERE 
		id = $1
	RETURNING 
		deleted_at;
`

const userQueryGetByID = `
	SELECT     
		id, 
		username,
		referrer_id,
		hash_password
	FROM 
		users
	WHERE
		id = $1
`

const userQueryGetByUsername = `
	SELECT     
		id,
		username,
		referrer_id,
		hash_password
	FROM 
		users
	WHERE
		username = $1
`

const userQueryList = `
	SELECT     
		id,
		username,
		referrer_id,
		hash_password
	FROM 
		users
`

const userQueryUpdate = `
	UPDATE 
		users 
	SET  
		username = $2,
		referrer_id = $3,
		hash_password = $4
	WHERE 
		id = $1
	RETURNING 
		id,
		username,
		referrer_id,
		hash_password
`
