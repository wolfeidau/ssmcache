// Package mocks provides mock implementations for AWS services
package mocks

//go:generate go run github.com/golang/mock/mockgen -destination awsmock/ssm.go -package awsmock github.com/aws/aws-sdk-go/service/ssm/ssmiface SSMAPI
