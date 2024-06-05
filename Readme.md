# Bobbox

A simple utility that listen for file changes and stores that information on a hidden file specified with --metadata_file.
Handles delete, create, update and write operation on the given folder and subfolders. 
Later we could implement an image as tree graph with that information.


## Motivation 
Working with fsnotify.


# 
```json
{
    "/Users/emdev/Downloads/coding/bobbox/test_folder/": {
        "name": ".",
        "full_path": "/Users/emdev/Downloads/coding/bobbox/test_folder/",
        "size": -1, // -1 indicates a directory
        "file_type": "dir"
    },
    "/Users/emdev/Downloads/coding/bobbox/test_folder/els": {
        "name": "els",
        "full_path": "/Users/emdev/Downloads/coding/bobbox/test_folder/els",
        "size": -1,
        "file_type": "dir"
    },
    "/Users/emdev/Downloads/coding/bobbox/test_folder/els/poil": {
        "name": "els/poil",
        "full_path": "/Users/emdev/Downloads/coding/bobbox/test_folder/els/poil",
        "size": 0,
        "file_type": "file"
    },
    "/Users/emdev/Downloads/coding/bobbox/test_folder/pols": {
        "name": "pols",
        "full_path": "/Users/emdev/Downloads/coding/bobbox/test_folder/pols",
        "size": -1,
        "file_type": "dir"
    },
    "/Users/emdev/Downloads/coding/bobbox/test_folder/pols/123": {
        "name": "pols/123",
        "full_path": "/Users/emdev/Downloads/coding/bobbox/test_folder/pols/123",
        "size": -1,
        "file_type": "dir"
    },
    "/Users/emdev/Downloads/coding/bobbox/test_folder/pols/123/ddd.txt": {
        "name": "ddd.txt",
        "full_path": "/Users/emdev/Downloads/coding/bobbox/test_folder/pols/123/ddd.txt",
        "size": 5, // size in bytes
        "file_type": "file"
    },
    "/Users/emdev/Downloads/coding/bobbox/test_folder/pols/123/sss": {
        "name": "pols/123/sss",
        "full_path": "/Users/emdev/Downloads/coding/bobbox/test_folder/pols/123/sss",
        "size": 0,
        "file_type": "file"
    }
}
```

## Usage

```sh 
bobbox watch --path ./user/downloads --metadata_file ./bobbox # creates a local hidden file with path metadata 
```
