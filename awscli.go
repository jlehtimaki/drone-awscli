package main

import (
	"archive/zip"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

const awsCliExe = "aws"

var (
	zipPath     = "/var/tmp/awscli.zip"
	path        = "/var/tmp/"
	downloadUrl = "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip"
)

func installAWSCli() error {
	err := downloadAWSCli()
	if err != nil {
		return err
	}
	err = Unzip(zipPath, path)
	if err != nil {
		return fmt.Errorf("error unzipping file: %w", err)
	}
	return installBinary()
}

func installBinary() error {
	cmd     := exec.Command("./install","-b","/bin/")
	cmd.Dir  = fmt.Sprintf("%s/aws", path)
	err     := cmd.Run()

	if err != nil  {
		return fmt.Errorf("error installing awsCli: %w", err)
	}
	return nil
}

func downloadAWSCli() error {
	return downloadFile(zipPath, downloadUrl)
}

func downloadFile(filepath string, url string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// Unzip a file to a destination
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
