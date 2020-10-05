package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"text/template"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

var (
	printHelp    bool
	printVersion bool
)

func init() {
	flag.BoolVar(&printHelp, "help", false, "print usage instructions")
	flag.BoolVar(&printVersion, "version", false, "print version information")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of dhall-docker-img: ...\n")
		fmt.Fprintln(os.Stderr, "OPTIONS:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, usageArgs())
	}
}

func versionString(version, commit, date string) string {
	b := bytes.Buffer{}
	w := tabwriter.NewWriter(&b, 0, 8, 1, ' ', 0)

	fmt.Fprintf(w, "version:\t%s", version)
	fmt.Fprintln(w, "")
	fmt.Fprintf(w, "commit:\t%s", commit)
	fmt.Fprintln(w, "")
	fmt.Fprintf(w, "build date:\t%s", date)
	w.Flush()

	return b.String()
}

func usageArgs() string {
	b := bytes.Buffer{}
	w := tabwriter.NewWriter(&b, 0, 8, 1, ' ', 0)

	fmt.Fprintln(w, "\t<path>\t(required) ")
	w.Flush()

	return fmt.Sprintf("ARGS:\n%s", b.String())
}

type ImageReference struct {
	Registry string
	Name string
	Version string
	Sha256 string
	Key string
}

func processReader(ir io.Reader, imgRefs *[]*ImageReference, seen map[string]struct{}) error {
	contents, err := ioutil.ReadAll(ir)
	if err != nil {
		return err
	}

	matches := NotAnchoredReferenceRegexp.FindAllStringSubmatch(string(contents), -1)

	for _, match := range matches {
		if len(match) != 4 {
			continue
		}

		imgRef := &ImageReference{}

		if strings.HasPrefix(match[3], "sha256:") {
			nameParts := strings.Split(match[1], "/")
			if len(nameParts) > 1 {
				imgRef.Registry = nameParts[0]
				imgRef.Name = strings.Join(nameParts[1:], "/")
			} else {
				imgRef.Name = match[1]
			}
			imgRef.Version = match[2]
			imgRef.Sha256 = strings.TrimPrefix(match[3], "sha256:")

			if strings.HasPrefix(imgRef.Name, "sourcegraph/") {
				imgRef.Key =
					strings.Replace(strings.TrimPrefix(imgRef.Name, "sourcegraph/"), ".", "_", -1)

				if _, ok := seen[imgRef.Key]; !ok {
					*imgRefs = append(*imgRefs, imgRef)
					seen[imgRef.Key] = struct{}{}
				}
			}
		}
	}
	return nil
}

func processFile(path string, imgRefs *[]*ImageReference, seen map[string]struct{}) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	bf := bufio.NewReader(f)

	return processReader(bf, imgRefs, seen)
}

func processInputs(inputs []string, imgRefs *[]*ImageReference, seen map[string]struct{}) error {
	for _, input := range inputs {
		err := filepath.Walk(input, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			if filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml" || filepath.Ext(path) == ".dhall" {
				return processFile(path, imgRefs, seen)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

const imageRecordTemplate = `let images =
{
  {{range $index, $imgRef := .}} {{if gt $index 0}},{{end}} {{$imgRef.Key}} = {
         registry = "{{$imgRef.Registry}}"
         , name = "{{$imgRef.Name}}"
         , version = "{{$imgRef.Version}}"
         , sha256 = "{{$imgRef.Sha256}}"
      }
  {{end}}
}
in images
`

var tmpl = template.Must(template.New("imageRecordDhall").Parse(imageRecordTemplate))

func main() {
	flag.Parse()

	if printHelp {
		flag.Usage()
		os.Exit(0)
	}

	if printVersion {
		output := versionString(version, commit, date)
		fmt.Fprintln(os.Stderr, output)
		os.Exit(0)
	}

	if len(flag.Args()) > 1 {
		flag.Usage()
		os.Exit(1)
	}

	var imgRefs []*ImageReference
	seen := make(map[string]struct{})

	if len(flag.Args()) == 0 {
		err := processReader(os.Stdin, &imgRefs, seen)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := processInputs(flag.Args(), &imgRefs, seen)
		if err != nil {
			log.Fatal(err)
		}
	}

	err := tmpl.Execute(os.Stdout, imgRefs)
	if err != nil {
		log.Fatal(err)
	}
}

