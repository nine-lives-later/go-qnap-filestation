package filestation

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

type FolderListEntry struct {
	Path          string `json:"id,omitempty"`
	CLS           string `json:"cls,omitempty"`
	Text          string `json:"text,omitempty"`
	Icon          string `json:"iconCls,omitempty"`
	RecycleBin    string `json:"recycle_bin,omitempty"`
	RecycleFolder string `json:"recycle_folder,omitempty"`
	MaxItemLimit  int    `json:"max_item_limit,omitempty"`
	ItemCount     int    `json:"real_total,omitempty"`
}

// GetShareList retrieves the list of shares.
func (s *FileStationSession) GetShareList() ([]FolderListEntry, error) {
	var result []FolderListEntry

	res, err := s.conn.NewRequest().
		ExpectContentType("application/json").
		SetQueryParam("func", "get_tree").
		SetQueryParam("node", "share_root").
		SetResult(&result).
		Get("cgi-bin/filemanager/utilRequest.cgi")
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %v", err)
	}
	if res.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to perform request: unexpected HTTP status code: %v", res.StatusCode())
	}

	return result, nil
}

type FileListEntry struct {
	Name              string `json:"filename,omitempty"`
	FullPath          string
	Exists            int    `json:"exist,omitempty"`
	IsFolder          int    `json:"isfolder,omitempty"`
	FileSize          int64  `json:"filesize,omitempty,string"`
	Group             string `json:"group,omitempty"`
	Owner             string `json:"owner,omitempty"`
	IsCompressed      int    `json:"iscommpressed,omitempty"`
	Privilege         string `json:"privilege,omitempty"`
	PrivilegeEx       int    `json:"privilege_ex,omitempty"`
	FileType          int    `json:"filetype,omitempty"`
	ModifiedDate      int    `json:"epochmt,omitempty"`
	HasStickyBit      int    `json:"sticky_bit,omitempty"`
	IsFolderEncrypted int    `json:"encrypt_folder,omitempty"`
	ProjectionType    int    `json:"projection_type,omitempty"`
}

type getFileListResponse struct {
	ItemCount       int             `json:"real_total,omitempty"`
	ACL             int             `json:"acl,omitempty"`
	IsACLEnabled    int             `json:"is_acl_enable,omitempty"`
	IsWinACLEnabled int             `json:"is_winacl_enable,omitempty"`
	Entries         []FileListEntry `json:"datas,omitempty"`
}

// GetFileList retrieves the list of files and folders of a share.
func (s *FileStationSession) GetFileList(path string) ([]FileListEntry, error) {
	return s.getFileListInternal(path, 1000)
}

func (s *FileStationSession) getFileListInternal(path string, limit int) ([]FileListEntry, error) {
	ret := make([]FileListEntry, 0)

	for true {
		var result getFileListResponse

		res, err := s.conn.NewRequest().
			ExpectContentType("application/json").
			SetQueryParam("func", "get_list").
			SetQueryParam("path", path).
			SetQueryParam("list_mode", "all").
			SetQueryParam("dir", "ASC").
			SetQueryParam("limit", strconv.Itoa(limit)).
			SetQueryParam("start", strconv.Itoa(len(ret))).
			SetResult(&result).
			Get("cgi-bin/filemanager/utilRequest.cgi")
		if err != nil {
			return nil, fmt.Errorf("failed to perform request: %v", err)
		}
		if res.StatusCode() != 200 {
			return nil, fmt.Errorf("failed to perform request: unexpected HTTP status code: %v", res.StatusCode())
		}

		// copy entries
		ret = append(ret, result.Entries...)

		// reached last entry?
		if len(result.Entries) < limit {
			break
		}
	}

	// inject full path
	for i := range ret {
		e := &ret[i]

		e.FullPath = filepath.ToSlash(filepath.Join(path, e.Name))
	}

	return ret, nil
}

