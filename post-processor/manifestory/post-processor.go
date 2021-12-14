//go:generate packer-sdc mapstructure-to-hcl2 -type Config

package manifestory

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/common"
	awscommon "github.com/hashicorp/packer-plugin-amazon/builder/common"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/service/s3/s3manager"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	awscommon.AccessConfig `mapstructure:",squash"`
	ctx                 interpolate.Context

	// Variables specific to this post processor
	S3Bucket        string            	`mapstructure:"s3_bucket_name"`
	S3Key           string            	`mapstructure:"s3_key_name"`
	MockOption      string 				`mapstructure:"mock"`

}

type PostProcessor struct {
	config Config
}

func (p *PostProcessor) ConfigSpec() hcldec.ObjectSpec { return p.config.FlatMapstructure().HCL2Spec() }

func (p *PostProcessor) Configure(raws ...interface{}) error {
	err := config.Decode(&p.config, &config.DecodeOpts{
		PluginType:         "packer.post-processor.manifestory",
		Interpolate:        true,
		InterpolateContext: &p.config.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{},
		},
	}, raws...)
	if err != nil {
		return err
	}
	
	if p.config.S3Key == "" {
		p.config.S3Key = "packer-import"
	}

	p.config.S3Key = "hudson-integration-test"

	packersdk.LogSecretFilter.Set(p.config.AccessKey, p.config.SecretKey, p.config.Token)
	return nil
}

func (p *PostProcessor) PostProcess(ctx context.Context, ui packersdk.Ui, source packersdk.Artifact) (packersdk.Artifact, bool, bool, error) {
	ui.Say(fmt.Sprintf("post-processor mock: %s", p.config.PackerBuilderType))
	session, err := p.config.Session()
	if err != nil {
		return nil, false, false, err
	}
	
	s3conn := s3.New(session)
	i := 0
	err = s3conn.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: aws.String("hudson-integration-dev"),
		Prefix: aws.String("Timesheet/"),
	}, func(p *s3.ListObjectsOutput, last bool) (shouldContinue bool) {
		ui.Say(fmt.Sprintf("Page, %s", i))
		
		i++

		for _, obj := range p.Contents {
			ui.Say(fmt.Sprintf("Object: %s", *obj.Key))
		}
		return true
	})
	if err != nil {
		return nil, false, false, fmt.Errorf("Failed to list s3: %s", err)
	}
	return source, true, true, nil
}
