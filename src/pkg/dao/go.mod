module dao

go 1.15

require (
	github.com/mariusmagureanu/burlan/src/pkg/entities v0.0.0-00010101000000-000000000000
	github.com/mattn/go-sqlite3 v1.14.6 // indirect
	github.com/stretchr/testify v1.7.0
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.20.12
)

replace github.com/mariusmagureanu/burlan/src/pkg/entities => ../entities
