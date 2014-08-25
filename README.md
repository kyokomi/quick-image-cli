# quick-image-cli

[![Build Status](https://drone.io/github.com/kyokomi/quick-image-cli/status.png)](https://drone.io/github.com/kyokomi/quick-image-cli/latest)
[![Coverage Status](https://img.shields.io/coveralls/kyokomi/quick-image-cli.svg)](https://coveralls.io/r/kyokomi/quick-image-cli?branch=master)

===============

terminal tool to upload quickly and easily image for golang（Go）

## Usage

```sh
$ quick-image-cli
NAME:
   quick-image-cli - terminal tool to upload quickly and easily image

USAGE:
   quick-image-cli [global options] command [command options] [arguments...]

VERSION:
   0.2.0

AUTHOR:
  kyokomi - <kyoko1220adword@gmail.com>

COMMANDS:
   add
   list
   delete-config
   help, h		Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h		show help
   --version, -v	print the version
```

## Install

```sh
$ brew tap kyokomi/homebrew-quick-image-cli
$ brew install quick-image-cli
```

### AccessToken

coming soon ...

## Demo

### list

```sh
$ quick-image-cli list
![/go.jpg](https://dl.dropboxusercontent.com/1/view/1i6rfuz10yxzsyf/%E3%82%A2%E3%83%97%E3%83%AA/kyokomi-sample/go.jpg)
![/gopher.png](https://dl.dropboxusercontent.com/1/view/b0ypv71rvg47kae/%E3%82%A2%E3%83%97%E3%83%AA/kyokomi-sample/gopher.png)
![/img_10.jpeg](https://dl.dropboxusercontent.com/1/view/6jnhs3gi77hex2b/%E3%82%A2%E3%83%97%E3%83%AA/kyokomi-sample/img_10.jpeg)
![/logo.png](https://dl.dropboxusercontent.com/1/view/g74s2wqu7kzt601/%E3%82%A2%E3%83%97%E3%83%AA/kyokomi-sample/logo.png)
```

```sh
$ alias quick-image='$(quick-image-cli list | peco | gocopy)'
$ quick-image
```

- [clipboard](https://github.com/atotto/clipboard)
- [peco](https://github.com/peco/peco)

### add

```sh
$ 
```


## Contribution
 
```sh
$ gox -osarch="darwin/amd64" -output="_obj/quick-image-cli" ./
$ zip _obj/quick-image-cli.zip _obj/quick-image-cli
$ rm _obj/quick-image-cli
$ ghr -u kyokomi -r quick-image-cli {tag} _obj/
```

## Lisence

[MIT](https://github.com/kyokomi/quick-image-cli/blob/master/LICENSE)

## Author

[kyokomi](https://github.com/kyokomi)

