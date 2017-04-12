[![Go Report Card](https://goreportcard.com/badge/github.com/onodera-punpun/punf)](https://goreportcard.com/report/github.com/onodera-punpun/punf)

Upload files/scrots/urls to sr.ht, punpun.xyz, or via ssh.


## SYNOPSIS

punf [arguments] [file/url]


## EXAMPLES

Upload stdin as a text file:
```
$ cat Pkgfile | punf
Punfed stdin: https://punpun.xyz/BMip.txt
```

Download URL, and upload it:
```
$ punf https://i.4cdn.org/g/1450659832892.png
Punfed url: https://punpun.xyz/6r2T.png
```


## AUTHORS

Camille Scholtz
