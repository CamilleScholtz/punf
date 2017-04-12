package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/go2c/optparse"
	homedir "github.com/mitchellh/go-homedir"
)

func getClipboard() (string, error) {
	fr, err := clipboard.ReadAll()
	if err != nil {
		return "", err
	}

	f := filepath.Join(os.TempDir(), "clipboard.txt")
	if err := ioutil.WriteFile(f, []byte(fr), 0644); err != nil {
		return "", err
	}

	return f, nil
}

func getScrot() (string, error) {
	args := strings.Fields(config.Scrot)

	cmd := exec.Command(args[0], args[1:len(args)]...)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("scrot: Selection cancelled")
	}

	return "/tmp/screenshot.png", nil
}

func getSelScrot() (string, error) {
	args := strings.Fields(config.SelScrot)

	cmd := exec.Command(args[0], args[1:len(args)]...)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("scrot: Selection cancelled")
	}

	return "/tmp/screenshot.png", nil
}

func upload(h string, fl ...string) ([]string, error) {
	args := []string{"--silent"}
	for _, f := range fl {
		args = append(args, "-F", "file=@"+f)
	}

	switch h {
	case "punpun.xyz":
		args = append(args, "-F", "key="+config.Key, "https://punpun.xyz/upload")
	case "sr.ht":
		args = append(args, "-F", "key="+config.Key, "https://sr.ht/api/upload")
	default:
		return []string{}, fmt.Errorf("upload %s: No such host", h)
	}

	b := new(bytes.Buffer)
	cmd := exec.Command("curl", args...)
	cmd.Stdout = b
	if err := cmd.Run(); err != nil {
		return []string{}, fmt.Errorf("upload %s: Something went wrong", fl)
	}

	return strings.Fields(b.String()), nil
}

func main() {
	// Initialize the config.
	if err := initConfig(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Define valid arguments.
	argc := optparse.Bool("clipboard", 'c', false)
	args := optparse.Bool("selection", 's', false)
	argq := optparse.Bool("quiet", 'q', false)
	argh := optparse.Bool("help", 'h', false)

	// Parse arguments.
	vals, err := optparse.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invaild argument, use -h for a list of arguments!")
		os.Exit(1)
	}

	// Print help.
	if *argh {
		fmt.Println("Usage: punf [arguments] [file/url]")
		fmt.Println("")
		fmt.Println("arguments:")
		fmt.Println("  -c,   --clipboard       upload your clipboard as text")
		fmt.Println("  -s,   --selection       upload selection scrot")
		fmt.Println("  -q,   --quiet           disable all feedback")
		fmt.Println("  -q,   --help            print help and exit")
		os.Exit(0)
	}

	var word string
	var fl []string
	switch {
	case *argc:
		word = "clipboard"

		fl[0], err = getClipboard()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer os.Remove(fl[0])
	case *args:
		word = "screenshot"

		fl[0], err = getSelScrot()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer os.Remove(fl[0])
	case len(vals) > 0:
		word = "file"

		for _, v := range vals {
			if _, err := os.Stat(v); os.IsNotExist(err) {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			fl = append(fl, v)
		}
		if len(vals) == 0 {
			os.Exit(1)
		}
	default:
		word = "screenshot"

		fl[0], err = getScrot()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer os.Remove(fl[0])
	}

	url, err := upload(config.Host, fl...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if config.Clipboard {
		// TODO: What should I do with multiple URLs?
		// TODO: Also copy to PRIMARY.
		err := clipboard.WriteAll(url[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
	if config.Log {
		hd, err := homedir.Dir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		lf, err := os.OpenFile(filepath.Join(hd, "punf/log"), os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		for i, u := range url {
			if _, err := lf.WriteString(u + "\t" + fl[i] + "\n"); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}

		lf.Close()
	}
	if config.Notification && !*argq {
		// TODO.
	}
	if config.Print && !*argq {
		if len(url) > 1 {
			fmt.Printf("Punfed files: %s\n", strings.Join(url, ", "))
		} else {
			fmt.Printf("Punfed %s: %s\n", word, url[0])
		}
	}
}
