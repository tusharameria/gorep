# gorep

### grep in go

```go
go run cmd/main.go [-i] [-r] [-time] [-totalLines] [-totalFiles] [-coreWorkers] [-workers]=n [searchQuery] [path1 path2 ... pathn]
```

#### -i : turns the search case insensetive

#### -r : turns the search recursive (directories inside directories, default is 1 layer ie if you provide directory as path then it will give result based on files inside that directory only)

#### -time : prints total time taken in search

#### -totalLines : returns total number of lines matched

#### -totalFiles : returns total number of files matched

#### -coreWorkers : pool workers would be equal to number of your CPU cores

#### -workers : provide the number of workers you want to initialise (default = 1)

```
Note : if -coreWorkers is flagged then number provided with -workers would be irrelevant
```

#### searchQuery : your search string

#### [path1 path2 ... pathn] : array of paths in which you want to search
