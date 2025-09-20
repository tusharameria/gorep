# gorep

## building grep in go

# Flag

```go
go run cmd/main.go [-i] [-r] [-tl] [-tf] [searchQuery] [path1 path2 ... pathn]
```

-i : turns the search case insensetive
-r : turns the search recursive (directories inside directories, default is 1 layer ie if you provide directory as path then it will give result based on files inside that directory only)
-tl : returns total number of lines matched
-tf : returns total number of files matched
searchQuery : your search string
[path1 path2 ... pathn] : array of paths in which you want to search
