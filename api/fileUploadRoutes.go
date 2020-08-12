package api

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//Do we need to utilize a seperate struct for this?
type experienceImage struct {
	Filename     string `json:"filename,omitempty"`
	Data         []byte `json:"data,omitempty"`
	ExperienceID string `json:"experienceId,omitempty"`
	FileType     string `json:"fileType,omitempty"`
}

func okContentType(contentType string) (string, bool) {
	if contentType == "image/png" || contentType == "image/jpeg" || contentType == "image/jpg" {
		contentExtension := strings.Split(contentType, "image/")
		return contentExtension[1], true
	}
	return "", false
}

/*Add error responses back to client*/
func imageUpload(w http.ResponseWriter, r *http.Request) {
	var imageResponse experienceImage

	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		errorResponse := fmt.Sprintf("Error retrieving file: %s", err)
		log.Println(err)
		imageResponse = experienceImage{
			Filename: errorResponse,
		}
		sendFailResponse(w, imageResponse)
		return
	}
	defer file.Close()

	experience := r.FormValue("experienceId")

	//This may be needed for a specific file structure
	log.Printf("Uploaded File: %+v\n", handler.Filename)
	log.Printf("File Size: %+v\n", handler.Size)
	log.Printf("MIME Header: %+v\n", handler.Header)

	//Check the file type before proceeding
	fileExtension, isValid := okContentType(handler.Header["Content-Type"][0])
	if !isValid {
		log.Println("Invalid Content type")
		imageResponse = experienceImage{
			FileType: "Invalid Content type",
		}
		sendFailResponse(w, imageResponse)
		return
	}

	//Read contents of file into a byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		errorResponse := fmt.Sprintf("Unable to read contents for file: %s", err)
		log.Println(errorResponse)
		imageResponse = experienceImage{
			Filename: errorResponse,
		}
		sendFailResponse(w, imageResponse)
		return
	}

	image := experienceImage{
		FileType:     fileExtension,
		Data:         fileBytes,
		ExperienceID: experience,
	}

	err = proccessImage(&image)
	if err != nil {
		errorResponse := fmt.Sprintf("Error processing image: %s", err)
		log.Println(errorResponse)
		imageResponse = experienceImage{
			Filename: errorResponse,
		}
		sendFailResponse(w, imageResponse)
		return
	}
	log.Println("Successfully uploaded file")

	imageResponse = experienceImage{
		Filename: image.Filename,
	}
	sendSuccessResponse(w, imageResponse)
	return

}

//do we want to use the tempfile API, simplifies things like filenaming
//However, may not be intended use case
func proccessImage(image *experienceImage) error {
	tempImageTemplate := fmt.Sprintf("upload-*.%s", image.FileType)

	//Create a temporary file with our images directory
	tempFile, err := ioutil.TempFile("images", tempImageTemplate)
	if err != nil {
		errorResponse := fmt.Sprintf("Issue creating file: %s", err)
		return errors.New(errorResponse)
		//send 500 response
	}
	defer tempFile.Close()

	//Gets the filename without the directory path
	fileName := strings.Split(tempFile.Name(), "images/")
	image.Filename = fileName[1]

	//File is written
	tempFile.Write(image.Data)

	err = updateExperienceImageInSlice(image.ExperienceID, image.Filename)
	if err != nil {
		return err
	}

	return nil
}
