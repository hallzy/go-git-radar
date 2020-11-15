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

### Recommended Git Alias

Git Radar uses a custom git configuration to keep track of what the parent
remote branch of your current branch is. This is so that Git Radar can tell you
when your current remote branch is ahead or falls behind the parent.

When I refer to a parent, I mean the branch that you branched from in order to
make the branch that you are in currently. The presumption here is that
eventually you will be merging back into that branch later, so knowing how far
ahead or behind you are from it would help.

This is a Git alias I created which automatically sets the git config option
when you create a branch.

```bash
  cob = "!f() { \
      currentTracking=\"$(git for-each-ref --format='%(upstream:short)' \"$(git symbolic-ref -q HEAD)\")\"; \
      if [ -z \"$currentTracking\" ]; then \
        echo \"Could not determine tracking info for current branch.\"; \
        return; \
      fi; \
      git checkout -b \"$1\"; \
      git config --local branch.\"$(git rev-parse --abbrev-ref HEAD)\".git-radar-parent-remote \"$currentTracking\"; \
    }; \
    f"
```

Now, when you want to create a new branch just do:

```bash
$ git cob <new-branch-name>
```

This will create a new branch and check it out for you, but will also set the
configuration in your local gitconfig file.

If you don't set this git config variable, it will always compare you to
origin/master.

## Testing

You can run the automated tests with:

```bash
$ make test
```

Every bug fix and added feature should be accompanied with tests that test the
change and all tests should pass before pushing.
