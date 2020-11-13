# Go-Git-Radar

## TODO

- Add testing

## Setup and Install

Later I will probably provide releases but that hasn't happened yet.

For now, there will be a binary in the repo you can use, with my configs or you
can clone this repo and build it yourself.

You will need to have golang installed in order to build it and you will need to
copy the `config.go.example` to `config.go` and change any config options you
want and then build it.

Build the program with:

```bash
$ make
```

or

```bash
$ make build
```

## Testing

You can run the automated tests with:

```bash
$ make test
```

Every bug fix and added feature should be accompanied with tests that test the
change and all tests should pass before pushing.
