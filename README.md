# FIXDECODER
A convenient and decoder library for FIX messages.

# install
`go get -u github.com/ilovelili/fixdecoder`

# usage
```go
    fd := fixdecoder.NewFixDecoder()
    fd.Decode("<your fix message>")
```

# dependencies
* [gjson](https://github.com/tidwall/gjson)