// GetFileStat checks if a file or folder exists.
func (s *FileStationSession) GetFileStat(path string) (*FileListEntry, error) {
	var result getFileListResponse

	res, err := s.conn.NewRequest().
		ExpectContentType("application/json").
		SetQueryParam("func", "stat").
		SetQueryParam("path", filepath.ToSlash(filepath.Dir(path))).
		SetQueryParam("file_name", filepath.Base(path)).
		SetQueryParam("file_total", "1").
		SetResult(&result).
		Get("cgi-bin/filemanager/utilRequest.cgi")
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %v", err)
	}
	if res.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to perform request: unexpected HTTP status code: %v", res.StatusCode())
	}

	if len(result.Entries) <= 0 {
		return nil, nil
	}

	// handle file entry
	entry := &result.Entries[0]

	if entry.Exists == 0 {
		return nil, nil
	}

	// inject full path
	entry.FullPath = filepath.ToSlash(path)

	return entry, nil
}

type createFolderResponse struct {
	Status FileStationStatus `json:"status,omitempty"`
}

// CreateFolder creates a new folder.
// The base directory must exists.
func (s *FileStationSession) CreateFolder(path string) (bool, error) {
	var result createFolderResponse

	res, err := s.conn.NewRequest().
		ExpectContentType("application/json").
		SetQueryParam("func", "createdir").
		SetQueryParam("dest_path", filepath.ToSlash(filepath.Dir(path))).
		SetQueryParam("dest_folder", filepath.Base(path)).
		SetResult(&result).
		Get("cgi-bin/filemanager/utilRequest.cgi")
	if err != nil {
		return false, fmt.Errorf("failed to perform request: %v", err)
	}
	if res.StatusCode() != 200 {
		return false, fmt.Errorf("failed to perform request: unexpected HTTP status code: %v", res.StatusCode())
	}

	switch result.Status {
	case WFM2_SUCCESS: // success
		return true, nil
	case WFM2_FILE_EXIST, WFM2_NAME_DUP: // folder already exists
		return false, nil
	}

	return false, result.Status
}

// EnsureFolder creates a new folder and its parent directories.
func (s *FileStationSession) EnsureFolder(path string) (int, error) {
	if !strings.HasPrefix(path, "/") {
		return 0, fmt.Errorf("path does not begin with a slash: %v", path)
	}

	// already exists?
	exists, err := s.GetFileStat(path)
	if err != nil {
		return 0, fmt.Errorf("failed to check for folder '%v': %v", path, err)
	}
	if exists != nil {
		return 0, nil
	}

	// create sub-folders
	parts := strings.Split(filepath.ToSlash(path), "/")[1:]
	if len(parts) < 2 {
		return 0, fmt.Errorf("path is not a subfolder of a share: %v", path)
	}

	createdOverall := 0
	for i := 1; i < len(parts); i++ {
		subPath := "/" + filepath.ToSlash(filepath.Join(parts[0:i+1]...))

		created, err := s.CreateFolder(subPath)
		if err != nil {
			return createdOverall, fmt.Errorf("failed to create sub-folder '%v': %v", subPath, err)
		}

		if created {
			createdOverall++
		}
	}

	return createdOverall, nil
}

// DeleteFile deletes a file or folder.
func (s *FileStationSession) DeleteFile(path string) (bool, error) {
	return s.deleteFileInternal(path, false)
}

// DeleteFileNoRecycleBin deletes a file or folder without moving them to the recycling bin.
func (s *FileStationSession) DeleteFileNoRecycleBin(path string) (bool, error) {
	return s.deleteFileInternal(path, true)
}

func (s *FileStationSession) deleteFileInternal(path string, force bool) (bool, error) {
	forceStr := "0"
	if force {
		forceStr = "1"
	}

	var result createFolderResponse

	res, err := s.conn.NewRequest().
		ExpectContentType("application/json").
		SetQueryParam("func", "delete").
		SetQueryParam("path", filepath.ToSlash(filepath.Dir(path))).
		SetQueryParam("file_name", filepath.Base(path)).
		SetQueryParam("file_total", "1").
		SetQueryParam("force", forceStr).
		SetResult(&result).
		Get("cgi-bin/filemanager/utilRequest.cgi")
	if err != nil {
		return false, fmt.Errorf("failed to perform request: %v", err)
	}
	if res.StatusCode() != 200 {
		return false, fmt.Errorf("failed to perform request: unexpected HTTP status code: %v", res.StatusCode())
	}

	switch result.Status {
	case WFM2_SUCCESS: // success
		return true, nil
	case WFM2_FAIL, WFM2_PERMISSION_DENY: // file not found
		return false, nil
	}

	return false, result.Status
}
