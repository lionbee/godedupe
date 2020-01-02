# godedupe

Simple command line tool for finding and deleting duplicate files.

## Usage

By default `godedupe` does a dry run.

```bash
./godedupe [options] directory
  -c Print duplicate values as a CSV to the console
  -d Delete all duplicate values
Calling without any options does a dry run and lists the files to be deleted
```

### Dry run

`./godedupe directory`

This is the default mode. Files that would be deleted are printer to stdout.

### CSV

`./godedupe -c directory`

A CSV is printed to stdout in the format: _original_, _duplicate_

### Delete

`./godedupe -d directory`

Duplicate files are deleted and the name of the deleted files are printed to stdout.

## Logic

`godedupe` works by walking the supplied directory. Each file that is encountered is MD5 hashed. The hash is stored in a map of hash to path. If there is a hash colission there is potentionally a duplicate file. To confirm the file is in fact a duplicate the files are compared byte by byte before proceding.

If a duplicate is found, it becomes the new original.

### Example

If file1.jpg, file2.jpg, file3.jpg are all the same and you choose the output to be CSV then the output is

> file1.jpg, file2.jpg\
> file2.jpg, file3.jpg

Indicating that file2.jpg is a duplicate of file1.jpg and file3.jpg is a duplicate of file2.jpg.

file2.jpg and file3.jpg are the files that would be deleted.
