package pgrepo

const referralCodeQueryCreate = `
	INSERT INTO referralCode (
		user_id,
	) VALUES (
		$1, $2
	)
	RETURNING 
		id,
		user_id
`

const referralCodeQueryDelete = `
	UPDATE 
		referralCode
	SET 
		deleted_at = COALESCE(deleted_at, NOW())
	WHERE 
		id = $1
	RETURNING 
		deleted_at;
`

const referralCodeQueryGetByUsername = `
	SELECT     
		refcode.id, 
		refcode.user_id
	FROM 
		users JOIN referral_code refcode ON users.username = $1 AND users.id = refcode.id
	WHERE 
		deleted_at IS NULL;
`
