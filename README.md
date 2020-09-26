Run the program
1. install go 1.15 on your computer
2. run `go run cmd/main.go`

Send HTTP request
the default port is 8088
temperature monitor url is `/tmp`
quality monitor url is `/qlt`

Example
send quality monitor request to `http://127.0.0.1:8088/qlt`
Or you could use go tests in cmd/main_test.go
`TestTemperature and TestQuality`
these test will fire 1000 requests for 10 devices in parallel 
