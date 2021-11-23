# deduplicate

Simple command line tool for finding and deleting duplicate files.
> Fork of [lionbee/godedupe](https://github.com/lionbee/godedupe)

## Usage

By default `deduplicate` does a dry run.

```bash
./deduplicate [options] directory
  -c Print duplicate values as a CSV to the console
  -d Delete all duplicate values
Calling without any options does a dry run and lists the files to be deleted
```

### Dry run

`./deduplicate directory`

This is the default mode. Files that would be deleted are printer to stdout.

### CSV

`./deduplicate -c directory`

A CSV is printed to stdout in the format: _original_, _duplicate_

### Delete

`./deduplicate -d directory`

Duplicate files are deleted and the name of the deleted files are printed to stdout.

## Logic

`deduplicate` works by walking the supplied directory. Files that have the same number of bytes are MD5 hashed. If there is a hash colission there is potentionally a duplicate file. To confirm the file is in fact a duplicate the files are compared byte by byte before proceding.
