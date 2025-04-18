package awshelper

import (
	"fmt"
	"log"
	"mime/multipart"
	"path"
	"tutr-backend/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type AwsInstance struct {
	BucketName string
	AwsRegion  string
	MyError    error
}

func (a *AwsInstance) AwsInit() (*session.Session, error) {

	envVars, errEnv := utils.GetEnvVars()
	if envVars.S3BucketName == "" {
		return nil, errEnv
	}

	a.AwsRegion = envVars.AwsRegion
	a.BucketName = envVars.S3BucketName

	awsSession, err := session.NewSession(
		&aws.Config{
			Region:      aws.String(a.BucketName),
			Credentials: credentials.NewEnvCredentials(),
		},
	)

	if err != nil {
		a.MyError = err
		return nil, a.MyError
	}

	return awsSession, a.MyError
}

// func (a *AwsInstance) PutObjectToAWSS3(file multipart.File, fileHeader *multipart.FileHeader, filePathS3 string) (bool, error) {
// 	envVars, err := utils.GetEnvVars()
// 	if envVars.S3BucketName == "" {
// 		return false, err
// 	}
// 	awsRegion := envVars.AwsRegion
// 	awsBucket := envVars.S3BucketName

// 	sess, err := session.NewSession(&aws.Config{
// 		Region: aws.String(awsRegion),
// 	})

// 	if err != nil {
// 		log.Fatal("Error creating session: ", err)
// 		return false, err
// 	}

// 	_, err = io.ReadAll(file)
// 	if err != nil {
// 		log.Fatal("Error in reading file: ", err)
// 		return false, err
// 	}

// 	fileSize, err := file.Seek(0, io.SeekEnd)
// 	if err != nil {
// 		log.Fatal("Error getting file size: ", err)
// 		return false, err
// 	}

// 	fmt.Printf("File size: %d bytes\n", fileSize)

// 	_, err = file.Seek(0, io.SeekStart)
// 	if err != nil {
// 		log.Fatal("Error resetting file pointer: ", err)
// 		return false, err
// 	}

// 	defer file.Close()

// 	svc := s3.New(sess)

// 	input := s3.PutObjectInput{
// 		Bucket:      aws.String(awsBucket),
// 		Key:         aws.String(path.Join(filePathS3, fileHeader.Filename)),
// 		Body:        file,
// 		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
// 	}

// 	_, err = svc.PutObject(&input)
// 	if err != nil {
// 		log.Fatal("Error in Uploading file: ", err)
// 		return false, err
// 	}

// 	fmt.Printf("File Uploaded Successfully!")
// 	return true, err
// }

func (a *AwsInstance) PutObjectToAWSS3(file multipart.File, fileName string, filePathS3 string) (bool, error) {
	// ✅ Load environment variables
	envVars, err := utils.GetEnvVars()
	if envVars.S3BucketName == "" {
		return false, err
	}
	awsRegion := envVars.AwsRegion
	awsBucket := envVars.S3BucketName

	// ✅ Create AWS S3 session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})

	if err != nil {
		log.Println("Error creating session:", err)
		return false, err
	}

	svc := s3.New(sess)

	// ✅ Prepare input for S3 upload
	input := &s3.PutObjectInput{
		Bucket:      aws.String(awsBucket),
		Key:         aws.String(path.Join(filePathS3, fileName)),
		Body:        file,
		ContentType: aws.String("application/octet-stream"), // ✅ Set generic content type
	}

	// ✅ Upload file to S3
	_, err = svc.PutObject(input)
	if err != nil {
		log.Println("Error uploading file:", err)
		return false, err
	}

	fmt.Printf("File %s uploaded successfully!\n", fileName)
	return true, nil
}

func GetAllFilesFromBucket() *s3.S3 {
	envVars, _ := utils.GetEnvVars()
	if envVars.S3BucketName == "" {
		return nil
	}
	awsRegion := envVars.AwsRegion
	// awsBucket := os.Getenv("BUCKETNAME")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		log.Fatal("Error creating session: ", err)
	}

	if err != nil {
		log.Fatal("Error in reading file: ", err)
	}

	svc := s3.New(sess)
	return svc
}
