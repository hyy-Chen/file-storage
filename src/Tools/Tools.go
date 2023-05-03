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

type CreateFolderRequest struct {
	Path string `json:"create_folder_path"`
}

type Test struct {
	FilePath string `json:"file_path"`
}

type Directory struct {
	CName    string      `json:"cname"`
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Children []Directory `json:"children"`
}

const (
	UploadFolderFiles                = "files"
	UploadFolderJson                 = "json"
	Folder                           = "directory" // type
	UploadFolderDestinationDirectory = "destination_directory"
)
