package main

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/fynxiu/dbd/internal/config"
	"github.com/fynxiu/dbd/internal/constant"
	"github.com/fynxiu/dbd/internal/driver"
	"github.com/fynxiu/dbd/internal/transformer"
	"github.com/fynxiu/dbd/z"

	"github.com/golang/glog"
	"golang.org/x/exp/slices"
)

var (
	inputDir          = flag.String("input-dir", "", "input directory")
	engineType        = flag.String("engine-name", constant.EngineMysql, "engine type, mysql, mongo, etc.")
	image             = flag.String("image", "", "docker image")
	fileSequenceRegex = flag.String("file-seq-regex", "^[0-9]+", "file sequence regex")
	fileFilterRegex   = flag.String("file-filter-regex", "\\.up\\.", "file filter regex")
	output            = flag.String("output", "", "output filename fo the generated schema")
	reuseDocker       = flag.Bool("reuse-docker", false, "reuse docker container")
)

var engineConfig *config.EngineConfig

func init() {
	flag.Parse()
	flag.Lookup("alsologtostderr").Value.Set("true")

	engineConfig = config.GetEngineConfig(*engineType)

	if *image == "" {
		*image = engineConfig.Image
	}

	if *inputDir == "" {
		glog.Fatalln("--input-dir is required")
	}

	if *output == "" {
		*output = path.Join(*inputDir, fmt.Sprintf("../%s", engineConfig.OutputFilename))
	}

	err := os.MkdirAll(path.Dir(*output), os.ModePerm)
	requireNoErr(err)
}

func main() {
	ctx := context.Background()
	driver, err := launchDriver(ctx, *reuseDocker)
	requireNoErr(err)
	if !*reuseDocker {
		defer driver.Dispose(ctx)
	}
	glog.V(1).Infof("dsn: %s", driver.DataSourceName())

	transformer, err := transformer.NewTransformer(*engineType, driver.DataSourceName())
	requireNoErr(err)
	r, err := transformer.Transform(mustGetUpFiles(*inputDir))
	requireNoErr(err)

	err = ioutil.WriteFile(*output,
		[]byte(engineConfig.Comment(constant.GeneratedHeader)+r),
		os.ModePerm)
	requireNoErr(err)
	glog.V(1).Infof("done")
}

func launchDriver(ctx context.Context, reuse bool) (driver.Driver, error) {
	driver, err := driver.NewDefaultDriver(*engineType, *image)
	if err != nil {
		return nil, err
	}
	if reuse {
		if err := driver.Reuse(ctx); err == nil {
			return driver, nil
		} else if err != constant.ErrContainerNotFound {
			return nil, err
		}
	}
	if err := driver.Launch(ctx); err != nil {
		return nil, err
	}
	return driver, nil
}

func mustGetUpFiles(inputDir string) []string {
	files, err := ioutil.ReadDir(inputDir)
	requireNoErr(err)
	fsr, err := regexp.Compile(*fileSequenceRegex)
	requireNoErr(err)
	ffr, err := regexp.Compile(*fileFilterRegex)
	requireNoErr(err)

	upfiles := z.FilterMap(files, func(fi fs.FileInfo, _ int) (string, bool) {
		return filepath.Join(inputDir, fi.Name()), ffr.MatchString(fi.Name())
	}).Slice()

	slices.SortFunc(upfiles, func(a, b string) bool {
		return fsr.FindString(a) < fsr.FindString(b)
	})
	glog.V(2).Infof("upfiles: %v", upfiles)

	return upfiles
}

func requireNoErr(err error) {
	if err != nil {
		glog.FatalDepth(2, err)
		panic(err)
	}
}
