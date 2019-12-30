# ssmcache

This is a simple cache library which only retrieves params from [SSM parameter store](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-parameter-store.html) when the cache time expires or when a version updated has occurred.

Note the versions checks are done without decryption to ensure we limit the calls to [KMS](https://aws.amazon.com/kms/) to keep the costs down.

[![GitHub Actions status](https://github.com/wolfeidau/ssmcache/workflows/Go/badge.svg?branch=master)](https://github.com/wolfeidau/ssmcache/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/wolfeidau/ssmcache)](https://goreportcard.com/report/github.com/wolfeidau/ssmcache)
[![Documentation](https://godoc.org/github.com/wolfeidau/ssmcache?status.svg)](https://godoc.org/github.com/wolfeidau/ssmcache)

# License

This code is released under the Apache 2.0 license.
