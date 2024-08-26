# my-ls

## Overview

`my-ls` is a custom implementation of the Unix `ls` command written in Go. This project replicates the behavior of the original `ls` command with the following variations:
- Supports the `-l`, `-R`, `-a`, `-r`, and `-t` flags.
- Displays files and folders in the specified directory, or in the current directory if none is specified.
- The `ls -l` display format is identical to the system command.

<br>

## Features

- **-l**: Long listing format.
- **-R**: Recursively list subdirectories.
- **-a**: Include directory entries whose names begin with a dot (`.`).
- **-r**: Reverse order while sorting.
- **-t**: Sort by modification time, newest first.

<br>

## Installation
To install my-ls, clone this repository and build the project using Go:
```bash
git clone https://github.com/yourusername/my-ls.git
cd my-ls
go build
```

<br>

## Usage

To use `my-ls`, simply run it from the terminal:

```bash
my-ls [options] [directory]
```

<br>

## Examples
- List files and directories in the current directory:
```
my-ls
```

- List files and directories in the specified directory:
```
my-ls /path/to/directory
```

- my-ls /path/to/directory
```
my-ls -l -a -t
```

<br>


