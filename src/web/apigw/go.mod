module github.com/mariusmagureanu/burlan/src/web/apigw

go 1.15

require (
	github.com/gorilla/mux v1.8.0
	github.com/mariusmagureanu/burlan/src/pkg/dao v0.0.0-00010101000000-000000000000
	github.com/mariusmagureanu/burlan/src/pkg/entities v0.0.0-00010101000000-000000000000
)

replace github.com/mariusmagureanu/burlan/src/pkg/dao => ../../pkg/dao

replace github.com/mariusmagureanu/burlan/src/pkg/entities => ../../pkg/entities
