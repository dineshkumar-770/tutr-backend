package addnotes

import (
	"encoding/json"
	"fmt"
	awshelper "main/aws_helper"
	"main/constants"
	"main/database"
	"main/helpers"
	"main/middlewares"
	"main/models/notes"
	u "main/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

var myAwsInstance = awshelper.AwsInstance{}

func GetAllNotesOfTheGroup(w http.ResponseWriter, r *http.Request) {
	resp := u.ResponseStr{
		Status:     "failed",
		Message:    "something went wrong",
		MyResponse: nil,
	}

	tokenString := r.Header.Get(constants.Authrization)
	if tokenString == "" {
		u.SendResponseWithUnauthorizedError(w)
		return
	}

	tokenString = tokenString[len("Bearer "):]
	claims, err := middlewares.VerifyToken(tokenString)
	if err != nil {
		resp.Message = "unable to verify your token"
		fmt.Println(err)
		u.SendResponseWithUnauthorizedError(w)
		return
	}

	teacherId := claims["user_id"]
	groupId := r.URL.Query().Get("group_id")

	fmt.Println(teacherId)

	if groupId == "" {
		resp.Message = "Missing values are not allowed"
		u.SendResponseWithMissingValues(w)
		return
	}

	db := database.GetDBInstance()

	rows, err := db.Query(`SELECT notes_id, notes_title, notes_desctription, class, topic, subject, file_url, uploaded_at, visibility, is_editable FROM tutrdevdb.teacher_notes WHERE group_id=?`, groupId)

	if err != nil {
		resp.Message = "Database error while fetching notes"
		u.SendResponseWithServerError(w, resp)
		return
	}
	defer rows.Close()

	// var notesList []map[string]interface{}
	envVars, _ := u.GetEnvVars()

	s3Svc := awshelper.GetAllFilesFromBucket()
	if s3Svc == nil {
		resp.Message = "S3 bucket configuration error"
		u.SendResponseWithServerError(w, resp)
		return
	}

	var listOfGroupNotes []notes.GroupNotes

	for rows.Next() {
		var groupNotes notes.GroupNotes

		err := rows.Scan(&groupNotes.NotesId, &groupNotes.NotesTitle, &groupNotes.NotesDescription, &groupNotes.ClassName, &groupNotes.NotesTopic, &groupNotes.NotesSubject, &groupNotes.FileNames, &groupNotes.UploadedAt, &groupNotes.NotesVisiblity, &groupNotes.IsEditableInt)
		if err != nil {
			continue
		}

		fileURLList := strings.Split(groupNotes.FileNames, ",")
		signedURLs := []string{}

		if groupNotes.IsEditableInt == 1 {
			groupNotes.IsEditable = true
		} else {
			groupNotes.IsEditable = false
		}

		for _, filePath := range fileURLList {
			req, _ := s3Svc.GetObjectRequest(&s3.GetObjectInput{
				Bucket: aws.String(envVars.S3BucketName),
				Key:    aws.String(filePath),
			})
			signedURL, err := req.Presign(168 * time.Hour)
			getFileMetadata(signedURL)
			if err == nil {
				signedURLs = append(signedURLs, signedURL)
			}
		}
		groupNotes.FileURLsList = signedURLs

		listOfGroupNotes = append(listOfGroupNotes, groupNotes)
	}

	resp.Status = "success"
	resp.Message = "Notes fetched successfully"
	resp.MyResponse = listOfGroupNotes
	u.SendResponseWithOK(w, resp)
}
func getFileMetadata(signedURL string) {
	req, err := http.NewRequest("HEAD", signedURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching metadata:", err)
		return
	}
	defer resp.Body.Close()

	// Displaying the headers
	fmt.Println("File Metadata:")
	fmt.Println("Content-Type:", resp.Header.Get("Content-Type"))     // MIME Type
	fmt.Println("Content-Length:", resp.Header.Get("Content-Length")) // File Size in bytes
	fmt.Println("Last-Modified:", resp.Header.Get("Last-Modified"))   // Last modified date
	fmt.Println("ETag:", resp.Header.Get("ETag"))                     // Unique identifier of the file

	// Additional headers for debugging
	fmt.Println("Cache-Control:", resp.Header.Get("Cache-Control"))
	fmt.Println("Content-Disposition:", resp.Header.Get("Content-Disposition"))
	fmt.Println("X-Amz-Meta-SomeCustomMetadata:", resp.Header.Get("X-Amz-Meta-SomeCustomMetadata"))
}

func AddTeacherClassNotesToStorage2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "multipart/form-data")

	resp := u.ResponseStr{
		Status:     "failed",
		Message:    "something went wrong",
		MyResponse: nil,
	}

	tokenString := r.Header.Get(constants.Authrization)
	if tokenString == "" {
		u.SendResponseWithUnauthorizedError(w)
		return
	}

	tokenString = tokenString[len("Bearer "):]
	claims, err := middlewares.VerifyToken(tokenString)
	if err != nil {
		resp.Message = "unable to verify your token"
		u.SendResponseWithUnauthorizedError(w)
		return
	}

	userId := claims["user_id"].(string)
	notesTitle := r.FormValue("notes_title")
	notesDescription := r.FormValue("notes_desc")
	class := r.FormValue("class")
	topic := r.FormValue("notes_topic")
	subject := r.FormValue("subject")
	visibility := r.FormValue("visibility")
	isEditable := r.FormValue("is_editable")
	groupId := r.FormValue("group_id")
	isEditableBool, _ := strconv.ParseBool(isEditable)

	if userId == "" || notesTitle == "" || notesDescription == "" || class == "" || topic == "" || subject == "" || visibility == "" || isEditable == "" || groupId == "" {
		resp.Message = "Missing fields are not allowed. Kindly fill all the details."
		u.SendResponseWithMissingValues(w)
		return
	}

	envVars, _ := u.GetEnvVars()
	if envVars.S3BucketName == "" {
		resp.Message = "Something went wrong with the storage to save your notes! Please try again later."
		u.SendResponseWithServerError(w, resp)
		return
	}

	s3BucketFolderPath := envVars.S3NotesFolder

	reqErr := r.ParseMultipartForm(32 << 20)
	if reqErr != nil {
		resp.Message = "Notes file needs to be attached to upload. Please attach a file."
		u.SendResponseWithStatusBadRequest(w, resp)
		return
	}

	files := r.MultipartForm.File["notes"]
	if len(files) == 0 {
		resp.Message = "No files found. Please attach at least one file."
		u.SendResponseWithStatusBadRequest(w, resp)
		return
	}

	db := database.GetDBInstance()
	timestamp := time.Now().Unix()

	var existingFileURLs string
	_ = db.QueryRow(`SELECT file_url FROM teacher_notes WHERE notes_title=? AND notes_desctription=? AND class=? AND topic=? AND subject=? AND teacher_id=? AND visibility=? AND is_editable=? AND group_id=?`,
		notesTitle, notesDescription, class, topic, subject, userId, visibility, isEditableBool, groupId).Scan(&existingFileURLs)

	uploadedFileURLs := []string{}

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		// extension := filepath.Ext(fileHeader.Filename)
		//TODO: rename this file as <timestamp_filename>
		// newFileName := fmt.Sprintf("%s_%d%s", fileHeader.Filename, timestamp, extension)
		newFileName := fmt.Sprintf("%d_%s", timestamp, fileHeader.Filename)

		status, err := myAwsInstance.PutObjectToAWSS3(file, newFileName, s3BucketFolderPath)
		if !status {
			resp.Message = "Unknown Error occured while uploading your notes."
			fmt.Println(err)
			u.SendResponseWithServerError(w, resp)
			return
		}

		fileURL := fmt.Sprintf("%s%s", s3BucketFolderPath, newFileName)
		uploadedFileURLs = append(uploadedFileURLs, fileURL)
	}

	if len(uploadedFileURLs) == 0 {
		resp.Message = "None of the files could be uploaded."
		u.SendResponseWithServerError(w, resp)
		return
	}

	finalFileURLs := strings.Join(uploadedFileURLs, ",")
	if existingFileURLs != "" {
		finalFileURLs = existingFileURLs + "," + finalFileURLs
	}

	if existingFileURLs == "" {
		notesID := helpers.GenerateUserID("notes")
		_, dbErr := db.Exec(`INSERT INTO teacher_notes (notes_id, notes_title, notes_desctription, class, topic, subject, teacher_id, uploaded_at, file_url, visibility, is_editable, group_id) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`,
			notesID, notesTitle, notesDescription, class, topic, subject, userId, timestamp, finalFileURLs, visibility, isEditableBool, groupId)

		if dbErr != nil {
			resp.Message = "Error inserting new record."
			u.SendResponseWithServerError(w, resp)
			return
		}
	} else {
		_, dbErr := db.Exec(`UPDATE teacher_notes SET file_url=?, uploaded_at=? WHERE notes_title=? AND notes_desctription=? AND class=? AND topic=? AND subject=? AND teacher_id=? AND visibility=? AND is_editable=? AND group_id=?`,
			finalFileURLs, timestamp, notesTitle, notesDescription, class, topic, subject, userId, visibility, isEditableBool, groupId)

		if dbErr != nil {
			resp.Message = "Error updating record."
			u.SendResponseWithServerError(w, resp)
			return
		}
	}

	resp.Status = "success"
	resp.Message = fmt.Sprintf("%d files uploaded successfully", len(uploadedFileURLs))
	resp.MyResponse = uploadedFileURLs
	u.SendResponseWithOK(w, resp)
}

