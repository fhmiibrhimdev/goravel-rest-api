package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250901143706CreatePostsTable struct{}

// Signature The unique signature for the migration.
func (r *M20250901143706CreatePostsTable) Signature() string {
	return "20250901143706_create_posts_table"
}

// Up Run the migrations.
func (r *M20250901143706CreatePostsTable) Up() error {
	if !facades.Schema().HasTable("posts") {
		return facades.Schema().Create("posts", func(table schema.Blueprint) {
            table.ID()
            table.String("title")
            table.Text("body")
            table.Text("publish_date")
            table.TimestampsTz()
        })
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250901143706CreatePostsTable) Down() error {
 	return facades.Schema().DropIfExists("posts")
}
