package pgrepo

const sessionQueryCreate = `
	INSERT INTO session (
		user_id,
		refresh_token,
		expired_at,
		is_revoked
	) VALUES (
		$1, $2, $3, $4
	)
	RETURNING 
		id,
		user_id,
		refresh_token,
		expired_at,
		is_revoked
`

const sessionQueryDelete = `
	UPDATE 
		session
	SET 
		deleted_at = COALESCE(deleted_at, NOW())
	WHERE 
		id = $1
	RETURNING 
		deleted_at;
`

const sessionQueryGetByID = `
	SELECT     
		id, 
		user_id,
		refresh_token,
		expired_at,
		is_revoked
	FROM 
		session
	WHERE
		id = $1
`

const sessionQueryGetByRefreshToken = `
	SELECT     
		id, 
		user_id,
		refresh_token,
		expired_at,
		is_revoked
	FROM 
		session
	WHERE
		refresh_token = $1
`

const sessionQueryRevoke = `
	UPDATE 
		session 
	SET  
		is_revoked = TRUE
	WHERE 
		user_id = $1
`

const sessionQueryUpdate = `
	UPDATE 
		session 
	SET  
		user_id = $2,
		refresh_token = $3,
		expired_at = $4,
		is_revoked = $5
	WHERE 
		id = $1
	RETURNING 
		id,
		user_id,
		refresh_token,
		expired_at,
		is_revoked
`
