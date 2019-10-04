package security

import (
	"bytes"
	"fmt"
	"github.com/go-delve/delve/pkg/config"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
)

// Simple Get
func CheckEmailRegistered(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "please provide the email in the query ex: ?email=123..."})
		return
	}
}

// POST

// UserValidatePhone ...
func UserValidatePhone(c *gin.Context) {

	// get the user id from jwt
	jwtClaims := jwt.ExtractClaims(c)
	id, ok := jwtClaims["id"].(string)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid request"})
		return
	}

	var payload models.ValidatePhone

	err := c.BindJSON(&payload)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = database.CheckPhoneExists(strconv.Itoa(payload.CountryCode) + payload.PhoneNumber)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

}

// upload file
func UploadFileExample(c *gin.Context) {

	// get the user id from jwt
	jwtClaims := jwt.ExtractClaims(c)
	id, ok := jwtClaims["id"].(string)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid request"})
		return
	}

	config, ok := c.MustGet("configuration").(*config.Config)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get the configuration"})
		return
	}

	user, err := database.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}


	//checks if the user doesn't try to upload too many documents
	totalDocs, err := database.GetUploadsNumberForUser(user.ID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "system error: " + err.Error()})
		return
	}
	fmt.Println("total docs this user uploaded", totalDocs)
	if totalDocs > 10 {
		c.JSON(http.StatusForbidden, gin.H{"error": "you reached the maximum number of documents you can upload."})
		return
	}

	newID := uuid.NewV4().String()[0:8]
	uploadedFile, header, err := c.Request.FormFile("upload")

	// fast file verifications
	err = verifyHeader(header)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// upload them to S3 ...............
	docPath := "kyc_" + user.ID + "_" + newID + filepath.Ext(header.Filename)

	sess, err := session.NewSession(&aws.Config{Region: aws.String("ca-central-1"),
		S3ForcePathStyle: aws.Bool(true), LogLevel: aws.LogLevel(aws.LogDebug),
		Credentials: credentials.NewStaticCredentials(config.AWSS3AcccessKeyID, config.AWSS3SecretAccessKey, "")})
	if err != nil {
		log.Error(err)
		return
	}
	s3client := s3.New(sess)

	buffer := make([]byte, header.Size)

	// read file content to buffer
	_, err = uploadedFile.Read(buffer)
	if err != nil {
		log.Error(err)
		return
	}

	fileBytes := bytes.NewReader(buffer) // convert to io.ReadSeeker type
	fileType := http.DetectContentType(buffer)

	params := &s3.PutObjectInput{
		Bucket:      aws.String("xxxxxxxxxxxx"), // required
		Key:         aws.String(docPath),      // required
		ACL:         aws.String("public-read"),
		Body:        fileBytes,
		ContentType: aws.String(fileType),
		Metadata: map[string]*string{
			"UserID": aws.String(user.ID), //required
		},
	}

	result, err := s3client.PutObject(params)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			log.Error(awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
			if reqErr, ok := err.(awserr.RequestFailure); ok {
				log.Error(reqErr.Code(), reqErr.Message(), reqErr.StatusCode(), reqErr.RequestID())
			}
		} else {
			log.Error(err)
		}
		return
	}

	fmt.Println("Uploaded ok to S3: " + awsutil.StringValue(result))

	// save it in the database
	err = database.CreateUpload(user.ID, docPath)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"error": "system error. please contact admin"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "uploaded ok"})
}

func verifyHeader(header *multipart.FileHeader) error {
	if header == nil {
		return fmt.Errorf("header is nil")
	}

	contentType := header.Header.Get("Content-Type")
	var re = regexp.MustCompile(`image/png|image/jpeg|image/jpg|application/pdf`)
	if !re.MatchString(contentType) {
		return fmt.Errorf("invalid image or pdf")
	}

	if header.Size > 5245880 {
		return fmt.Errorf("file is too large. please upload a 5mb maximum image")
	}

	return nil
}