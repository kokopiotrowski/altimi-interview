
# Directory Syncer

This program compares a source directory with a target directory. When run it performs the following actions:
  1. Copies files from source to target if they do not exist in target
  2. Overwrites files in target if they differ from source
     (comparing size, mode and modification date)
  3. (Optionally) Deletes files in target that are missing in source, if
     `--delete-missing` is specified


 All operations are logged to a .txt log file. If an error occurs
 (permissions or other IO issues), the error is logged and the program
 continues to process remaining files.

 ## Flags:
1. `--source` - Specify source directory **(REQUIRED)**
2. `--target` - Specify target directory **(REQUIRED)**
3. `--delete-missing` - Specify if you want to remove files in target directory, not present in source directory *(OPTIONAL)*

 ## Usage Example:

	go run main.go --source /path/to/source --target /path/to/target --delete-missing