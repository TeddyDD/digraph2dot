# digraph2dot

Converts digraph format to Graphviz.

## Install

```
go install go.teddydd.me/digraph2dot@latest
```

## Usage

```sh
go mod graph | digraph2dot -attr shape=plaintext | dot -Tsvg >out.svg
```

## See also

[digraph](https://pkg.go.dev/golang.org/x/tools@v0.1.10/cmd/digraph), tsort(1)
