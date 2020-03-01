package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go2c/optparse"
	"github.com/zyedidia/clipboard"
	"mvdan.cc/xurls/v2"
)

func curl(fl ...string) (string, error) {
	args := []string{"--silent"}
	for _, f := range fl {
		args = append(args, "-F", "files[]=@"+f)
	}

	args = append(args, "-F", "user="+config.User, "-F", "pass="+config.Pass,
		config.URL)
	cmd := exec.Command("curl", args...)
	b := new(bytes.Buffer)
	cmd.Stdout = b
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("curl %s: Something went wrong", fl)
	}

	return b.String(), nil
}

func getClipboard() ([]string, error) {
	// TODO: What about primary here?
	fr, err := clipboard.ReadAll("clipboard")
	if err != nil {
		return []string{}, err
	}

	f := filepath.Join(os.TempDir(), "clipboard.txt")
	if err := ioutil.WriteFile(f, []byte(fr), 0644); err != nil {
		return []string{}, err
	}

	return []string{f}, nil
}

func getFiles(l []string) ([]string, error) {
	var fl []string
	for _, f := range l {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			return []string{}, err
		}

		fl = append(fl, f)
	}

	return fl, nil
}

func getStdin() ([]string, error) {
	fr, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return []string{}, err
	}

	f := filepath.Join(os.TempDir(), "stdin.txt")
	if err := ioutil.WriteFile(f, fr, 0644); err != nil {
		return []string{}, err
	}

	return []string{f}, nil
}

func getURLs(l []string) ([]string, error) {
	var fl []string
	for _, u := range l {
		f := filepath.Join(os.TempDir(), filepath.Base(u))
		cmd := exec.Command("curl", "-L", "--fail", "--ftp-pasv", "-C", "-",
			"-o", f, u)
		if err := cmd.Run(); err != nil {
			return []string{}, fmt.Errorf(
				"getURLs %s: Could not download source", u)
		}

		fl = append(fl, f)
	}

	return fl, nil
}

func getSelScrot() ([]string, error) {
	args := strings.Fields(config.SelScrot)

	cmd := exec.Command(args[0], args[1:]...)
	if err := cmd.Run(); err != nil {
		return []string{}, fmt.Errorf("scrot: Selection cancelled")
	}

	return []string{filepath.Join(os.TempDir(), "screenshot.png")}, nil
}

func getScrot() ([]string, error) {
	args := strings.Fields(config.Scrot)

	cmd := exec.Command(args[0], args[1:]...)
	if err := cmd.Run(); err != nil {
		return []string{}, fmt.Errorf("scrot: Selection cancelled")
	}

	return []string{filepath.Join(os.TempDir(), "screenshot.png")}, nil
}

func upload(fl ...string) ([]string, error) {
	var urls []string

	url, err := curl(fl...)
	if err != nil {
		return []string{}, err
	}
	urls = strings.Fields(url)

	return urls, nil
}

func view() error {
	args := []string{"--silent", "-F", "function=view", "-F", "user=" + config.
		User, "-F", "pass=" + config.Pass, config.URL}

	cmd := exec.Command("curl", args...)
	cmd.Stdout = os.Stdout
	return cmd.Run()
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

	var fl []string
	switch {
	case *argc:
		fl, err = getClipboard()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer os.Remove(fl[0])
	case *args:
		fl, err = getSelScrot()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer os.Remove(fl[0])
	case *argv:
		if err := view(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case len(vals) > 0:
		urls := xurls.Strict().FindAllString(strings.Join(vals, " "), -1)
		if len(urls) > 0 {
			fl, err = getURLs(urls)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			for _, f := range fl {
				defer os.Remove(f)
			}
		} else {
			fl, err = getFiles(vals)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}
	case (std.Mode() & os.ModeNamedPipe) != 0:
		fl, err = getStdin()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer os.Remove(fl[0])
	default:
		fl, err = getScrot()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer os.Remove(fl[0])
	}

	urls, err := upload(fl...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if config.Clipboard {
		// TODO: What should I do when there are multiple URLs?
		if err := clipboard.WriteAll(urls[0], "clipboard"); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		if err := clipboard.WriteAll(urls[0], "primary"); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
	if config.Print {
		fmt.Println(strings.Join(urls, "\n"))
	}
}
