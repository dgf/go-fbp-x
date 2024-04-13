# Golang Flow Based Programming Experiment

## Resources

Some helpful links about FBP and related topics

- [Flow-Based Programming (Chapters from the First Edition)](https://www.jpaulmorrison.com/fbp/index_old.shtml#book)
- [Tutorial on Flow-Based Programming (Filter File application)](https://github.com/jpaulm/fbp-tutorial-filter-file)

## Development

Run a process benchmark three times

```sh
go test -bench=. ./pkg/process/core/tick_test.go -benchmem -count=3
```
