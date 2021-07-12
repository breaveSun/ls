package zip

import (
	"archive/zip"
	"fmt"
	"go.uber.org/zap"
	"io"
	"log"
	"ls/internal/pkg/common"
	"ls/internal/pkg/lib/logger"
	"os"
	"path/filepath"
	"strings"
	"sync"
)
//压缩文件
//files 文件数组，可以是不同dir下的文件或者文件夹
//dest 压缩文件存放地址

func Compress(source, dest string) error {
	// 创建准备写入的文件
	fw, err := os.Create(dest)
	defer fw.Close()
	if err != nil {
		return err
	}

	// 通过 fw 来创建 zip.Write
	zw := zip.NewWriter(fw)
	defer func() {
		// 检测一下是否成功关闭
		if err := zw.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// 下面来将文件写入 zw ，因为有可能会有很多个目录及文件，所以递归处理
	return filepath.Walk(source, func(path string, fi os.FileInfo, errBack error) (err error) {
		if errBack != nil {
			return errBack
		}

		if fi.IsDir(){
			return
		}

		// 通过文件信息，创建 zip 的文件信息
		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return
		}

		// 替换文件信息中的文件名
		fh.Name = strings.TrimPrefix(filepath.ToSlash(path), filepath.ToSlash(source))

		fh.Name = strings.TrimPrefix(fh.Name, "/")

		// 写入文件信息，并返回一个 Write 结构
		w, err := zw.CreateHeader(fh)
		if err != nil {
			return
		}

		// 检测，如果不是标准文件就只写入头信息，不写入文件数据到 w
		// 如目录，也没有数据需要写
		if !fh.Mode().IsRegular() {
			return nil
		}

		// 打开要压缩的文件
		fr, err := os.Open(path)
		defer fr.Close()
		if err != nil {
			return
		}

		// 将打开的文件 Copy 到 w
		n, err := io.Copy(w, fr)
		if err != nil {
			return
		}

		// 输出压缩的内容
		fmt.Printf("成功压缩文件： %s, 共写入了 %d 个字符的数据\n", path, n)

		return nil
	})
}
//解压缩文件
//files 文件数组，可以是不同dir下的文件或者文件夹
//dest 压缩文件存放地址
func Decompress(source, dest string) error {
	if _, err := os.Stat(source); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(source, 0755)
			if err != nil{
				logger.Logger.Error("Decompress mkdir err",zap.String("errMsg",err.Error()))
				return err
			}
		}
	}
	reader, err := zip.OpenReader(dest)
	if err != nil {
		logger.Logger.Error("Decompress Openzip err",zap.String("errMsg",err.Error()))
		return err
	}

	defer  reader.Close()
	var wg = sync.WaitGroup{}
	wgNum := len(reader.File)
	wg.Add(wgNum)
	var fileErr = make(chan error,wgNum)
	for _, file := range reader.File {
		fileName := source+"\\"+file.Name
		if file.FileInfo().IsDir() {
			wg.Done()
			err := os.MkdirAll(fileName, 0755)
			if err != nil {
				fileErr<-err
				logger.Logger.Error("Decompress (IsDir) Close zip err",zap.String("errMsg",err.Error()))
				return err
			}
			fileErr<-nil
		} else {
			p := common.BaseHandler{}.GetDirPath(fileName)
			err = os.MkdirAll(p, 0755)
			if err != nil {
				wg.Done()
				logger.Logger.Error("Decompress Close zip err",zap.String("errMsg",err.Error()))
				return err
			}
			go OneFileCopy(&wg,file,fileName,fileErr)
		}
	}
	wg.Wait()
	for i:=0;i<wgNum;i++ {
		o,ok := <-fileErr
		if !ok || o!=nil{
			return o
		}
	}
	return nil
}
func OneFileCopy(wg *sync.WaitGroup,file *zip.File,fileName string,fileErr chan error) {
	defer func(wg *sync.WaitGroup) {
		wg.Done()
	}(wg)
	rc, err := file.Open()
	if err != nil {
		logger.Logger.Error("Decompress file.Open err", zap.String("errMsg", err.Error()))
		fileErr<-err
		return
	}
	defer rc.Close()
	w, err := os.Create(fileName)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Decompress Create %s err",fileName), zap.String("errMsg", err.Error()))
		fileErr<-err
		return
	}
	defer w.Close()
	_, err = io.Copy(w, rc)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Decompress copy %s err",fileName), zap.String("errMsg", err.Error()))
		fileErr<-err
		return
	}
	fileErr <- nil
	return
}
