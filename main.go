package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type organizer struct {
	TargetDir       string
	OutputDir       string
	Recursive       bool
	DeleteOnSuccess bool
	SuppressErrors  bool
}

type node struct {
	IsDir      bool
	TargetPath string
	OutputPath string
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
	var nodes []node

	fi, err := os.Stat(o.TargetDir)
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	nodes = append(nodes, node{
		IsDir:      fi.IsDir(),
		TargetPath: o.TargetDir,
		OutputPath: o.OutputDir,
	})

	for len(nodes) > 0 {
		currNode := nodes[len(nodes)-1]
		nodes = nodes[:len(nodes)-1]

		if currNode.IsDir {
			if o.Recursive {
				addNodes, err := ioutil.ReadDir(currNode.TargetPath)
				if err != nil {
					errs = append(errs, err)
					if !o.SuppressErrors {
						return errs
					}
					continue
				}

				for _, an := range addNodes {
					if an.IsDir() {
						name := an.Name()
						nodes = append(nodes, node{
							IsDir:      true,
							TargetPath: filepath.Join(currNode.TargetPath, name),
							OutputPath: filepath.Join(currNode.OutputPath, name),
						})
					} else {
						name := an.Name()
						ext := filepath.Ext(name)
						if ext == "" {
							ext = "NO_EXTENSION"
						}

						nodes = append(nodes, node{
							IsDir:      false,
							TargetPath: filepath.Join(currNode.TargetPath, name),
							OutputPath: filepath.Join(currNode.OutputPath, ext, name),
						})
					}
				}
			}
		} else {
			err := os.MkdirAll(filepath.Dir(currNode.OutputPath), os.ModePerm)
			if err != nil {
				errs = append(errs, err)
				if !o.SuppressErrors {
					return errs
				}
				continue
			}

			r, err := ioutil.ReadFile(currNode.TargetPath)
			if err != nil {
				errs = append(errs, err)
				if !o.SuppressErrors {
					return errs
				}
				continue
			}

			err = ioutil.WriteFile(currNode.OutputPath, r, 0755)
			if err != nil {
				errs = append(errs, err)
				if !o.SuppressErrors {
					return errs
				}
				continue
			}

			if o.DeleteOnSuccess {
				err := os.Remove(currNode.TargetPath)
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
	flag.BoolVar(&o.Recursive, "r", false, "[Optional, Default = false] Recurse through subdirectories.")
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
