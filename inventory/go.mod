module inventory

go 1.25.1

require (
	github.com/mllbll/space-manufacture/shared v0.0.0-20251122152415-84e18afe73f2
	google.golang.org/grpc v1.76.0
	google.golang.org/protobuf v1.36.10
)

replace github.com/mllbll/space-manufacture/shared => ../shared

require (
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
)
