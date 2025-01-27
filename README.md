# Logdy - terminal logs in web browser

<p align="center">
<img src="https://github.com/logdyhq/logdy-core/assets/1653294/9ec8cb3f-0b8f-4523-b600-377444734b9d" height=300/>
</p>

Supercharge terminal logs with web browser UI and low-code (TypeScript snippets). Save 90% of time searching and browsing logs. It's like **jq, tail, less, grep and awk merged together** and available in a clean UI. **Self-hosted, single binary.**

![autogenerate](https://github.com/logdyhq/logdy-core/assets/1653294/bfe09fa8-bbba-46fa-b54d-503f796c7b57)

Visit [logdy.dev](http://logdy.dev) for more info.

## Download precompiled binary

Naviage to [releases](https://github.com/logdyhq/logdy-core/releases) Github page and download the latest release for your architecture.
In addition you can [add the binary to your PATH](https://logdy.dev/docs/how-tos#how-to-add-logdy-to-path) for easier access.

## Quick start
Below are few examples of what Logdy can do. Whatever the below commands will produce, will be forwarded to a Web UI.
```bash
logdy stdin 'npm run dev'
logdy stdin 'node index.js'
logdy stdin 'go run .'
logdy stdin 'python script.py'
logdy stdin 'tail -f /var/log/nginx/access.log'
```
The following should appear
```
INFO[2024-02...] WebUI started, visit http://localhost:8080    port=8080
```

Open the URL Address and start building parsers, columns and filters.

## Documentation

For product documentation navigate to the [official docs](https://logdy.dev/docs).

## Usage

```bash
Usage:
  logdy [command] [flags]
  logdy [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  demo        Starts a demo mode, random logs will be produced, the [number] defines a number of messages produced per second
  help        Help about any command
  socket      Sets up a port to listen on for incoming log messages. Example ./logdy socket 8233
  stdin       Listens to STDOUT/STDERR of a provided command. Example ./logdy stdin "npm run dev"

Flags:
  -h, --help           help for logdy
  -n, --no-analytics   Opt-out from sending anonymous analytical data that help improve this product
  -p, --port string    Port on which the Web UI will be served (default "8080")
  -v, --verbose        Verbose logs
      --version        version for logdy

Use "logdy [command] --help" for more information about a command.
```

## Development
For development, we recommend running `demo` mode
```bash
go run . demo 1
```

The above command will start Logdy in `demo` mode with 1 log message produced per second.
You can read more about [demo mode](https://logdy.dev/docs/demo-mode).

## Building

For a local architecture build
```bash
go build
```

## Releasing
For a cross architecture build use `gox`. This will generate multiple binaries (in `bin/` dir) for specific architectures, don't forget to update `main.Version` tag.
```bash
gox \
    -ldflags "-X 'main.Version=0.2.0'" \
    -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}" \
    -osarch="linux/amd64 windows/386 windows/amd64 darwin/amd64 darwin/arm64 linux/arm64"
```

Once it's ready, publish the binaries in a new Github release. Again, don't forget to update the version.

```bash
ghr v0.2.0 bin/
```
