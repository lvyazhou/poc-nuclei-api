package nuclei

import (
	"context"

	"github.com/logrusorgru/aurora"
	"github.com/pkg/errors"
	"github.com/projectdiscovery/nuclei/v2/pkg/catalog/config"
	"github.com/projectdiscovery/nuclei/v2/pkg/catalog/loader"
	"github.com/projectdiscovery/nuclei/v2/pkg/core"
	"github.com/projectdiscovery/nuclei/v2/pkg/core/inputs/hybrid"
	"github.com/projectdiscovery/nuclei/v2/pkg/output"
	"github.com/projectdiscovery/nuclei/v2/pkg/parsers"
	"github.com/projectdiscovery/nuclei/v2/pkg/progress"
	"github.com/projectdiscovery/nuclei/v2/pkg/protocols"
	"github.com/projectdiscovery/nuclei/v2/pkg/protocols/common/hosterrorscache"
	"github.com/projectdiscovery/nuclei/v2/pkg/protocols/common/interactsh"
	"github.com/projectdiscovery/nuclei/v2/pkg/protocols/common/protocolinit"
	"github.com/projectdiscovery/nuclei/v2/pkg/templates"
	"github.com/projectdiscovery/nuclei/v2/pkg/types"
	"github.com/projectdiscovery/nuclei/v2/pkg/utils/ratelimit"
)

func NewWrapper(options *types.Options) (*Wrapper, error) {
	w := &Wrapper{}

	if err := w.prepareEngine(options); err != nil {
		return nil, err
	}

	return w, nil
}

type Wrapper struct {
	engine         *core.Engine
	finalTemplates []*templates.Template
}

func (w *Wrapper) prepareEngine(options *types.Options) error {
	options, optionsErr := w.prepareOptions(options)
	if optionsErr != nil {
		return optionsErr
	}

	err := protocolinit.Init(options)
	if err != nil {
		return err
	}

	executerOpts, executerErr := w.prepareExecuterOptions(options)
	if executerErr != nil {
		return executerErr
	}

	// Create engine
	engine := core.New(options)
	engine.SetExecuterOptions(*executerOpts)

	// Load templates and workflows
	workflowLoader, err := parsers.NewLoader(executerOpts)
	if err != nil {
		return errors.Wrap(err, "Could not create loader.")
	}
	executerOpts.WorkflowLoader = workflowLoader

	templateConfig := &config.Config{}

	store, err := loader.New(loader.NewConfig(options, templateConfig, executerOpts.Catalog, *executerOpts))
	if err != nil {
		return errors.Wrap(err, "could not load templates from config")
	}

	store.Load()

	// Cluster templates
	finalTemplates, _ := templates.ClusterTemplates(store.Templates(), engine.ExecuterOptions())
	finalTemplates = append(finalTemplates, store.Workflows()...)

	// Ready to go
	w.engine = engine
	w.finalTemplates = finalTemplates

	return nil
}

func (w *Wrapper) prepareOptions(options *types.Options) (*types.Options, error) {
	options.Stdin = false
	options.NoColor = true
	options.Stream = false
	options.StopAtFirstMatch = true
	options.Silent = true
	options.Headless = false
	options.TemplatesDirectory = "embed_fs/*"
	options.Templates = []string{"embed_fs/templates/*"}
	options.Workflows = []string{"embed_fs/workflows/*"}
	options.Debug = false
	options.FollowHostRedirects = true
	options.Timeout = 30

	return options, nil
}

func (w *Wrapper) prepareExecuterOptions(options *types.Options) (*protocols.ExecuterOptions, error) {
	executerOptions := &protocols.ExecuterOptions{}

	executerOptions.Options = options

	writer, writerErr := NewWriter()
	if writerErr != nil {
		return nil, writerErr
	} else {
		executerOptions.Output = writer
	}

	// TODO: Generate token randomly for interactsh client
	interactshClient, interactErr := w.prepareInteractshClient(writer, "")
	if interactErr != nil {
		return nil, interactErr
	}
	executerOptions.Interactsh = interactshClient

	var cache *hosterrorscache.Cache
	if options.MaxHostError > 0 {
		cache = hosterrorscache.New(options.MaxHostError, hosterrorscache.DefaultMaxHostsCount)
		cache.SetVerbose(options.Verbose)
	}
	executerOptions.HostErrorsCache = nil

	// Creates the progress tracking object
	prog, progressErr := progress.NewStatsTicker(options.StatsInterval, options.EnableProgressBar, options.StatsJSON, options.Metrics, options.MetricsPort)
	if progressErr != nil {
		return nil, progressErr
	}

	executerOptions.Progress = prog
	executerOptions.RateLimiter = ratelimit.NewUnlimited(context.Background())
	executerOptions.ResumeCfg = types.NewResumeCfg()
	executerOptions.StopAtFirstMatch = false
	executerOptions.Catalog = NewCatalog()
	executerOptions.Colorizer = aurora.NewAurora(!options.NoColor)

	workflowLoader, workflowErr := parsers.NewLoader(executerOptions)
	if workflowErr != nil {
		return nil, errors.Wrap(workflowErr, "Could not create loader.")
	}
	executerOptions.WorkflowLoader = workflowLoader

	return executerOptions, nil
}

func (w *Wrapper) prepareInteractshClient(writer *Writer, interactshToken string) (*interactsh.Client, error) {
	opts := interactsh.NewDefaultOptions(writer, nil, nil)
	opts.Debug = false
	opts.NoColor = true
	opts.Authorization = interactshToken
	opts.NoInteractsh = false
	opts.StopAtFirstMatch = false
	opts.DebugRequest = false
	opts.DebugResponse = false

	return interactsh.New(opts)
}

func (w *Wrapper) RunEnumeration(ctx context.Context, urls []string) ([]*output.ResultEvent, error) {
	target, inputErr := w.prepareInput(urls)
	if inputErr != nil {
		return nil, errors.Wrap(inputErr, "could not create input provider")
	}
	defer target.Close()

	// Engine start :)
	_ = w.engine.ExecuteWithOpts(w.finalTemplates, target, true)
	writer := w.engine.ExecuterOptions().Output.(*Writer)
	return writer.ResultEvents, nil
}

func (w *Wrapper) prepareInput(urls []string) (*hybrid.Input, error) {
	options := &types.Options{
		Stream:          false,
		Targets:         urls,
		Stdin:           false,
		TargetsFilePath: "",
	}
	return hybrid.New(options)
}
