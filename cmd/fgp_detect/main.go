package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/ProninIgorr/fingerprint/internal/contexts"
	erf "github.com/ProninIgorr/fingerprint/internal/errflow"
	"github.com/ProninIgorr/fingerprint/internal/fh"
	"github.com/ProninIgorr/fingerprint/internal/imgs/helpers"
	"github.com/ProninIgorr/fingerprint/internal/logging"
	out "github.com/ProninIgorr/fingerprint/internal/output"
	"github.com/ProninIgorr/fingerprint/internal/workflow"
	"github.com/ProninIgorr/fingerprint/internal/workflow/inputs"
	"github.com/ProninIgorr/fingerprint/internal/workflow/loaders"
	conf "github.com/ilyakaznacheev/cleanenv"

	"os"
	"os/signal"
	"os/user"
	fp "path/filepath"
	"time"

	"github.com/ProninIgorr/fingerprint/internal/errs"
	log "github.com/sirupsen/logrus"
)

var (
	AppName           = fp.Base(os.Args[0])
	DefaultConfigFile = fmt.Sprintf("%s.yml", AppName)
)

// Config is application configuration structure
type Config struct {
	Log struct {
		Level  string `yaml:"level" env:"FINGERS_LOG_LEVEL" env-default:"info" env-description:"log level: error, info, warn"`
		Format string `yaml:"format" env:"FINGERS_LOG_FORMAT" env-default:"text" env-description:"log format: text, json, csv"`
		File   string `yaml:"file" env:"FINGERS_LOG_FILE" env-default:"log.txt" env-description:"log file"`
	} `yaml:"log"`
	TraceFile       string        `yaml:"trace_file" env:"FINGERS_TRACE_FILE" env-description:"trace file"`
	StatsUpdateRate time.Duration `yaml:"stats_update_rate" env-default:"2s" env-description:"stats refresh interval"`

	Fingers struct {
		InputDir  string `yaml:"input_dir" env:"FINGERS_INPUTS" env-default:"inputs" env-description:"inputs"`
		OutputDir string `yaml:"output_dir" env:"FINGERS_OUTPUTS" env-default:"outputs" env-description:"outputs"`
		ReportDir string `yaml:"report_dir" env:"FINGERS_REPORTS" env-default:"reports" env-description:"reports"`
		Image     struct {
			Format string `yaml:"format" env:"FINGERS_IMAGE_FORMAT" env-default:"*.bmp}" env-description:"image format"`
			Width  int    `yaml:"width" env:"FINGERS_IMAGE_WIDTH" env-default:"96" env-description:"image width"`
			Height int    `yaml:"height" env:"FINGERS_IMAGE_HEIGHT" env-default:"103" env-description:"image height"`
		} `yaml:"image"`
	} `yaml:"fingers"`

	Greeting string `env:"GREETING" env-description:"Greeting phrase" env-default:"Privet!"`
}

type Args struct {
	ConfigPath string
}

func ProcessArgs(cfg interface{}) Args {
	var a Args

	f := flag.NewFlagSet("Example server", 1)
	f.StringVar(&a.ConfigPath, "c", DefaultConfigFile, "Path to configuration file")

	fu := f.Usage
	f.Usage = func() {
		fu()
		envHelp, _ := conf.GetDescription(cfg, nil)
		fmt.Fprintln(f.Output())
		fmt.Fprintln(f.Output(), envHelp)
	}

	f.Parse(os.Args[1:])
	return a
}

// TODO: move global vars to app.config
var (
	cfg         Config
	startTime   time.Time
	currentUser *user.User
)

