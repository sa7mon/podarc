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
	err := DownloadFile("test_file_10MB.zip", "http://speedtest.tele2.net/10MB.zip", nil, false)
	if err != nil {
		t.Error(err)
	}
	fi, err := os.Stat("test_file_10MB.zip")
	if err != nil {
		t.Error(err)
	}
	// get the size
	size := fi.Size()
	if size != 10485760 {
		t.Errorf("Expected size of downloaded file to be 10485760 got: %v", size)
	}
	err = os.Remove("test_file_10MB.zip")
	if err != nil {
		fmt.Println("Couldn't delete test file")
	}
}