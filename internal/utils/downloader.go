/*
*  Taken in large part from: https://golangcode.com/download-a-file-with-progress/ by Edd Turtle
*  License: MIT (https://github.com/eddturtle/golangcode-site/blob/master/LICENSE)
 */

package utils

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"io"
	"net/http"
	"os"
	"strings"
)

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer interface
// and we can pass this into io.TeeReader() which will report progress on each write cycle.
type WriteCounter struct {
	TotalBytes uint64
	DoPrintProgress bool
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.TotalBytes += uint64(n)

	if wc.DoPrintProgress {
		wc.PrintProgress()
	}
	return n, nil
}

func bytesToMegabytes(numBytes uint64) uint64 {
	return numBytes / 100000000
}

func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.TotalBytes))
}


// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory. We pass an io.TeeReader
// into Copy() to report progress on the download.
func DownloadFile(filepath string, url string, headers map[string]string, printProgress bool) error {
	// Create the file, but give it a tmp file extension, this means we won't overwrite a
	// file until it's downloaded, but we'll remove the tmp extension once downloaded.
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}

	client := &http.Client{}

	// Get the data
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	//resp, err := http.Get(url)
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		out.Close()
		return err
	}
	defer resp.Body.Close()

	// Create our progress reporter and pass it to be used alongside our writer
	counter := &WriteCounter{DoPrintProgress: printProgress}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return err
	}

	// Move "print cursor" back to the beginning of the line so the progress line gets overwritten
	if printProgress {
		fmt.Printf("\r")
	}

	// Close the file without defer so it can happen before Rename()
	out.Close()

	if err = os.Rename(filepath+".tmp", filepath); err != nil {
		return err
	}
	return nil
}