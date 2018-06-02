package filestation

import (
	"fmt"
	"path/filepath"
	"strconv"
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

	err := s.getForEntity(&result, "cgi-bin/filemanager/utilRequest.cgi", QueryParameters{
		"func": "get_tree",
		"node": "share_root",
	})
	if err != nil {
		return nil, err
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
		var result *getFileListResponse

		err := s.getForEntity(&result, "cgi-bin/filemanager/utilRequest.cgi", QueryParameters{
			"func":      "get_list",
			"path":      path,
			"list_mode": "all",
			"dir":       "ASC",
			"limit":     strconv.Itoa(limit),
			"start":     strconv.Itoa(len(ret)),
		})
		if err != nil {
			return nil, err
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
	var result *getFileListResponse

	err := s.getForEntity(&result, "cgi-bin/filemanager/utilRequest.cgi", QueryParameters{
		"func":       "stat",
		"path":       filepath.ToSlash(filepath.Dir(path)),
		"file_name":  filepath.Base(path),
		"file_total": "1",
	})
	if err != nil {
		return nil, err
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
	Status int `json:"status,omitempty"`
}

// CreateFolder creates a new folder.
// The base directory must exists.
func (s *FileStationSession) CreateFolder(path string) (bool, error) {
	var result *createFolderResponse

	err := s.getForEntity(&result, "cgi-bin/filemanager/utilRequest.cgi", QueryParameters{
		"func":        "createdir",
		"dest_path":   filepath.ToSlash(filepath.Dir(path)),
		"dest_folder": filepath.Base(path),
	})
	if err != nil {
		return false, err
	}

	switch result.Status {
	case 1: // success
		return true, nil
	case 3: // session expired
		return false, fmt.Errorf("Session expired")
	case 5: // base directory does not exist
		return false, fmt.Errorf("Base directory does not exist")
	case 2: // folder already exists
		return false, nil
	case 33: // folder already exists
		return false, nil
	case 4: // permission denied
		return false, fmt.Errorf("Permission denied")
	}

	return false, fmt.Errorf("Unknown status code: %v", result.Status)
}

// DeleteFile deletes a file or folder.
func (s *FileStationSession) DeleteFile(path string) (bool, error) {
	var result *createFolderResponse

	err := s.getForEntity(&result, "cgi-bin/filemanager/utilRequest.cgi", QueryParameters{
		"func":       "delete",
		"path":       filepath.ToSlash(filepath.Dir(path)),
		"file_name":  filepath.Base(path),
		"file_total": "1",
	})
	if err != nil {
		return false, err
	}

	switch result.Status {
	case 1: // success
		return true, nil
	case 0: // file not found
		return false, nil
	case 25: // base directory does not exist
		return false, fmt.Errorf("Base directory does not exist")
	case 3: // session expired
		return false, fmt.Errorf("Session expired")
	case 4: // permission denied
		return false, fmt.Errorf("Permission denied")
	}

	return false, fmt.Errorf("Unknown status code: %v", result.Status)
}
