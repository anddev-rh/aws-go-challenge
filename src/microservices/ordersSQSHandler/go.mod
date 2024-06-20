require (
	github.com/aws/aws-lambda-go v1.47.0
	github.com/aws/aws-sdk-go v1.54.4
	microservices/utils v0.0.0-00010101000000-000000000000
)

replace gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.2.8

replace microservices/utils => ../utils

module orders

go 1.16
