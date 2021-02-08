module github.com/mariusmagureanu/burlan/src/web/apigw

go 1.15

require (
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/mariusmagureanu/burlan/src/pkg/auth v0.0.0-20210208095453-e0713562b900
	github.com/mariusmagureanu/burlan/src/pkg/dao v0.0.0-00010101000000-000000000000
	github.com/mariusmagureanu/burlan/src/pkg/entities v0.0.0-00010101000000-000000000000
	github.com/mariusmagureanu/burlan/src/pkg/errors v0.0.0-20210208095453-e0713562b900
	github.com/mariusmagureanu/burlan/src/pkg/log v0.0.0-00010101000000-000000000000
	github.com/segmentio/kafka-go v0.4.9
)

replace github.com/mariusmagureanu/burlan/src/pkg/dao => ../../pkg/dao

replace github.com/mariusmagureanu/burlan/src/pkg/errors => ../../pkg/errors

replace github.com/mariusmagureanu/burlan/src/pkg/entities => ../../pkg/entities

replace github.com/mariusmagureanu/burlan/src/pkg/log => ../../pkg/log

replace github.com/mariusmagureanu/burlan/src/pkg/auth => ../../pkg/auth
