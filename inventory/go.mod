module github.com/danilfaer/golang/inventory

go 1.24.4

replace github.com/danilfaer/golang/shared => ../shared

require (
	github.com/danilfaer/golang/shared v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.73.0
	google.golang.org/protobuf v1.36.6
)

require (
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
)
