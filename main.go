package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/d-tsuji/clipboard"
	"github.com/mxmCherry/multipartbuilder"
	"github.com/tzvetkoff-go/optparse"
)

type punf struct {
	Name   string
	Reader io.Reader
}

func getClipboard() (punf, error) {
	s, err := clipboard.Get()
	if err != nil {
		return punf{}, err
	}

	return punf{"clipboard.txt", strings.NewReader(s)}, nil
}

func getFiles(fnl []string) ([]punf, error) {
	var pl []punf
	for _, fn := range fnl {
		if _, err := os.Stat(fn); os.IsNotExist(err) {
			return []punf{}, err
		}

		f, err := os.Open(fn)
		if err != nil {
			return []punf{}, err
		}

		pl = append(pl, punf{path.Base(fn), f})
	}

	return pl, nil
}

func getScrot(sel bool) (punf, error) {
	args := strings.Fields(config.Scrot)
	if sel {
		args = strings.Fields(config.SelScrot)
	}

	cmd := exec.Command(args[0], args[1:]...)
	if err := cmd.Run(); err != nil {
		return punf{}, err
	}

	f, err := os.Open(filepath.Join(os.TempDir(), "screenshot.png"))
	if err != nil {
		return punf{}, err
	}

	return punf{"screenshot.png", f}, nil
}

func upload(pl []punf) ([]string, error) {
	var urls []string

	mpb := multipartbuilder.New()

	mpb.AddField("user", config.User)
	mpb.AddField("pass", config.Pass)

	for _, p := range pl {
		mpb.AddReader("files[]", p.Name, p.Reader)
	}

	ct, br := mpb.Build()
	defer br.Close()

	res, err := http.Post(config.URL, ct, br)
	if err != nil {
		return urls, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return urls, fmt.Errorf("incorrect HTTP status code %b", res.StatusCode)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return urls, err
	}

	return strings.Split(strings.TrimSuffix(string(b), "\n"), "\n"), nil
}

func view() (string, error) {
	mpb := multipartbuilder.New()

	mpb.AddField("user", config.User)
	mpb.AddField("pass", config.Pass)
	mpb.AddField("function", "view")

	ct, br := mpb.Build()
	defer br.Close()

	res, err := http.Post(config.URL, ct, br)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("incorrect HTTP status code %b", res.StatusCode)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(string(b), "\n"), nil
}

func main() {
	// Initialize the config.
	if err := parseConfig(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Define valid arguments.
	argc := optparse.Bool("clipboard", 'c', false)
	args := optparse.Bool("selection", 's', false)
	argv := optparse.Bool("view", 'v', false)
	argh := optparse.Bool("help", 'h', false)

	// Parse arguments.
	vals, err := optparse.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr,
			"Invaild argument, use -h for a list of arguments!")
		os.Exit(1)
	}

	// Print help.
	if *argh {
		fmt.Println("Usage: punf [arguments] [file/url]")
		fmt.Println("")
		fmt.Println("arguments:")
		fmt.Println("  -c,   --clipboard       upload your clipboard as text")
		fmt.Println("  -s,   --selection       upload selection scrot")
		fmt.Println("  -v,   --view            view all uploaded files")
		fmt.Println("  -h,   --help            print help and exit")
		os.Exit(0)
	}

	std, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var pl []punf
	switch {
	case *argc:
		s, err := getClipboard()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		pl = []punf{s}
	case *args:
		s, err := getScrot(true)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer os.Remove(filepath.Join(os.TempDir(), "screenshot.png"))
		pl = []punf{s}
	case *argv:
		v, err := view()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Println(v)
		os.Exit(0)
	case len(vals) > 0:
		pl, err = getFiles(vals)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case (std.Mode() & os.ModeNamedPipe) != 0:
		pl = []punf{punf{"stdin.txt", os.Stdin}}
	default:
		s, err := getScrot(false)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer os.Remove(filepath.Join(os.TempDir(), "screenshot.png"))
		pl = []punf{s}
	}

	urls, err := upload(pl)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if config.Clipboard {
		// TODO: This doesn't work.
		if err := clipboard.Set(urls[0]); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
	if config.Print {
		fmt.Println(strings.Join(urls, "\n"))
	}
}
