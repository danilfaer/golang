module github.com/danilfaer/golang/payment

replace github.com/danilfaer/golang/shared => ../shared

replace github.com/danilfaer/golang/platform => ../platform

go 1.25.5

require (
	github.com/brianvoe/gofakeit/v7 v7.3.0
	github.com/caarlos0/env/v11 v11.4.0
	github.com/danilfaer/golang/platform v0.0.0-00010101000000-000000000000
	github.com/danilfaer/golang/shared v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/stretchr/testify v1.10.0
	google.golang.org/grpc v1.73.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
