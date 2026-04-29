module github.com/danilfaer/golang/inventory

go 1.25.5

replace github.com/danilfaer/golang/shared => ../shared

replace github.com/danilfaer/golang/platform => ../platform

require (
	github.com/caarlos0/env/v11 v11.4.0
	github.com/danilfaer/golang/platform v0.0.0-00010101000000-000000000000
	github.com/danilfaer/golang/shared v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
	github.com/stretchr/testify v1.10.0
	go.mongodb.org/mongo-driver/v2 v2.5.1
	google.golang.org/grpc v1.73.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/klauspost/compress v1.17.6 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.2.0 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
