package imageuploader

import "mime/multipart"

type ImageUploader interface {
	Upload(fileHeader *multipart.FileHeader) (string, error)
	Delete(url string) error
}
