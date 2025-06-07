## Rutube downloader

# RU

## Использование

Если у вас установлен go, то вы можете использовать утилиту так:

````bash
go run ./cmd/main.go --threads 100 RUTUBE_LINK
````

Или же так, с помощью скомпилированного файла:

````bash
rutube_dwld.exe --threads 100 RUTUBE_LINK
````

В данный момент доступен только один флаг --threads для указания количества потоков на скачивание

# EN
## Usage

If you have go installed you can run it like this:

````bash
go run ./cmd/main.go --threads 100 RUTUBE_LINK
````

Or like this with compiled binary

````bash
rutube_dwld.exe --threads 100 RUTUBE_LINK
````

Currently only --threads flag is available
