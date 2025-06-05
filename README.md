# scriptura

A CLI app for reading the Bible by book, chapter(s), or verse(s); uses the Douay-Rheims 1899 American Edition (DRA) version of the Bible.

## Usage

This usage information can be found by running `scriptura --help`:

```
Usage: scriptura <book>
       scriptura <book> [start_chapter]-[end_chapter]
       scriptura <book> <chapter>
       scriptura <book> <chapter>:[start_verse]-[end_verse]
       scriptura <book> <chapter>:<verse>

Options:
       --books    Print the list of books and exit
       --help     Print this message and exit
       --license  Print license, citation information and exit
       --version  Print version and exit
```

## Install

### Binary releases

Binaries for Windows and Linux on x86_64 are [available](https://github.com/john-bieren/scriptura/releases).

### Build from source

1. Clone the repository:
    ```
    git clone https://github.com/john-bieren/scriptura.git
    ```
2. Navigate to the project directory and compile the program:
    ```
    go build
    ```
