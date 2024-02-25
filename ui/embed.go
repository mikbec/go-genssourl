package ui

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strconv"

	cp "github.com/otiai10/copy"
)

// Content_static holds our static web server content.
//
//go:embed static
var Content_static embed.FS

// content_templates holds our templates for our web server.
//
//go:embed templates
var Content_templates embed.FS

func getAllDirnames(efs *embed.FS) (dirs []string, err error) {
	if err := fs.WalkDir(efs, ".", func(path string, dir fs.DirEntry, err error) error {
		if dir.IsDir() == false {
			return nil
		}

		dirs = append(dirs, path)

		return nil
	}); err != nil {
		return nil, err
	}

	return dirs, nil
}

func getAllFilenames(efs *embed.FS) (files []string, err error) {
	if err := fs.WalkDir(efs, ".", func(path string, dir fs.DirEntry, err error) error {
		if dir.IsDir() == true {
			return nil
		}

		files = append(files, path)

		return nil
	}); err != nil {
		return nil, err
	}

	return files, nil
}

func ListEmbeddedDir(efs *embed.FS) {
	dirs, err := getAllDirnames(efs)
	if err != nil {
		log.Fatal(err)
	}
	files, err := getAllFilenames(efs)
	if err != nil {
		log.Fatal(err)
	}

	for _, fName := range dirs {
		//fmt.Println("efs-dir: " + fName)
		fStat, err := fs.Stat(efs, fName)
		if err != nil {
			log.Fatal(err)
			return
		}
		permString := fmt.Sprintf("%v", (fStat.Mode() & os.ModePerm))
		fmt.Println("efs-dirs: (" + permString +
			"," + fStat.ModTime().String() +
			"," + strconv.FormatInt(fStat.Size(), 10) +
			") " + fName)
		//") " + fStat.Name())
	}
	for _, fName := range files {
		//fmt.Println("efs-file: " + fName)
		fStat, err := fs.Stat(efs, fName)
		if err != nil {
			log.Fatal(err)
			return
		}
		permString := fmt.Sprintf("%v", (fStat.Mode() & os.ModePerm))
		fmt.Println("efs-files: (" + permString +
			"," + fStat.ModTime().String() +
			"," + strconv.FormatInt(fStat.Size(), 10) +
			") " + fName)
		//") " + fStat.Name())
	}

	return
}

func CopyEmbeddedDirWithoutDst(srcEmbedFS *embed.FS, dstRootPath string) error {
	_, err := os.Stat(dstRootPath)
	if err == nil {
		log.Print("efs-dir: destination \"" + dstRootPath + "\" exists ... could not copy")
		return errors.New("destination \"" + dstRootPath + "\" exists ... could not copy")
	}

	//log.Print("efs-dir: copy to destination \"" + dstRootPath + "\"")

	// set our source
	opt := cp.Options{
		FS: *srcEmbedFS,
		Skip: func(info os.FileInfo, src, dst string) (bool, error) {
			log.Print("do copy: " + src + " -> " + dst)
			return false, nil
		},
	}
	return cp.Copy(".", dstRootPath, opt)
}

func CopyEmbeddedDirWithDst(srcEmbedFS *embed.FS, dstRootPath string) error {
	fStat, err := os.Stat(dstRootPath)
	if err != nil {
		log.Print("efs-dir: destination \"" + dstRootPath + "\" does not exist ... could not copy")
		return errors.New("destination \"" + dstRootPath + "\"does not exist ... could not copy")
	}

	// if destination is a directory check if it is writeable
	err = nil
	if fStat.IsDir() && (fStat.Mode().Perm()&(1<<(uint(7))) == 0) {
		// write bit is not set .. set it now
		err := os.Chmod(dstRootPath, os.FileMode(0755))
		if err != nil {
			log.Print("efs-dir: Could not set mode 0755 on destination \"" +
				dstRootPath + "\", error msg: \"" +
				err.Error() + "\" ... could not copy")
			return err
		}
	}

	//log.Print("efs-dir: copy to destination \"" + dstRootPath + "\"")

	// set our source
	opt := cp.Options{
		FS: *srcEmbedFS,
		Skip: func(info os.FileInfo, src, dst string) (bool, error) {
			log.Print("do copy: " + src + " -> " + dst)
			return false, nil
		},
	}
	return cp.Copy(".", dstRootPath, opt)
}
