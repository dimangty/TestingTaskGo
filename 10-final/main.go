package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"module-name/api"
	"module-name/bins"
	"module-name/config"
	localfile "module-name/file"
	"module-name/storage"
)

const storagePath = "bins.json"

type options struct {
	create bool
	update bool
	delete bool
	get    bool
	list   bool
	file   string
	name   string
	id     string
}

type binAPI interface {
	Create(context.Context, []byte, string) (bins.Bin, error)
	Get(context.Context, string) ([]byte, error)
	Update(context.Context, string, []byte) ([]byte, error)
	Delete(context.Context, string) error
}

type inputFiles interface {
	Read(path string) ([]byte, error)
	IsJSON(path string) bool
}

func main() {
	if err := run(context.Background(), os.Args[1:], os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, args []string, output io.Writer) error {
	appOptions, err := parseOptions(args)
	if err != nil {
		return err
	}

	localFiles := localfile.NewLocal()
	binStorage := storage.NewJSONStorage(storagePath, localFiles)
	binList, err := bins.NewBinList(binStorage)
	if err != nil {
		return err
	}

	// Listing is entirely local and therefore must not require an API key.
	if appOptions.list {
		return execute(ctx, appOptions, output, nil, binList, localFiles)
	}

	appConfig, err := config.Load()
	if err != nil {
		return err
	}

	return execute(ctx, appOptions, output, api.New(appConfig), binList, localFiles)
}

func parseOptions(args []string) (options, error) {
	var parsed options
	flags := flag.NewFlagSet("jsonbin", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.BoolVar(&parsed.create, "create", false, "create a bin from a JSON file")
	flags.BoolVar(&parsed.update, "update", false, "update a bin from a JSON file")
	flags.BoolVar(&parsed.delete, "delete", false, "delete a bin")
	flags.BoolVar(&parsed.get, "get", false, "get a bin")
	flags.BoolVar(&parsed.list, "list", false, "list locally saved bins")
	flags.StringVar(&parsed.file, "file", "", "path to a JSON file")
	flags.StringVar(&parsed.name, "name", "", "bin name")
	flags.StringVar(&parsed.id, "id", "", "bin ID")

	if err := flags.Parse(args); err != nil {
		return options{}, fmt.Errorf("parse flags: %w", err)
	}
	if flags.NArg() != 0 {
		return options{}, fmt.Errorf("unexpected positional arguments: %s", strings.Join(flags.Args(), " "))
	}

	actionCount := 0
	for _, selected := range []bool{parsed.create, parsed.update, parsed.delete, parsed.get, parsed.list} {
		if selected {
			actionCount++
		}
	}
	if actionCount != 1 {
		return options{}, errors.New("select exactly one action: --create, --update, --delete, --get, or --list")
	}

	parsed.file = strings.TrimSpace(parsed.file)
	parsed.name = strings.TrimSpace(parsed.name)
	parsed.id = strings.TrimSpace(parsed.id)

	switch {
	case parsed.create:
		if parsed.file == "" || parsed.name == "" {
			return options{}, errors.New("--create requires --file and --name")
		}
	case parsed.update:
		if parsed.file == "" || parsed.id == "" {
			return options{}, errors.New("--update requires --file and --id")
		}
	case parsed.delete, parsed.get:
		if parsed.id == "" {
			return options{}, errors.New("--delete and --get require --id")
		}
	}

	return parsed, nil
}

func execute(
	ctx context.Context,
	appOptions options,
	output io.Writer,
	client binAPI,
	binList *bins.BinList,
	files inputFiles,
) error {
	switch {
	case appOptions.list:
		for _, bin := range binList.Bins() {
			if _, err := fmt.Fprintf(output, "%s\t%s\n", bin.ID, bin.Name); err != nil {
				return fmt.Errorf("print bin list: %w", err)
			}
		}
		return nil

	case appOptions.create:
		document, err := readJSONFile(files, appOptions.file)
		if err != nil {
			return err
		}
		createdBin, err := client.Create(ctx, document, appOptions.name)
		if err != nil {
			return err
		}
		if err := binList.Add(createdBin); err != nil {
			return fmt.Errorf("save created bin locally: %w", err)
		}
		return printValue(output, createdBin)

	case appOptions.update:
		document, err := readJSONFile(files, appOptions.file)
		if err != nil {
			return err
		}
		response, err := client.Update(ctx, appOptions.id, document)
		if err != nil {
			return err
		}
		return printJSON(output, response)

	case appOptions.delete:
		if err := client.Delete(ctx, appOptions.id); err != nil {
			return err
		}
		if err := binList.Delete(appOptions.id); err != nil {
			return fmt.Errorf("delete bin locally: %w", err)
		}
		_, err := fmt.Fprintf(output, "deleted bin %s\n", appOptions.id)
		return err

	case appOptions.get:
		response, err := client.Get(ctx, appOptions.id)
		if err != nil {
			return err
		}
		return printJSON(output, response)
	}

	return errors.New("no action selected")
}

func readJSONFile(files inputFiles, path string) ([]byte, error) {
	if files == nil {
		return nil, errors.New("file dependency is required")
	}
	if !files.IsJSON(path) {
		return nil, fmt.Errorf("input file %q must have a .json extension", path)
	}

	document, err := files.Read(path)
	if err != nil {
		return nil, err
	}
	if !json.Valid(document) {
		return nil, fmt.Errorf("input file %q: %w", path, api.ErrInvalidJSON)
	}
	return document, nil
}

func printJSON(output io.Writer, document []byte) error {
	var formatted bytes.Buffer
	if err := json.Indent(&formatted, document, "", "  "); err != nil {
		return fmt.Errorf("format API response: %w", err)
	}
	formatted.WriteByte('\n')
	if _, err := output.Write(formatted.Bytes()); err != nil {
		return fmt.Errorf("print API response: %w", err)
	}
	return nil
}

func printValue(output io.Writer, value any) error {
	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(value); err != nil {
		return fmt.Errorf("print result: %w", err)
	}
	return nil
}
