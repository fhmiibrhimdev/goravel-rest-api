package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20210101000001CreateUsersTable struct{}

// Signature The unique signature for the migration.
func (r *M20210101000001CreateUsersTable) Signature() string {
	return "20210101000001_create_users_table"
}

// Up Run the migrations.
func (r *M20210101000001CreateUsersTable) Up() error {
	if !facades.Schema().HasTable("users") {
		return facades.Schema().Create("users", func(table schema.Blueprint) {
			table.ID()
			table.String("name")
			table.String("email")
			table.String("password")
			table.Integer("active").Default(1)
			table.Timestamps()
		})
	}
	return nil
}

// Down Reverse the migrations.
func (r *M20210101000001CreateUsersTable) Down() error {
	return facades.Schema().DropIfExists("users")
}