func DeleteTeacherNotes(w http.ResponseWriter, r *http.Request) {
	resp := u.ResponseStr{
		Status:     "failed",
		Message:    "something went wrong",
		MyResponse: nil,
	}

	tokenString := r.Header.Get(constants.Authrization)
	if tokenString == "" {
		u.SendResponseWithUnauthorizedError(w)
		return
	}

	tokenString = tokenString[len("Bearer "):]
	claims, err := middlewares.VerifyToken(tokenString)
	if err != nil {
		resp.Message = "unable to verify your token"
		u.SendResponseWithUnauthorizedError(w)
		return
	}

	userId := claims["user_id"].(string)

	if userId == "" {
		resp.Message = "unable to verify your token"
		u.SendResponseWithUnauthorizedError(w)
		return
	}

	var deleteNotes notes.DeletedGroupNotes

	deocodeErr := json.NewDecoder(r.Body).Decode(&deleteNotes)
	if deocodeErr != nil {
		resp.Message = "Some field are missing. kindly provide the proper fields"
		resp.Status = constants.FAILED
		u.SendResponseWithMissingValues(w)
	}

	deletequery := `
	DELETE FROM tutrdevdb.teacher_notes WHERE teacher_notes.notes_id = (?);
	`

	addquery := `
		INSERT INTO tutrdevdb.notes_trash (trash_id,deleted_notes_id,group_id,deleted_at,notes_title,reason,notes_description,file_urls) VALUES (?,?,?,?,?,?,?,?)
	`

	db := database.GetDBInstance()
	_, dbErr := db.Exec(deletequery, deleteNotes.DeletdNotesID)

	if dbErr != nil {
		resp.Message = "Error in deleting the notes record."
		u.SendResponseWithServerError(w, resp)
		return
	}

	trashId := helpers.GenerateUserID("trash")
	deletedAt := time.Now().Unix()

	_, err1 := db.Exec(addquery, deleteNotes, trashId, deleteNotes.DeletdNotesID, deleteNotes.GroupID, deletedAt, deleteNotes.NotesTitle, deleteNotes.Reason, deleteNotes.NotesDescription, deleteNotes.FileUrls)

	if err1 != nil {
		resp.Message = "Error in deleting the notes"
		u.SendResponseWithServerError(w, resp)
		return
	}

	resp.Status = constants.SUCCESS
	resp.Message = "Notes deleted successfully"
	u.SendResponseWithOK(w, resp)
}

