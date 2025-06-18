package owner

// list of sql insert statements of various types
const (
	stmtInsert string = `INSERT INTO owner (name,created_at) VALUES (:name, :created_at) ON CONFLICT (name) DO UPDATE SET name=excluded.name RETURNING id;`
)