// TODO: after moving global variables, refactoring of this method is required (most likely it will disappear as unnecessary ? logger ?)
func init() {
	var (
		err error
		//ok  bool
	)
	startTime = time.Now()
	ctx := contexts.BuildContext(
		context.Background(),
		contexts.SetContextOperation("00.init"),
		errs.SetDefaultErrsSeverity(errs.SeverityCritical),
		errs.SetDefaultErrsKind(errs.KindInvalidValue),
	)
	args := ProcessArgs(&cfg)
	if err = conf.ReadConfig(args.ConfigPath, &cfg); err != nil {
		logging.LogError(ctx, fmt.Errorf("invalid config: %w", err))
		log.Exit(1)
	}
	if err = logging.Initialize(ctx, cfg.Log.File, cfg.Log.Level, cfg.Log.Format, cfg.TraceFile, currentUser); err != nil {
		logging.LogError(err)
		log.Exit(1)
	}
	// logger is initialized

	// dirs validation
	if cfg.Fingers.InputDir, err = fh.SafeParentResolvePath(cfg.Fingers.InputDir, currentUser, 0700); err != nil {
		logging.LogError(ctx, fmt.Errorf("invalid input dir: %w", err))
		log.Exit(1)
	}
	if err = os.MkdirAll(cfg.Fingers.InputDir, 0755); err != nil {
		logging.LogError(ctx, fmt.Errorf("create input dir [%s] failed: %w", cfg.Fingers.InputDir, err))
		log.Exit(1)
	}

	// dirs validation
	if cfg.Fingers.OutputDir, err = fh.SafeParentResolvePath(cfg.Fingers.OutputDir, currentUser, 0700); err != nil {
		logging.LogError(ctx, fmt.Errorf("invalid output dir: %w", err))
		log.Exit(1)
	}
	if err = os.MkdirAll(cfg.Fingers.OutputDir, 0755); err != nil {
		logging.LogError(ctx, fmt.Errorf("create output dir [%s] failed: %w", cfg.Fingers.OutputDir, err))
		log.Exit(1)
	}

	// dirs validation
	if cfg.Fingers.ReportDir, err = fh.SafeParentResolvePath(cfg.Fingers.ReportDir, currentUser, 0700); err != nil {
		logging.LogError(ctx, fmt.Errorf("invalid report dir: %w", err))
		log.Exit(1)
	}
	if err = os.MkdirAll(cfg.Fingers.ReportDir, 0755); err != nil {
		logging.LogError(ctx, fmt.Errorf("create report dir [%s] failed: %w", cfg.Fingers.ReportDir, err))
		log.Exit(1)
	}
	logging.LogMsg(ctx).Debugf(fmt.Sprintf("cfg: %v", cfg))
}

func main() {
	defer logging.Finalize()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	ctx = contexts.BuildContext(ctx, contexts.SetContextOperation("0.main"))
	logging.LogMsg(ctx).Debug("start listening for signals")
	go func() {
		<-ctx.Done()
		cancel() // stop listening for signed signals asap
		logging.LogMsg(ctx).Debugf("stop listening for signals: %v", ctx.Err())
	}()
	defer cancel() // in case of early return (on error) - signal to close already running goroutines

	// pipeline building
	inputFiles := inputs.NewInputs(
		ctx,
		cfg.Fingers.InputDir,
		cfg.Fingers.Image.Format,
		100,
	)

	inputFileStats := inputs.NewValidator(
		ctx,
		inputFiles.FoundFilePathsCh(),
		100,
	)

	imageLoader := loaders.NewImageLoader(ctx, inputFileStats.ValidatedFileStatCh(), helpers.LoadImage, 100)

	//loader := loader.NewLoader(
	//	ctx,
	//	searcher.FoundFilePathsCh(),
	//	statMetaKeyFunc,
	//	priorDupsFunc,
	//	statValidatorFunc,
	//	cfg.SLinkEnabled,
	//	cfg.PatternFoundFilesInitCapacity,
	//)

	errModerator, err := erf.NewErrorModerator(
		ctx,
		cancel,
		inputFiles,
		inputFileStats,
		imageLoader,
	)
	if err != nil {
		logging.LogError(err)
		return
	}
	pipeline := []workflow.Pipeliner{
		errModerator,
		inputFiles,
		inputFileStats,
		imageLoader,
	}

	// run pipeline
	finish := workflow.Run(ctx, workflow.Pipelines(pipeline).Runners()...)

monitoring:
	for {
		select {
		case <-finish:
			break monitoring
		case <-ctx.Done():
			fmt.Println("\nProcessing stopped")
			<-finish
			return
		case <-time.After(cfg.StatsUpdateRate):
			out.PrintStats(ctx, startTime, workflow.Pipelines(pipeline).StatProducers()...)
		}
	}
	out.PrintStats(ctx, startTime, workflow.Pipelines(pipeline).StatProducers()...)
	//SaveResults(ctx, contentFilter.Stats().(*filtering.ContentFilterStats))
}

//func SaveResults(ctx context.Context, dups *filtering.ContentFilterStats) {
//	reports := out.SaveDupsResults(ctx, cfg.OutputDir, cfg.OutputFilePrefix, cfg.MaxGroupsPerOutputFile, dups)
//	for report := range reports {
//		if report.Err != nil {
//			logging.LogError(report.Err)
//		} else {
//			logging.LogMsg(ctx).Info(
//				fmt.Sprintf("results witten to file [%s]: %d(indexFrom) %d(dups) %d(files) %d(bytes)",
//					report.FileName,
//					report.IndexFrom,
//					report.DupGroupsCount,
//					report.FilesCount,
//					report.Bytes,
//				))
//		}
//	}
//}
