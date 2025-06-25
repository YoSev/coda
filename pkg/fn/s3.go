package fn

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/yosev/coda/internal/utils"
)

type fnS3 struct {
	category FnCategory
}

func (f *fnS3) init(fn *Fn) {
	fn.register("s3.upload", &FnEntry{
		Handler:     f.upload,
		Name:        "Upload to S3",
		Description: "Uploads a file or folder to an S3-compatible bucket",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "endpoint", Description: "The S3 endpoint to use", Mandatory: true},
			{Name: "bucket", Description: "The S3 bucket to use", Mandatory: true},
			{Name: "region", Description: "The S3 region to use", Mandatory: true},
			{Name: "key_id", Description: "The S3 key ID to use", Mandatory: true},
			{Name: "key_secret", Description: "The S3 key secret to use", Mandatory: true},
			{Name: "local_path", Description: "The local path to upload", Mandatory: true},
			{Name: "remote_path", Description: "The remote path in the S3 bucket", Mandatory: false},
			{Name: "remote_prefix", Description: "The remote prefix in the S3 bucket (for recursive upload)", Mandatory: false},
			{Name: "invisible_files", Description: "If true, invisible files will be uploaded", Type: "boolean", Mandatory: false},
		},
	})

	fn.register("s3.download", &FnEntry{
		Handler:     f.download,
		Name:        "Download from S3",
		Description: "Downloads a file or folder from S3 to the local filesystem",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "endpoint", Description: "The S3 endpoint to use", Mandatory: true},
			{Name: "bucket", Description: "The S3 bucket to use", Mandatory: true},
			{Name: "region", Description: "The S3 region to use", Mandatory: true},
			{Name: "key_id", Description: "The S3 key ID to use", Mandatory: true},
			{Name: "key_secret", Description: "The S3 key secret to use", Mandatory: true},
			{Name: "local_path", Description: "The local path to download to", Mandatory: true},
			{Name: "remote_path", Description: "The remote path in the S3 bucket", Mandatory: false},
		},
	})
}

type s3Params struct {
	Endpoint  string `json:"endpoint" yaml:"endpoint"`
	Bucket    string `json:"bucket" yaml:"bucket"`
	Region    string `json:"region" yaml:"region"`
	KeyId     string `json:"key_id" yaml:"key_id"`
	KeySecret string `json:"key_secret" yaml:"key_secret"`

	LocalPath    string `json:"local_path" yaml:"local_path"`                           // Local file or directory
	RemotePath   string `json:"remote_path,omitempty" yaml:"remote_path,omitempty"`     // S3 file key (for single file)
	RemotePrefix string `json:"remote_prefix,omitempty" yaml:"remote_prefix,omitempty"` // S3 folder key (for recursive)

	InvisibleFiles bool `json:"invisible_files,omitempty" yaml:"invisible_files,omitempty"` // Whether to include files starting with a dot (.)
}

// UploadToS3 uploads a single file or folder to an S3-compatible bucket
func (f *fnS3) upload(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *s3Params) (json.RawMessage, error) {
		client, err := buildS3Client(params)
		if err != nil {
			return nil, err
		}

		info, err := os.Stat(params.LocalPath)
		if err != nil {
			return nil, fmt.Errorf("cannot access local path: %w", err)
		}

		prefix := strings.Trim(params.RemotePrefix, "/")
		if prefix != "" {
			prefix += "/"
		}

		var uploaded []string

		if info.IsDir() {
			err = filepath.Walk(params.LocalPath, func(path string, fi os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if fi.IsDir() {
					return nil
				}
				if !params.InvisibleFiles && strings.HasPrefix(fi.Name(), ".") {
					return nil // Skip invisible files
				}

				relPath, err := filepath.Rel(params.LocalPath, path)
				if err != nil {
					return err
				}

				key := filepath.ToSlash(filepath.Join(prefix, relPath))
				if err := uploadFile(client, params.Bucket, path, key); err != nil {
					return err
				}
				uploaded = append(uploaded, key)
				return nil
			})
			if err != nil {
				return nil, err
			}
		} else {
			key := filepath.ToSlash(filepath.Join(prefix, filepath.Base(params.LocalPath)))
			if params.RemotePath != "" {
				key = filepath.ToSlash(params.RemotePath)
			}
			if err := uploadFile(client, params.Bucket, params.LocalPath, key); err != nil {
				return nil, err
			}
			uploaded = append(uploaded, key)
		}

		return json.Marshal(map[string]interface{}{
			"message":  "upload successful",
			"uploaded": uploaded,
		})
	})
}

