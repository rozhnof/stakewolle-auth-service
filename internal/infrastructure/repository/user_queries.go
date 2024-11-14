package repository

const userQueryCreate = `
	INSERT INTO users (
		username,
		hash_password
	) VALUES (
	 	$1, $2
	)
	RETURNING 
		id,
		username,
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
		hash_password
	FROM 
		users
`

const userQueryUpdate = `
	UPDATE 
		users 
	SET  
		username = $2,
		hash_password = $3
	WHERE 
		id = $1
	RETURNING 
		id,
		username,
		hash_password
`
