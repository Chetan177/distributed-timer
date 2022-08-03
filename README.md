# distributed-timer
A distributed timer service in Golang with activeMQ

### 
Steps:
1. First Run active mq docker [Link](activemq/README.md)
2. Run Producer `go run producer/producer.go`
3. Run Consumer `go run consumer/consumer.go`
4. Now Execute API

Note Access producer swagger doc at: http://localhost:7070/v1/swagger/index.html