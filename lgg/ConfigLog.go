package lgg

import (
	"github.com/Sirupsen/logrus"
	_ "github.com/gogap/logrus_mate/hooks/file"
	_ "github.com/gogap/logrus_mate/hooks/slack"
	"github.com/gogap/logrus_mate"
	_"path/filepath"
	_"os"
	"path/filepath"
	"os"
	"github.com/kardianos/service"
	"runtime"
	"path"
	"strings"
	"github.com/kardianos/osext"
)

var Lgdef *logrus.Logger

func InitLogger(svclgg service.Logger) (err error) {

	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				svclgg.Error(r)
				panic(r)
			}
			err = r.(error);
			if err!=nil {
				svclgg.Error(err)
			}
		}
	}()

	var file string = ""

	if service.Interactive() {
		//Comprobar que existe la carpeta log.
		folderLog,_ := filepath.Abs("./log")
		if _,err := os.Stat(folderLog); os.IsNotExist(err) {
			//folder no existe
			os.Mkdir(folderLog,os.ModePerm)
		}

		fconfig := "logger.conf"
		file,_ = filepath.Abs("./" + fconfig)

	} else {

		//Comprobar que existe la carpeta log.
		pathexec,_ :=osext.ExecutableFolder();
		folderLog := strings.Replace(path.Join(pathexec,"log"),"/","\\",-1)
		if _,err := os.Stat(folderLog); os.IsNotExist(err) {
			//folder no existe
			os.Mkdir(folderLog,os.ModePerm)
		}

		fconfig := "logger.conf"
		file=strings.Replace(path.Join(pathexec, fconfig),"/","\\",-1)
	}

	mate, err := logrus_mate.NewLogrusMate(
		logrus_mate.ConfigFile(file),
	)

	if (err!=nil) {
		svclgg.Error(err)
		return err
	}

	Lgdef = mate.Logger("default")

	return nil
}