// Package mocks provides mock implementations for AWS services
package mocks

//go:generate env GOBIN=$PWD/.bin GO111MODULE=on go install github.com/golang/mock/mockgen
//go:generate $PWD/.bin/mockgen -destination awsmock/ssm.go -package awsmock github.com/aws/aws-sdk-go/service/ssm/ssmiface SSMAPI
