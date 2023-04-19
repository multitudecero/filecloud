package models

type UploadStep string

const (
	UploadStepStart   UploadStep = "start"
	UploadStepAppend  UploadStep = "append"
	UploadStepReceipt UploadStep = "receipt"

	DownloadReceipt = "download-receipt"
	DownloadFile    = "download-file"

	UploadStepHeader = "upload-step"
	FileNameHeader   = "upload-filename"
)
