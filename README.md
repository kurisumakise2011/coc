# coc
Count OCcurrences and highlight them

From source

```
go get -u github.com/kurisumakise2011/coc
go build -o coc $GOPATH/src/github.com/kurisumakise2011/coc/wordcount.go
mv $GOPATH/src/github.com/kurisumakise2011/coc/coc /usr/bin

coc -f $GOPATH/src/github.com/kurisumakise2011/coc/randomtext.txt -i "out" -c

```

Then check occurrances

![Example](https://github.com/kurisumakise2011/coc/blob/main/ex.png)
