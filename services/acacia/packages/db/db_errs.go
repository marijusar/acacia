package db

// PostgreSQL error codes
const (
	// UniqueViolation occurs when a unique constraint is violated
	PgErrUniqueViolation = "23505"

	// ForeignKeyViolation occurs when a foreign key constraint is violated
	PgErrForeignKeyViolation = "23503"

	// NotNullViolation occurs when a NOT NULL constraint is violated
	PgErrNotNullViolation = "23502"

	// CheckViolation occurs when a CHECK constraint is violated
	PgErrCheckViolation = "23514"
)
