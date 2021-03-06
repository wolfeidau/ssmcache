package ssmcache

import (
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/pkg/errors"
)

var defaultExpiry = 30 * time.Second

// SetDefaultExpiry update the default expiry for all cached parameters
//
// Note this will update expires value on the next refresh of entries.
func SetDefaultExpiry(expires time.Duration) {
	defaultExpiry = expires
}

// Entry an SSM entry in the cache
type Entry struct {
	value       string
	expires     time.Time
	lastVersion int64
}

// Cache SSM cache which provides read access to parameters
type Cache interface {
	GetKey(string, bool) (string, error)
	PutKey(string, string, bool) error
}

type cache struct {
	ssm       sync.Mutex
	ssmValues map[string]*Entry
	ssmSvc    ssmiface.SSMAPI
}

// New new SSM cache
func New(cfgs ...*aws.Config) Cache {
	sess := session.Must(session.NewSession(cfgs...))
	return &cache{
		ssmSvc:    ssm.New(sess),
		ssmValues: make(map[string]*Entry),
	}
}

// GetKey retrieve a parameter from SSM and cache it.
func (ssc *cache) GetKey(key string, encrypted bool) (string, error) {

	ssc.ssm.Lock()
	defer ssc.ssm.Unlock()

	ent, ok := ssc.ssmValues[key]
	if !ok {
		// record is missing
		return ssc.updateParam(key, encrypted)
	}

	if time.Now().After(ent.expires) {
		// we have expired and need to refresh
		// log.Println("expired cache refreshing value")

		return ssc.updateParam(key, encrypted)
	}

	// return the value
	return ent.value, nil
}

func (ssc *cache) PutKey(key string, val string, encrypted bool) error {
	putParams := &ssm.PutParameterInput{
		Name:      aws.String(key),
		Overwrite: aws.Bool(true),
		Value:     aws.String(val),
	}

	if encrypted {
		putParams.Type = aws.String(ssm.ParameterTypeSecureString)
	}

	_, err := ssc.ssmSvc.PutParameter(putParams)
	if err != nil {
		return errors.Wrapf(err, "failed to store key %s from ssm", key)
	}

	_, err = ssc.updateParam(key, encrypted)
	if err != nil {
		return errors.Wrapf(err, "failed to retrieve key %s from ssm", key)
	}

	return nil
}

func (ssc *cache) updateParam(key string, encrypted bool) (string, error) {

	// log.Println("updating key from ssm:", key)

	if _, ok := ssc.ssmValues[key]; ok {
		// check if the parameter has been modified recently, if not we can just return the existing
		// value and avoid another decryption call which costs money.
		modified, err := ssc.parameterModified(key)
		if err != nil {
			return "", errors.Wrapf(err, "failed to get parameter key %s in ssm", key)
		}

		if !modified {
			// log.Println("skipping refresh returning unmodified value")
			ssc.ssmValues[key].expires = time.Now().Add(defaultExpiry)
			return ssc.ssmValues[key].value, nil
		}
	}

	resp, err := ssc.ssmSvc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(encrypted),
	})
	if err != nil {
		return "", errors.Wrapf(err, "failed to retrieve key %s from ssm", key)
	}

	ssc.ssmValues[key] = &Entry{
		value:       aws.StringValue(resp.Parameter.Value),
		expires:     time.Now().Add(defaultExpiry), // reset the expiry
		lastVersion: aws.Int64Value(resp.Parameter.Version),
	}

	// log.Println("key value refreshed from ssm at:", time.Now())

	return aws.StringValue(resp.Parameter.Value), nil
}

func (ssc *cache) parameterModified(key string) (bool, error) {
	desc, err := ssc.ssmSvc.GetParameter(&ssm.GetParameterInput{
		Name: aws.String(key),
	})
	if err != nil {
		return false, err
	}

	// log.Printf("check last modified param: %t cache: %d", aws.TimeValue(desc.Parameter.LastModifiedDate), ssc.ssmValues[key].lastVersion)

	// check if the parameter has been modified recently, if not we can just return the existing
	// value and avoid another decryption call which costs money.
	return aws.Int64Value(desc.Parameter.Version) > ssc.ssmValues[key].lastVersion, nil
}
