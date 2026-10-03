## Rutube downloader

# RU

## Использование

Для использования необходим установленный ffmpeg

Если у вас установлен go, то вы можете использовать утилиту так:

````bash
go run ./cmd/main.go --threads 100 RUTUBE_LINK
````

Или же так, с помощью скомпилированного файла:

````bash
rutube_dwld.exe --threads 100 RUTUBE_LINK
````

Доступны флаги:

```
--threads для указания количества потоков на скачивание
--private для скачивания приватного видео
```

# EN

## Usage

To use this you need to have ffmpeg installed

If you have go installed you can run it like this:

````bash
go run ./cmd/main.go --threads 100 RUTUBE_LINK
````

Or like this with compiled binary

````bash
rutube_dwld.exe --threads 100 RUTUBE_LINK
````

Available flags:

```
--threads for specifying how many concurrent threads to use
--private to download private videos
```
