package utils

import (
	"fmt"
	"os"
	"testing"
)

/*
	Try to download a file whose size is known. Assert the files gets downloaded and is the known size.
 */
func TestDownloadFile(t *testing.T) {
	err := DownloadFile("test_1KB.bin", "https://fastest.fish/lib/downloads/1KB.bin", nil, true)
	if err != nil {
		t.Error(err)
	}
	fi, err := os.Stat("test_1KB.bin")
	if err != nil {
		t.Error(err)
	}
	// get the size
	size := fi.Size()
	if size != 8409 {
		t.Errorf("Expected size of downloaded file to be 8409 got: %v", size)
	}
	err = os.Remove("test_1KB.bin")
	if err != nil {
		fmt.Println("Couldn't delete test file")
	}
}

/*
	Make sure DownloadFile() throws an error when given an invalid filepath
 */
func TestDownloadFileBadPath(t *testing.T) {
	err := DownloadFile("/,./&&&$$**", "https://my.site/file", nil, false)
	if err == nil {
		t.Error("Invalid filepath didn't throw an error")
	}

}

func TestDownloadFileBadURL(t *testing.T) {
	err := DownloadFile("testfile.txt", ")(*&^%$#@!@#%^&*()()", nil, false)
	if err == nil {
		t.Error("Invalid url didn't throw an error")
	}
	err = os.Remove("testfile.txt.tmp")
	if err != nil {
		fmt.Println("Couldn't delete test file")
	}

}

func TestDownloadFileWithHeaders(t *testing.T) {
	headers := make(map[string]string, 1)
	headers["User-Agent"] = "podarc_testing"
	err := DownloadFile("test_1KB_2.bin", "https://fastest.fish/lib/downloads/1KB.bin", headers, true)
	if err != nil {
		t.Error(err)
	}
	fi, err := os.Stat("test_1KB_2.bin")
	if err != nil {
		t.Error(err)
	}
	// get the size
	size := fi.Size()
	if size != 8409 {
		t.Errorf("Expected size of downloaded file to be 8409 got: %v", size)
	}
	err = os.Remove("test_1KB_2.bin")
	if err != nil {
		fmt.Println("Couldn't delete test file")
	}
}