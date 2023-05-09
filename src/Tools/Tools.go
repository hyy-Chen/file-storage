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
	Type     int         `json:"type"`
	Children []Directory `json:"children"`
}

type FileNode struct {
	Name string `json:"name"`
	Type int    `json:"type"`
}

type Node struct {
	Name     string      `json:"name"`
	Type     int         `json:"type"`
	Children []*FileNode `json:"children"`
}

const (
	UploadFolderFiles                = "files"
	UploadFolderJson                 = "json"
	UploadFolderDestinationDirectory = "destination_directory"
	DownloadFolderPath               = "folder_path"
	QueryFolderStructurePath         = "folder_path"
	DirectoryType                    = 1
	FileType                         = 2
)
