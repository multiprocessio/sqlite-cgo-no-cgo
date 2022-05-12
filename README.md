# SQLite: CGO vs no CGO

This repo benchmarks mattn/go-sqlite3 against modernc.org/sqlite which
is a translation of SQLite3 from C to Go. This translation allows the
latter package to avoid CGO since there is no C.

My initial observations showed it being twice as slow as
mattn/go-sqlite3 and this repo is to test that observation.

See the [blog post for details](https://datastation.multiprocess.io/blog/2022-05-12-sqlite-in-go-with-and-without-cgo.html).
