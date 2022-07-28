package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type organizer struct {
	TargetDir       string
	OutputDir       string
	DeleteOnSuccess bool
	SuppressErrors  bool
}

func (o organizer) validate() error {
	if o.TargetDir == "" {
		return errors.New("-t flag is required.")
	}

	if _, err := os.Stat(o.TargetDir); os.IsNotExist(err) {
		return errors.New("-t flag is invalid. Directory does not exist.")
	}

	if o.OutputDir == "" {
		return errors.New("-o flag is required.")
	}

	err := os.MkdirAll(o.OutputDir, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (o organizer) cleanup() []error {
	var errs []error

	files, err := ioutil.ReadDir(o.TargetDir)
	if err != nil {
		errs = append(errs, err)
		return errs
	}

	for len(files) > 0 {
		file := files[len(files)-1]
		files = files[:len(files)-1]

		if !file.IsDir() {
			name := strings.TrimSpace(file.Name())
			ext := strings.TrimSpace(filepath.Ext(file.Name()))

			if ext == "" {
				ext = "NO_EXTENSION"
			}

			err := os.MkdirAll(filepath.Join(o.OutputDir, ext), os.ModePerm)
			if err != nil {
				errs = append(errs, err)
				if !o.SuppressErrors {
					return errs
				}
				continue
			}

			r, err := ioutil.ReadFile(filepath.Join(o.TargetDir, name))
			if err != nil {
				errs = append(errs, err)
				if !o.SuppressErrors {
					return errs
				}
				continue
			}

			err = ioutil.WriteFile(filepath.Join(o.OutputDir, ext, name), r, 0755)
			if err != nil {
				errs = append(errs, err)
				if !o.SuppressErrors {
					return errs
				}
				continue
			}

			if o.DeleteOnSuccess {
				err := os.Remove(filepath.Join(o.TargetDir, name))
				if err != nil {
					errs = append(errs, err)
					if !o.SuppressErrors {
						return errs
					}
					continue
				}
			}
		}
	}

	return errs
}

func main() {
	o := &organizer{}

	flag.StringVar(&o.TargetDir, "t", "", "[Required] Target directory to clean up. Directory must exist.")
	flag.StringVar(&o.OutputDir, "o", "", "[Required] Output directory to place files. Directory is created if does not exist.")
	flag.BoolVar(&o.DeleteOnSuccess, "d", false, "[Optional, Default = false] Delete files after successful move.")
	flag.BoolVar(&o.SuppressErrors, "s", false, "[Optional, Default = false] Suppress errors.")
	flag.Parse()

	err := o.validate()
	if err != nil {
		fmt.Printf("%s\n\n", err)
		flag.Usage()
		os.Exit(1)
	}

	errs := o.cleanup()
	if !o.SuppressErrors && len(errs) != 0 {
		for _, err := range errs {
			fmt.Printf("%s\n", err)
		}
		fmt.Print("\n")
		os.Exit(1)
	}

	os.Exit(0)
}
