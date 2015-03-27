package aws

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/hashicorp/aws-sdk-go/aws"
	"github.com/hashicorp/aws-sdk-go/gen/iam"
)

func resourceAwsInstanceProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsInstanceProfileCreate,
		Read:   resourceAwsInstanceProfileRead,
		Update: resourceAwsInstanceProfileUpdate,
		Delete: resourceAwsInstanceProfileDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceAwsInstanceProfileCreate(d *schema.ResourceData, meta interface{}) error {
	iamconn := meta.(*AWSClient).iamconn
	//awsRegion := meta.(*AWSClient).region

	// Get the name
	name := d.Get("name").(string)

	log.Printf("[DEBUG] InstanceProfile create: %s", name)

	req := &iam.CreateInstanceProfileRequest{
		InstanceProfileName: aws.String(name),
		Path:    aws.String("/"),
	}

	_, err := iamconn.CreateInstanceProfile(req)
	if err != nil {
		return fmt.Errorf("Error creating InstanceProfile: %s", err)
	}

	// Assign the InstanceProfileName as the resource ID
	d.SetId(name)

	//return resourceAwsS3BucketUpdate(d, meta)
        return nil;
}

func resourceAwsInstanceProfileUpdate(d *schema.ResourceData, meta interface{}) error {
        return nil;
/*
	s3conn := meta.(*AWSClient).s3conn
	if err := setTagsS3(s3conn, d); err != nil {
		return err
	}
	return resourceAwsS3BucketRead(d, meta)
*/
}

func resourceAwsInstanceProfileRead(d *schema.ResourceData, meta interface{}) error {
        return nil;
/*
	s3conn := meta.(*AWSClient).s3conn

	err := s3conn.HeadBucket(&s3.HeadBucketRequest{
		Bucket: aws.String(d.Id()),
	})
	if err != nil {
		return err
	}

	resp, err := s3conn.GetBucketTagging(&s3.GetBucketTaggingRequest{
		Bucket: aws.String(d.Id()),
	})
	if err != nil {
		return err
	}

	if err := d.Set("tags", tagsToMapS3(resp.TagSet)); err != nil {
		return err
	}
*/
	return nil
}

func resourceAwsInstanceProfileDelete(d *schema.ResourceData, meta interface{}) error {
        return nil;
/*
	s3conn := meta.(*AWSClient).s3conn

	log.Printf("[DEBUG] S3 Delete Bucket: %s", d.Id())
	err := s3conn.DeleteBucket(&s3.DeleteBucketRequest{
		Bucket: aws.String(d.Id()),
	})
	if err != nil {
		return err
	}
	return nil
*/
}
