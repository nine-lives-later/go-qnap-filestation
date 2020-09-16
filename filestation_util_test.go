package filestation

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestRoundtrip(t *testing.T) {
	s := createTestSession(t)

	shares, err := s.GetShareList()
	if err != nil {
		t.Fatalf("Failed retrieve share list: %v", err)
	}

	var unitTestShare *FolderListEntry
	for _, s := range shares {
		if s.Path != "/home" && s.Path != "/Public" {
			unitTestShare = &s
			break
		}
	}
	if unitTestShare == nil {
		t.Fatal("Failed to determine unit test share")
	}

	t.Logf("Using unit test share: %v", unitTestShare.Path)

	// create test folder
	rand.Seed(time.Now().UnixNano())

	testFolderName := "unit-test-" + strconv.Itoa(int(rand.Int31()))
	testFolderPath := unitTestShare.Path + "/" + testFolderName

	t.Logf("Using unit test folder: %v", testFolderPath)

	t.Run("TestFolderExists1", func(t *testing.T) {
		exists, err := s.GetFileStat(testFolderPath)
		if err != nil {
			t.Fatalf("Failed retrieve file stat: %v", err)
		}

		if exists != nil {
			t.Fatal("Expected folder to not exist")
		}
	})

	t.Run("CreateTestFolder", func(t *testing.T) {
		created, err := s.CreateFolder(testFolderPath)
		if err != nil {
			t.Fatalf("Failed create test folder: %v", err)
		}

		if !created {
			t.Fatalf("Expected the test folder to not exist")
		}
	})

	defer s.Logout()
	defer s.DeleteFile(testFolderPath)

	t.Run("TryCreateTestFolder", func(t *testing.T) {
		created, err := s.CreateFolder(testFolderPath)
		if err != nil {
			t.Fatalf("Failed create test folder: %v", err)
		}

		if created {
			t.Fatalf("Expected the test folder to already exist")
		}
	})

	// test file/folder exists
	t.Run("TestFolderExists2", func(t *testing.T) {
		exists, err := s.GetFileStat(testFolderPath)
		if err != nil {
			t.Fatalf("Failed retrieve file stat: %v", err)
		}

		if exists == nil {
			t.Fatal("Expected folder to exist")
		}
	})

	t.Run("CreateTestFolder-Level2", func(t *testing.T) {
		created, err := s.CreateFolder(testFolderPath + "/test")
		if err != nil {
			t.Fatalf("Failed create test folder: %v", err)
		}

		if !created {
			t.Fatalf("Expected the test folder to not exist")
		}
	})

	t.Run("TryCreateTestFolder-Level4", func(t *testing.T) {
		_, err := s.CreateFolder(testFolderPath + "/test/do0esNotEx1st/test4")
		if err == nil {
			t.Fatal("Creating test folder should fail")
		}
	})

	t.Run("EnsureTestFolder-Level4", func(t *testing.T) {
		created, err := s.EnsureFolder(testFolderPath + "/test/ensure/test4")
		if err != nil {
			t.Fatalf("Failed create test folder: %v", err)
		}

		if created <= 0 {
			t.Fatalf("Expected the folder to not exist")
		}
	})

	t.Run("TryEnsureTestFolder-Level4", func(t *testing.T) {
		created, err := s.EnsureFolder(testFolderPath + "/test/ensure/test4")
		if err != nil {
			t.Fatalf("Failed create test folder: %v", err)
		}

		if created != 0 {
			t.Fatalf("Expected the folder to already exist")
		}
	})

	// test modifying privilege
	t.Run("SetPrivilegeRecursive-Level3", func(t *testing.T) {
		err := s.SetPrivilege(testFolderPath+"/test/ensure", 0751, true)
		if err != nil {
			t.Fatalf("Failed set privilege on parent folder: %v", err)
		}
	})

	t.Run("GetAndCheckPrivilege-Level4", func(t *testing.T) {
		stat, err := s.GetFileStat(testFolderPath + "/test/ensure/test4")
		if err != nil {
			t.Fatalf("Failed get stats on folder: %v", err)
		}
		p := NewPrivilegeFromOctal(stat.Privilege)
		if p != 0751 {
			t.Fatalf("Expected changed privilege on sub-folder")
		}
	})

	// list folders
	t.Run("GetFileList-Level2", func(t *testing.T) {
		folders, err := s.GetFileList(testFolderPath)
		if err != nil {
			t.Fatalf("Failed retrieve folder list: %v", err)
		}

		if len(folders) != 1 {
			t.Fatal("Expected one single file entry to exist")
		}
		if folders[0].IsFolder == 0 {
			t.Fatal("Expected file entry to be a folder")
		}
		if folders[0].Name != "test" {
			t.Fatal("Expected folder name to be 'test'")
		}
	})

	// fill test folder
	t.Run("FillTestFolder-Level2", func(t *testing.T) {
		for i := 0; i < 9; i++ {
			_, err := s.CreateFolder(testFolderPath + "/test-" + strconv.Itoa(int(rand.Int31())))
			if err != nil {
				t.Fatalf("Failed create test folder: %v", err)
			}
		}
	})

	// test paging
	t.Run("FileListPaging", func(t *testing.T) {
		folders, err := s.GetFileList(testFolderPath)
		if err != nil {
			t.Fatalf("Failed retrieve folder list: %v", err)
		}

		folders2, err := s.getFileListInternal(testFolderPath, 1)
		if err != nil {
			t.Fatalf("Failed retrieve folder list: %v", err)
		}

		if len(folders) != len(folders2) {
			t.Fatal("Getting file list via paging did not return the same amount of files")
		}
	})

	// test file/folder does not exists
	t.Run("FileDoesNotExist", func(t *testing.T) {
		exists, err := s.GetFileStat(testFolderPath + "/D0esN0tEx1st!!__")
		if err != nil {
			t.Fatalf("Failed retrieve file stat: %v", err)
		}

		if exists != nil {
			t.Fatal("Expected file to be missing")
		}
	})

	t.Run("TryDeleteFolder-Level2", func(t *testing.T) {
		deleted, err := s.DeleteFile(testFolderPath + "/D0esN0tEx1st456")
		if err != nil {
			t.Fatalf("Failed to delete folder: %v", err)
		}

		if deleted {
			t.Fatal("Expected folder to not be deleted")
		}
	})

	t.Run("TryDeleteFolder-Level3", func(t *testing.T) {
		deleted, err := s.DeleteFile(testFolderPath + "/D0esN0tEx1st/test3")
		if err == nil {
			t.Fatal("Expected deleting folder to fail")
		}

		if deleted {
			t.Fatal("Expected folder to not be deleted")
		}
	})

	t.Run("DeleteTestFolder", func(t *testing.T) {
		deleted, err := s.DeleteFile(testFolderPath)
		if err != nil {
			t.Fatalf("Failed to delete test folder: %v", err)
		}

		if !deleted {
			t.Fatal("Expected test folder to be deleted")
		}
	})

	t.Run("TestFolderExists3", func(t *testing.T) {
		exists, err := s.GetFileStat(testFolderPath)
		if err != nil {
			t.Fatalf("Failed retrieve file stat: %v", err)
		}

		if exists != nil {
			t.Fatal("Expected folder to not exist")
		}
	})
}