// func AddTeacherClassNotesToStorage(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "multipart/form-data")

// 	resp := u.ResponseStr{
// 		Status:     "failed",
// 		Message:    "something went wrong",
// 		MyResponse: nil,
// 	}

// 	tokenString := r.Header.Get(constants.Authrization)
// 	if tokenString == "" {
// 		u.SendResponseWithUnauthorizedError(w)
// 		return
// 	}

// 	tokenString = tokenString[len("Bearer "):]
// 	claims, err := middlewares.VerifyToken(tokenString)
// 	if err != nil {
// 		resp.Message = "unable to verify your token"
// 		fmt.Println(err)
// 		u.SendResponseWithUnauthorizedError(w)
// 		return
// 	}

// 	userId := claims["user_id"]
// 	fmt.Println(userId)

// 	envVars, errEnv := u.GetEnvVars()
// 	if envVars.S3BucketName == "" {
// 		resp.Status = "failed"
// 		fmt.Println("no environment found error: ", errEnv.Error())
// 		resp.Message = "Something wnet wrong with the storage to save your notes! please try again later."
// 		u.SendResponseWithServerError(w, resp)
// 		return
// 	}

// 	s3BucketFolderPath := envVars.S3NotesFolder
// 	notesFileName := r.FormValue("notes_file_name")
// 	reqErr := r.ParseMultipartForm(32 << 20)
// 	if reqErr != nil {
// 		resp.Status = "failed"
// 		resp.Message = "Notes file needs to be attach to upload in your group. Please attach a file"
// 		u.SendResponseWithStatusBadRequest(w, resp)
// 		return
// 	}

// 	file, handler, err := r.FormFile("notes")
// 	if err != nil {
// 		resp.Status = "Failed"
// 		resp.Message = "some error occured while parsing your uploaded file, or may be the file format is unsupported. Please attach the correct format of the file."
// 		fmt.Println("Error in forming file OS : ", err.Error())
// 		u.SendResponseWithStatusBadRequest(w, resp)
// 		return
// 	}
// 	defer file.Close()

// 	timestamp := time.Now().Unix()
// 	extension := filepath.Ext(handler.Filename)
// 	newFileName := fmt.Sprintf("%s_%d%s", notesFileName, timestamp, extension)
// 	handler.Filename = newFileName

// 	status, err := myAwsInstance.PutObjectToAWSS3(file, handler, s3BucketFolderPath)

// 	if !status {
// 		resp.Status = "Failed"
// 		resp.Message = "cannot save your file into our system right now!"
// 		fmt.Println("Unable to Save the file in Cloud : ", err.Error())
// 		u.SendResponseWithServerError(w, resp)
// 		return
// 	}
// 	resp.Status = "success"
// 	resp.Message = "File saved successfully"
// 	u.SendResponseWithOK(w, resp)
// }
