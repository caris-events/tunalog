# Tunalog [![](https://img.shields.io/github/v/release/caris-events/tunalog)](https://github.com/caris-events/tunalog/releases) [![](https://img.shields.io/badge/license-MIT-green)](https://github.com/caris-events/tunalog/blob/master/LICENSE)

A no-installation, easy-to-use blogging system written in Go. Powered by [SimpleMDE](https://github.com/sparksuite/simplemde-markdown-editor), [SQLite3](https://www.sqlite.org/). [Gin](https://github.com/gin-gonic/gin), [TocasUI](https://tocas-ui.com/5.0/en-us/index.html).

-   📘 Official Webiste: [tunalog.org](https://tunalog.org) ([Source Code](https://github.com/caris-events/tunalog-docs))
-   📦 Source Code: [GitHub](https://github.com/caris-events/tunalog)
-   🌎 Supported Languages: English (US), 台灣正體

---

![](https://tunalog.org/en-us/assets/screenshot-2_light.png)

&nbsp;

## Getting Started

> Launch with the `-p :8123` argument if you want to change the port.
>
> -   http://localhost:8080 - 🐟 Site
> -   http://localhost:8080/admin - 👩‍💼 Admin Panel

&nbsp;

### Using Docker

```
$ docker run --name tunalog -v tunalog:/var/lib/tunalog -p 8080:8080 yamiodymel/tunalog
```

### Using on Linux

```
$ wget -c https://github.com/caris-events/tunalog/releases/latest/download/tunalog_linux_x64
$ chmod +x tunalog_linux_x64
$ ./tunalog_linux_x64
```

### Using on Windows

1. Download latest Tunalog: [tunalog_windows_x64.exe](https://github.com/caris-events/tunalog/releases/latest/download/tunalog_windows_x64.exe)

2. Double-click the downloaded `tunalog_windows_x64.exe` file to run.

### Using Go

```
$ go install github.com/caris-events/tunalog@latest
```

&nbsp;

## Help

```

████████╗██╗   ██╗███╗   ██╗ █████╗ ██╗      ██████╗  ██████╗
╚══██╔══╝██║   ██║████╗  ██║██╔══██╗██║     ██╔═══██╗██╔════╝
   ██║   ██║   ██║██╔██╗ ██║███████║██║     ██║   ██║██║  ███╗
   ██║   ╚██████╔╝██║ ╚████║██║  ██║███████╗╚██████╔╝╚██████╔╝
   ╚═╝    ╚═════╝ ╚═╝  ╚═══╝╚═╝  ╚═╝╚══════╝ ╚═════╝  ╚═════╝

NAME:
   tunalog - A simple blogging system written in Golang ✨

USAGE:
   tunalog [global options] command [command options]

VERSION:
   1.0.0

COMMANDS:
   reset-password  reset the password of a user, email address is required
   help, h         Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --port value, -p value  port to listen on (default: ":8080")
   --tls-key value         path to TLS key file
   --tls-crt value         path to TLS certificate file
   --help, -h              show help
   --version, -v           print the version
```

&nbsp;

## Contributions

The project was meant to be a simple side project, and may have strong opinions on accepting pull requests. But you can always fork it to make a better version 😎👍

&nbsp;

## Links

-   [twitter/twemoji](https://github.com/twitter/twemoji) - Default fish favicon because I was too lazy to draw one myself.
-   [teacat/tocas](https://github.com/teacat/tocas) - The UI used in the default theme and admin panel.
-   [cznic/sqlite](https://gitlab.com/cznic/sqlite) - For the CGO-free SQLite driver, making it easier to cross-compile Tunalog.
-   [alpinejs/alpine](https://github.com/alpinejs/alpine) - Minimal data-binding for a simpler post creation/edit form.

And Tunalog is inspired by [Ghost Blog](https://ghost.org/) and [WordPress](https://wordpress.org/).

&nbsp;

## ❤️ Love from Taiwan

٩(ˊᗜˋ\*)و Developed by Yami Odymel from <span class="ts-flag is-taiwan-flag is-small"></span> Taiwan, along with the ❤️ from [contributors](https://github.com/caris-events/tunalog/graphs/contributors). The source code is licensed under [MIT](https://github.com/caris-events/tunalog/blob/master/LICENSE). Feel free to use, share, or contribute. Tunalog is co-developed by [Caris Events](https://caris.events), part of [Sorae & Co., Ltd.](https://sorae.co)
