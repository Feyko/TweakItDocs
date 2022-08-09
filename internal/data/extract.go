package data

import (
	"TweakItDocs/internal/data/properties"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gammazero/workerpool"
	"github.com/pkg/errors"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
)

var validFilenameRegexp = `^FactoryGame.*/(Build|Desc|Recipe|Schematic)_.*$`
var validFilenameRegexpCompiled = regexp.MustCompile(`^FactoryGame.*/(Build|Desc|Recipe|Schematic)_.*$`)

func ExtractSplit(dir string) ([]Package, error) {
	extractor := newSplitExtractor(dir)
	return extractor.Extract()
}

type splitExtractor struct {
	rootDir        string
	packages       []Package
	pool           *workerpool.WorkerPool
	out            chan *Package
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
		packages:     make([]Package, 0, 25000),
		pool:         workerpool.New(runtime.NumCPU() * 2),
		out:          make(chan *Package),
		doneEvent:    make(chan struct{}),
		stoppedEvent: make(chan struct{}),
		errored:      make(chan error),
		tasks:        0,
		processed:    0,
	}
	return e
}

func (e *splitExtractor) Extract() ([]Package, error) {
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

func extractionTask(path string, errored chan error, out chan *Package) func() {
	return func() {
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

func ExtractAll(data []byte) ([]Package, error) {
	raw, err := extractAllRaw(data)
	if err != nil {
		return nil, err
	}
	extracted := rawRecordsToRecordSlice(raw)
	extracted = FilterForFilename(extracted, validFilenameRegexp)
	resolve(extracted)
	return extracted, nil
}

func ExtractOne(data []byte) (*Package, error) {
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

func resolve(packages []Package) {
	for _, p := range packages {
		p.Resolve()
	}
}

func FilterForFilename(data []Package, match string) []Package {
	r := make([]Package, 0, len(data))
	reg := regexp.MustCompile(match)
	for _, e := range data {
		if reg.MatchString(e.Filename) {
			r = append(r, e)
		}
	}
	return r
}

func extractAllRaw(data []byte) ([]rawRecord, error) {
	var out []rawRecord
	err := json.Unmarshal(data, &out)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal the json: %w", err)
	}
	return out, nil
}

func extractOneRaw(data []byte) (rawRecord, error) {
	var out rawRecord
	err := json.Unmarshal(data, &out)
	if err != nil {
		return rawRecord{}, fmt.Errorf("could not unmarshal the json: %w", err)
	}
	return out, nil
}

func rawRecordsToRecordSlice(raw []rawRecord) []Package {
	return mapSlice(raw, rawRecordToRecord)
}

func mapSlice[T any, R any](s []T, f func(T) R) []R {
	out := make([]R, len(s))
	for i, elem := range s {
		out[i] = f(elem)
	}
	return out
}

func rawRecordToRecord(raw rawRecord) Package {
	return Package{
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
