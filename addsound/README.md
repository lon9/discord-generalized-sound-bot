# addsound

addsound adds sounds from specified directory.

## Usage 

```shell
  Usage of ./addsound:
    -d string
        Distination directory (default "dist")
    -p string
        Prefix of the path to save sounds (default "/sounds_dca")
    -s string
        Source directory (default "src")
```

### Docker

```shell
docker build -t addsound
docker run --rm -v $(pwd)/src:/src -v $(pwd)/dist:/dist addsound addsound [OPTIONS]
```