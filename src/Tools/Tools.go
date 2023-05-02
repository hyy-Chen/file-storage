package Tools

type DownloadRequest struct {
	Filename string `json:"file_name"`
}

type MoveFilesRequest struct {
	SourcePath      string `json:"source_path"`
	DestinationPath string `json:"destination_path"`
}

type DeleteFilesRequest struct {
	Path string `json:"delete_files_path"`
}
