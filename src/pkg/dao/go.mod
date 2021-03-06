module github.com/mariusmagureanu/burlan/src/pkg/dao

go 1.15

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/google/uuid v1.2.0
	github.com/mariusmagureanu/burlan/src/pkg/entities v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	gorm.io/gorm v1.20.12
)

replace github.com/mariusmagureanu/burlan/src/pkg/entities => ../entities
