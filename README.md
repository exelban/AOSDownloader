# AOSDownloader

[Apple OpenSource](https://opensource.apple.com) download tool

## Install
```bash
go get -u github.com/exelban/AOSDownloader
```

## Usage
```bash
AOSDownloader https://opensource.apple.com/source/top/top-125
AOSDownloader https://opensource.apple.com/source/top/top-125 ./top

AOSDownloader https://opensource.apple.com/source/top/top-125 --out ./top
```

### Parameters

**Long** | **Short** | **Type** | **Description**
--- | --- | --- | ---
--url | -u | string | url to project which you want to download
--out | -o | string | destination path for project
--debug | -d  | bool | debug mode

#### Positional parameters
```bash
AOSDownloader [url] [out]
```

## License
[MIT License](https://github.com/exelban/AOSDownloader/blob/master/LICENSE)
