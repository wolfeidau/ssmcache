package ssmcache

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/wolfeidau/ssmcache/mocks/awsmock"
)

func TestGetKey(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gpo := &ssm.GetParameterOutput{
		Parameter: &ssm.Parameter{
			Name:    aws.String("testtest"),
			Value:   aws.String("sup"),
			Version: aws.Int64(1),
		},
	}

	ssmMock := awsmock.NewMockSSMAPI(ctrl)

	ssmMock.EXPECT().GetParameter(&ssm.GetParameterInput{
		Name:           aws.String("testtest"),
		WithDecryption: aws.Bool(true),
	}).Return(gpo, nil).MaxTimes(1)

	cache := &cache{
		ssmSvc:    ssmMock,
		ssmValues: make(map[string]*Entry),
	}

	SetDefaultExpiry(1 * time.Second)

	val, err := cache.GetKey("testtest", true)
	require.Nil(t, err)
	require.Equal(t, "sup", val)

	time.Sleep(1 * time.Second)

	dpo := &ssm.GetParameterOutput{
		Parameter: &ssm.Parameter{
			Version: aws.Int64(1),
		},
	}

	ssmMock.EXPECT().GetParameter(&ssm.GetParameterInput{
		Name: aws.String("testtest"),
	}).Return(dpo, nil).MaxTimes(2)

	val, err = cache.GetKey("testtest", true)
	require.Nil(t, err)
	require.Equal(t, "sup", val)

	time.Sleep(1 * time.Second)

	// simulate an update of key where a subsequent change ot the parameter will
	// trigger retrieval from SSM
	gpo.Parameter.Version = aws.Int64(2)
	val, err = cache.GetKey("testtest", true)
	require.Nil(t, err)
	require.Equal(t, "sup", val)
}

func TestPutKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ssmMock := awsmock.NewMockSSMAPI(ctrl)

	ppo := &ssm.PutParameterOutput{
		Version: aws.Int64(1),
	}
	ssmMock.EXPECT().PutParameter(&ssm.PutParameterInput{
		Name:      aws.String("testtest"),
		Overwrite: aws.Bool(true),
		Type:      aws.String("SecureString"),
		Value:     aws.String("sup"),
	}).Return(ppo, nil)
	cache := &cache{
		ssmSvc:    ssmMock,
		ssmValues: make(map[string]*Entry),
	}
	gpo := &ssm.GetParameterOutput{
		Parameter: &ssm.Parameter{
			Name:  aws.String("testtest"),
			Value: aws.String("sup"),
		},
	}

	ssmMock.EXPECT().GetParameter(&ssm.GetParameterInput{
		Name:           aws.String("testtest"),
		WithDecryption: aws.Bool(true),
	}).Return(gpo, nil)

	err := cache.PutKey("testtest", "sup", true)
	require.Nil(t, err)
}
