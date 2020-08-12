package api

import (
	"io/ioutil"
	"log"
	"net/http"
)

//Do we need to utilize a seperate struct for this?
type image struct {
	Filename    string
	ContentType string
	Data        []byte
	Size        int
}

func okContentType(contentType string) bool {
	return contentType == "image/png" || contentType == "image/jpeg" || contentType == "image/jpg"
}

/*Add error responses back to client*/
func imageUpload(w http.ResponseWriter, r *http.Request) {
	log.Println("Upload Image endpoint hit")

	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		log.Println("Error retrieving file")
		log.Println(err)
		return
		//send request back
	}
	defer file.Close()

	//This may be needed for a specific file structure
	log.Printf("Uploaded File: %+v\n", handler.Filename)
	log.Printf("File Size: %+v\n", handler.Size)
	log.Printf("MIME Header: %+v\n", handler.Header)

	//Check the file type before proceeding
	if !okContentType(handler.Header["Content-Type"][0]) {
		log.Println("Invalid Content type")
		return
	}

	//Create a temporary file with our images directory
	tempFile, err := ioutil.TempFile("images", "upload-*.png")
	if err != nil {
		log.Println(err)
		return
		//send 500 response back?
	}
	defer tempFile.Close()

	//Read contents of file into a byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		return
		//send response?
	}
	//File is written
	tempFile.Write(fileBytes)
	log.Println("Successfully uploaded file")

	//need to be able to create a link to our experience
	//Should be determined by our other parameter
}
