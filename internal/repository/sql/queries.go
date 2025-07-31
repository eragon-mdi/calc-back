package sqlrepo

const getCalcsWithMax = `
SELECT 
	id, expression, result
FROM
	calculations
ORDER BY id DESC
LIMIT $1
`

const getCalcById = `
SELECT 
	id, expression, result
FROM
	calculations
WHERE id = $1
`

const deleteCalcById = `
DELETE 
FROM
	calculations
WHERE id = $1
`

const insertCalc = `
INSERT INTO 
	calculations
	(id, expression, result)
VALUES
	($1, $2, $3)
RETURNING
	id, expression, result
`

const updateCalc = `
UPDATE
	calculations
SET
	expression = $2, result = $3
WHERE
	id = $1
RETURNING id, expression, result
`
