punf
====

Upload files/scrots/urls to sr.ht, punpun.xyz, or via ssh.


Help
----

```
Usage: punf [options] [file/url]

options:
  -c,   --clipboard       upload your clipboard as text
  -d,   --desktop         force desktop scrot
  -s,   --selection       upload selection scrot
  -q,   --quiet           disable all feedback (for scripts using punf)
  -h,   --help            print help and exit
```


Dependencies
------------

* fish (2.3.0+)
* getopts (https://github.com/fisherman/getopts)
* xsel (optional, for copying links to clipboard)
* maim (optional, for screenshots)
* slop (optional, for selection screenshot)
* randstr (https://github.com/onodera-punpun/randstr, optional, for ssh uploads)


Installation
------------

Run `make install` inside the `punf` directory to install the script.
`punf` can be uninstalled easily using `make uninstall`.
`punf` can also be run from any directory like a normal script.
Be sure to copy `./configs/config` or `/usr/share/punf/config` to `$HOME/.punf`.

If you use CRUX you can also install using this port: https://github.com/6c37/crux-ports-git/tree/3.2/punf