// DownloadFromS3 downloads a file or folder from S3 to the local filesystem
func (f *fnS3) download(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *s3Params) (json.RawMessage, error) {
		client, err := buildS3Client(params)
		if err != nil {
			return nil, err
		}

		// Normalize remote path (strip leading slash)
		params.RemotePath = strings.TrimPrefix(params.RemotePath, "/")
		remotePath := params.RemotePath

		isFolder := true
		prefix := remotePath

		// Determine if it's a single file, unless remotePath is empty
		if remotePath != "" && !strings.HasSuffix(remotePath, "/") {
			// Try to determine if it is a folder by checking for keys under remotePath + "/"
			testPrefix := remotePath + "/"
			paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
				Bucket:  aws.String(params.Bucket),
				Prefix:  aws.String(testPrefix),
				MaxKeys: aws.Int32(1),
			})

			if paginator.HasMorePages() {
				page, err := paginator.NextPage(context.TODO())
				if err == nil && len(page.Contents) > 0 {
					isFolder = true
					prefix = testPrefix
				} else {
					isFolder = false
				}
			} else {
				isFolder = false
			}
		}

		if !isFolder {
			// === Single file download ===
			target := params.LocalPath
			if err := downloadFile(client, params.Bucket, remotePath, target); err != nil {
				return nil, fmt.Errorf("failed to download file '%s': %w", remotePath, err)
			}
			return json.Marshal(map[string]interface{}{
				"message": "single file download successful",
				"file":    target,
			})
		}

		// === Folder or full bucket download ===
		if prefix != "" && !strings.HasSuffix(prefix, "/") {
			prefix += "/"
		}

		var downloaded []string
		paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
			Bucket: aws.String(params.Bucket),
			Prefix: aws.String(prefix),
		})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(context.TODO())
			if err != nil {
				return nil, fmt.Errorf("list objects: %w", err)
			}

			for _, obj := range page.Contents {
				key := *obj.Key

				// Skip empty "folder marker" objects
				if strings.HasSuffix(key, "/") && obj.Size == aws.Int64(0) {
					continue
				}

				relPath := strings.TrimPrefix(key, prefix)
				if remotePath == "" {
					relPath = key // full path
				}
				if relPath == "" {
					continue
				}

				localPath := filepath.Join(params.LocalPath, relPath)
				if err := downloadFile(client, params.Bucket, key, localPath); err != nil {
					return nil, fmt.Errorf("failed to download key '%s': %w", key, err)
				}
				downloaded = append(downloaded, localPath)
			}
		}

		return json.Marshal(map[string]interface{}{
			"message":    "folder or bucket download successful",
			"downloaded": downloaded,
		})
	})
}

// Helpers

func buildS3Client(params *s3Params) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(params.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			params.KeyId, params.KeySecret, "")),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				if service == s3.ServiceID {
					return aws.Endpoint{
						URL:           params.Endpoint,
						SigningRegion: params.Region,
					}, nil
				}
				return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
			}),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	}), nil
}

func uploadFile(client *s3.Client, bucket, localPath, key string) error {
	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(strings.TrimPrefix(key, "/")),
		Body:   file,
		ACL:    types.ObjectCannedACLPrivate,
	})
	if err != nil {
		return fmt.Errorf("put object: %w", err)
	}
	return nil
}

func downloadFile(client *s3.Client, bucket, key, targetPath string) error {
	resp, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(strings.TrimPrefix(key, "/")),
	})
	if err != nil {
		return fmt.Errorf("get object: %w", err)
	}
	defer resp.Body.Close()

	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	outFile, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}
