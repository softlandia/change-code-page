//Convert text file code page
//>cpc 866 1251 .las
//convert all las files from 866 to 1251
//search from current folder with recursion
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/softlandia/xLib"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func main() {

	if len(os.Args) < 4 {
		fmt.Println("using >cpc 866 1251")
		os.Exit(-1)
	}

	icp, err := strconv.ParseInt(os.Args[1], 10, 32)
	if err != nil {
		fmt.Printf("input code page '%s' not recognize\n", os.Args[1])
		os.Exit(-2)
	}

	ocp, err := strconv.ParseInt(os.Args[2], 10, 32)
	if err != nil {
		fmt.Printf("output code page '%s' not recognize\n", os.Args[2])
		os.Exit(-3)
	}

	switch icp {
	case 866:
		icp = xLib.Cp866
	case 1251:
		icp = xLib.CpWindows1251
	default:
		fmt.Printf("input code page '%s' not support\n", os.Args[1])
		os.Exit(-4)
	}

	switch ocp {
	case 866:
		ocp = xLib.Cp866
	case 1251:
		ocp = xLib.CpWindows1251
	default:
		fmt.Printf("output code page '%s' not support\n", os.Args[2])
		os.Exit(-5)
	}

	if icp == ocp {
		fmt.Printf("input '%s' and output code page '%s' not equal\n", os.Args[1], os.Args[2])
		os.Exit(-6)
	}

	fileList := make([]string, 0, 10)
	i, _ := findFiles(&fileList, filepath.Dir(os.Args[0]), os.Args[3])
	fmt.Printf("founded :%v files\n", i)
	t0 := time.Now()
	i, err = convertFiles(&fileList, int(icp), int(ocp))
	fmt.Printf("elapsed time: %v\n", time.Since(t0))
	os.Exit(0)
}

func fileTrust(ext, path string, i os.FileInfo) bool {
	if i.IsDir() { //skip dir
		return false
	}

	if strings.ToUpper(filepath.Ext(path)) != ext { //skip files with extention not equal extFileName
		return false
	}
	return true
}

func findFiles(fileList *[]string, path, extFile string) (int, error) {
	log.Println("start search")
	log.Println("start path: " + path)
	log.Println("file name mask: " + extFile)

	extFile = strings.ToUpper(extFile)
	i := 0 //index founded files
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if !fileTrust(extFile, path, info) {
			return nil
		}
		//file found
		i++
		*fileList = append(*fileList, path)
		return nil
	})
	return i, err
}

func convertFiles(fileList *[]string, fromCP, toCP int) (int, error) {
	fmt.Println("<start convert>")
	fmt.Println("file count to look: ", len(*fileList))
	for _, fn := range *fileList {
		if (fn[0] == '~') && (fn[1] == '~') {
			fmt.Printf("file '%s' removed from computer!\n", fn)
			os.Remove(fn)
			continue
		}

		newPath := ""
		newFileName := "~~" + filepath.Base(fn)
		dir := filepath.Dir(fn)
		if dir == "." {
			newPath = newFileName
		} else {
			newPath = fmt.Sprintf("%s\\%s", dir, newFileName)
		}

		fcp, err := xLib.CodePageDetect(fn)
		if err != nil {
			continue
		}

		fmt.Printf("file '%s' have code page: %v and ", fn, fcp)
		if fcp == fromCP {
			changeCpFile(fn, newPath, fromCP, toCP)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
				//continue
			}
			fmt.Printf("convert from cp: %v to %v, file replaced'\n", fcp, toCP)
			os.Rename(newPath, fn)
		} else {
			fcp, _ = xLib.CodePageDetect(fn)
			fmt.Println("skip")
		}
	}
	return 0, nil
}

func changeCpFile(iFileName, oFileName string, fromCP, toCP int) error {
	iFile, err := os.Open(iFileName)
	if err != nil {
		return err
	}
	defer iFile.Close()
	oFile, err := os.Create(oFileName)
	if err != nil {
		return err
	}
	defer oFile.Close()

	s := ""
	iScanner := bufio.NewScanner(iFile)
	for i := 0; iScanner.Scan(); i++ {
		s = iScanner.Text()
		s, _, err = transform.String(charmap.CodePage866.NewDecoder(), s)
		s, _, err = transform.String(charmap.Windows1251.NewEncoder(), s)
		if err != nil {
			fmt.Printf("error on file '%s' convert\n", iFileName)
			return err
		}
		fmt.Fprintf(oFile, "%s\n", s)
	}
	return nil
}
