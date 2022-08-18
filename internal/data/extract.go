package data

import (
	"TweakItDocs/internal/data/properties"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gammazero/workerpool"
	"github.com/pkg/errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var validFilenameRegexp = `^FactoryGame.*/(Build|Desc|Recipe|Schematic)_.*$`

const clean = false

//var validFilenameRegexpCompiled = regexp.MustCompile(`^FactoryGame.*/(Build|Desc|Recipe|Schematic)_.*$`)
var validFilenameRegexpCompiled = regexp.MustCompile(`.`)

func ExtractAllFromDir(dir string) ([]Asset, error) {
	extractor := newSplitExtractor(dir)
	return extractor.Extract()
}

type splitExtractor struct {
	rootDir        string
	packages       []Asset
	pool           *workerpool.WorkerPool
	out            chan *Asset
	outEvent       chan struct{}
	doneEvent      chan struct{}
	stoppedEvent   chan struct{}
	errored        chan error
	tasks          int
	processed      int
	waitingForLast bool
}

func newSplitExtractor(path string) splitExtractor {
	e := splitExtractor{
		rootDir:      path,
		packages:     make([]Asset, 0, 25000),
		pool:         workerpool.New(runtime.NumCPU() * 2),
		out:          make(chan *Asset),
		doneEvent:    make(chan struct{}),
		stoppedEvent: make(chan struct{}),
		errored:      make(chan error),
		tasks:        0,
		processed:    0,
	}
	return e
}

func (e *splitExtractor) Extract() ([]Asset, error) {
	e.startReceiver()
	err := e.startExtractionTasks()
	if err != nil {
		return nil, errors.Wrap(err, "could not start the extraction tasks")
	}

	select {
	case err = <-e.errored:
	case <-e.doneEvent:
	}
	return e.packages, err
}

func (e *splitExtractor) startExtractionTasks() error {
	ctx, cancel := context.WithCancel(context.Background())
	e.pool.Pause(ctx)
	err := filepath.WalkDir(e.rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		e.pool.Submit(extractionTask(path, e.errored, e.out))
		e.tasks++
		return nil
	})
	if err != nil {
		cancel()
		return err
	}
	cancel()
	return nil
}

func extractionTask(path string, errored chan error, out chan *Asset) func() {
	return func() {
		if clean {
			cleanData(path)
		}
		b, err := os.ReadFile(path)
		if err != nil {
			errored <- errors.Wrap(err, "could not read file")
		}

		extracted, err := ExtractOne(b)
		if err != nil {
			errored <- errors.Wrap(err, "could not decode data")
		}
		out <- extracted
	}
}

func (e *splitExtractor) startReceiver() {
	go func() {
		for {
			pkg := <-e.out
			e.processed++
			if pkg != nil {
				e.packages = append(e.packages, *pkg)
			}
			if e.processed == e.tasks {
				e.doneEvent <- struct{}{}
			}
		}
	}()
}

func ExtractOne(data []byte) (*Asset, error) {
	raw, err := extractOneRaw(data)
	if err != nil {
		return nil, err
	}
	extracted := rawRecordToRecord(raw)
	if !isValidFilename(extracted.Filename) {
		return nil, nil
	}
	extracted.Resolve()
	return &extracted, nil
}

func isValidFilename(filename string) bool {
	return validFilenameRegexpCompiled.MatchString(filename)
}

func FilterForFilename(data []Asset, match string) []Asset {
	r := make([]Asset, 0, len(data))
	reg := regexp.MustCompile(match)
	for _, e := range data {
		if reg.MatchString(e.Filename) {
			r = append(r, e)
		}
	}
	return r
}

func extractOneRaw(data []byte) (rawRecord, error) {
	var out rawRecord
	err := json.Unmarshal(data, &out)
	if err != nil {
		return rawRecord{}, fmt.Errorf("could not unmarshal the json: %w", err)
	}
	return out, nil
}

func mapSlice[T any, R any](s []T, f func(T) R) []R {
	out := make([]R, len(s))
	for i, elem := range s {
		out[i] = f(elem)
	}
	return out
}

func rawRecordToRecord(raw rawRecord) Asset {
	return Asset{
		Filename: raw.ExportRecord.FileName,
		Exports:  rawExportsToExportSlice(raw.Exports),
		Imports:  raw.Summary.Imports,
	}
}

func rawExportsToExportSlice(raw []rawExport) []Export {
	return mapSlice(raw, rawExportToExport)
}

func rawExportToExport(raw rawExport) Export {
	return Export{
		ClassIndex:    raw.Export.ClassIndex.Convert(),
		SuperIndex:    raw.Export.SuperIndex.Convert(),
		TemplateIndex: raw.Export.TemplateIndex.Convert(),
		OuterIndex:    raw.Export.OuterIndex.Convert(),
		ObjectName:    raw.Export.ObjectName,
		Properties:    rawPropertiesToPropertySlice(raw.Data.Properties),
	}
}

func rawPropertiesToPropertySlice(raw []properties.RawProperty) []properties.Property {
	return mapSlice(raw, properties.DataToProperty)
}

func cleanData(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not read"))
	}
	data = []byte(strings.ReplaceAll(string(data), "\\u0000", ""))
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not write"))
	}
}
