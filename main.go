//(c) softlandia@gmail.com
//Convert text file code page
//using
//>cpc 866 1251 .las
//convert all las files from 866 to 1251
//search from current folder with recursion
//>cpc 866 1251 x:\prj\data\plat-1.las
//convert file "x:\prj\data\plat-1.las" from 866 to 1251 code page, if input file already 1251 code page then nothing to do

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/softlandia/xLib"
)

func parsParameters(icp, ocp *int64) error {
	if len(os.Args) < 4 {
		fmt.Println("using >cpc 866 1251 .las")
		fmt.Println("where 866 input code page")
		fmt.Println("where 1251 output code page")
		fmt.Println("where .las file extention")
		return errors.New("no enought parameters")
	}

	var err error

	*icp, err = strconv.ParseInt(os.Args[1], 10, 32)
	if err != nil {
		fmt.Printf("input code page '%s' not recognize\n", os.Args[1])
		return errors.New("input code page is bad")
	}

	*ocp, err = strconv.ParseInt(os.Args[2], 10, 32)
	if err != nil {
		fmt.Printf("output code page '%s' not recognize\n", os.Args[2])
		return errors.New("output code page is bad")
	}

	switch *icp {
	case 866:
		*icp = xlib.Cp866
	case 1251:
		*icp = xlib.CpWindows1251
	default:
		fmt.Printf("input code page '%s' not support\n", os.Args[1])
		return errors.New("error input code page")
	}

	switch *ocp {
	case 866:
		*ocp = xlib.Cp866
	case 1251:
		*ocp = xlib.CpWindows1251
	default:
		fmt.Printf("output code page '%s' not support\n", os.Args[2])
		return errors.New("error output code page")
	}

	if *icp == *ocp {
		fmt.Printf("input '%s' and output code page '%s' is equal, nothing to do\n", os.Args[1], os.Args[2])
		return errors.New("last parameter error")
	}

	if len(os.Args[3]) < 2 {
		fmt.Printf("path to search or ext of files '%s' to search to small\n", os.Args[3])
		return errors.New("last parameter error")
	}
	return nil
}

func main() {

	var (
		icp,
		ocp int64
		path string
	)

	err := parsParameters(&icp, &ocp)
	if err != nil {
		log.Printf("program stop, error: %s", err)
		os.Exit(1)
	}

	path = os.Args[3]
	log.Printf("input cp: %d\n", icp)
	log.Printf("output cp: %d\n", ocp)
	log.Printf("ext: %s", path)

	t0 := time.Now()
	if xlib.FileExists(path) {
		log.Println("change code page at one file")
		fcp, err := xlib.CodePageDetect(path)
		if err != nil {
			log.Printf("error file '%s'. %v", path, err)
		}
		if int64(fcp) == icp {
			err = xlib.ReplaceCpFile(path, icp, ocp)
			if err != nil {
				log.Printf("error file '%s' convert. %v", path, err)
			} else {
				log.Printf("file '%s' converted", path)
			}
		} else {
			log.Printf("file '%s' already converted.", path)
		}
	} else {
		log.Println("change code page at many files")
		fileList := make([]string, 0, 10)
		i, _ := xlib.FindFilesExt(&fileList, filepath.Dir(os.Args[0]), path)
		fmt.Printf("founded :%v files\n", i)
		i, err = convertFiles(&fileList, icp, ocp)
	}
	fmt.Printf("elapsed time: %v\n", time.Since(t0))
}

func convertFiles(fileList *[]string, fromCP, toCP int64) (int, error) {
	fmt.Println("<start convert>")
	fmt.Println("file count to look: ", len(*fileList))
	for _, fn := range *fileList {
		/*if (fn[0] == '~') && (fn[1] == '~') {
			fmt.Printf("file '%s' removed from computer!\n", fn)
			os.Remove(fn)
			continue
		}*/

		fcp, err := xlib.CodePageDetect(fn)
		if err != nil {
			continue
		}

		fmt.Printf("file '%s' have code page: %v and ", fn, fcp)
		if int64(fcp) == fromCP {
			fmt.Printf("convert from cp: %v to %v, file replaced'\n", fcp, toCP)
			err = xlib.ReplaceCpFile(fn, fromCP, toCP)
			if err != nil {
				fmt.Println(err)
				os.Exit(1) //TODO need parameter to control. stop on error or continue
			}
		} else {
			fmt.Println("skip")
		}
	}
	return 0, nil
}
