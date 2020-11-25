# Go-Git-Radar

This is a HUD for git repos.

This project was inspired heavily on
[michaeldfallen/git-radar](https://github.com/michaeldfallen/git-radar), except
that I have rewritten it in golang instead of using bash scripts.

The display is also almost identical because I like the way that it worked but I
have changed a few things about how the remotes are shown.

There are 2 main reasons for the existence of this project:

1. I was interested in learning Go and needed something to do with it
2. I really like git-radar, but it isn't being actively developed right now and
   hasn't seen any activity in some time.

## Differences between git-radar and go-git-radar

As mentioned, most of the output is the same, but there are some differences.

1. There is built in support for a git config option for your branches called
   `git-radar-parent-remote`. This holds the name of the remote branch that your
   current branch is based off of. go-git-radar uses the parent remote to
   replace the fancy `m` that git-radar showed when remotes diverged. This is to
   give you more information about how your branch differs from your parent. It
   also uses this parent to compare to instead of always using origin/master.
2. Remotes branches are now shown with their remote name, unless the remote name
   is `origin` and then it is just implied, in order to make the prompt shorter.
3. The parent remote is always shown if you are not in a branch that tracks
   `origin/master` now.

See the below `help` output below for some examples.

Note: that your colours may differ, as I am using a gruvbox theme for my
terminal.

Note: that the help output uses your config file, which means that changes you
make to your config file will be reflected in the help output, so you can see
how your formatting or colour changes will affect how certain situations are
displayed.

![](images/help.png)

## Usage

### Show the Help

```bash
$ git-radar help
```

### Run Git Radar

```bash
$ git-radar
```

Note that this won't output anything if you are not in a git repository.

### Run Git Radar With Auto Fetching

```bash
$ git-radar fetch
```

Note that this won't output anything if you are not in a git repository.

### Add git-radar to your prompt

Add something like this to your .bashrc file.

```bash
export PS1+="\$(go-git-radar fetch)"
```

The `fetch` is optional if you do not want auto fetching.

The `\` before the `$(` is important as it escapes the execution of the call.
This means that every time your prompt loads up it will be called and therefore,
updated.

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

### Explanation of config.go

The `config.go` file has several things that you can configure.

#### GIT_RADAR_FETCH_TIME

This is the number of seconds to wait before doing an auto fetch. This variable
has no effect if you run git-radar without `fetch`.

#### PREFIX and SUFFIX

These are strings to prepend and append to the prompt in case you want some
extra padding or something else.

#### COLOUR_PREFIX and COLOUR_SUFFIX

These are not intended to change, but they can be if you so wish. These are
escape codes to start and end a change in colour of text.

These variables are ONLY used inside this config file, so you can choose to
remove them completely as long as you remove all references to them in the
config.go file.

#### BOLD, BLACK, RED etc

These are colour codes that are used in git-radar to colour the text.

These are only used inside the config.go file, so you can choose to remove them
or add more colours as needed so long as you update their usage elsewhere in the
config.go file.

#### REMOTE_AHEAD

This defines a string that is used if your current branch's remote is ahead of
your parent. The `%s` and `%d` are placeholders for data that will be filled in
from the code.

`REMOTE_AHEAD` must exist, but you can change it to be whatever you want. Just
make sure that the placeholders are in the same expected order.

#### REMOTE_BEHIND

This defines a string that is used if your current branch's remote is behind
your parent. The `%s` and `%d` are placeholders for data that will be filled in
from the code.

`REMOTE_BEHIND` must exist, but you can change it to be whatever you want. Just
make sure that the placeholders are in the same expected order.

#### REMOTE_DIVERGED

This defines a string that is used if your current branch's remote is both
behind and ahead of your parent. The `%s` and `%d` are placeholders for data
that will be filled in from the code.

`REMOTE_DIVERGED` must exist, but you can change it to be whatever you want.
Just make sure that the placeholders are in the same expected order.

#### REMOTE_EQUAL

This defines a string that is used if your current branch's remote neither
behind, nor ahead of your parent. The `%s` and `%d` are placeholders for data
that will be filled in from the code.

`REMOTE_EQUAL` must exist, but you can change it to be whatever you want.
Just make sure that the placeholders are in the same expected order.

#### REMOTE_NOT_UPSTREAM

This defines a string that is used if your branch isn't tracking a remote.
The `%s` and `%d` are placeholders for data that will be filled in from the
code.

`REMOTE_NOT_UPSTREAM` must exist, but you can change it to be whatever you want.
Just make sure that the placeholders are in the same expected order.

#### REMOTE_SAME

This defines a string that is used if your branch's remote, and its parent's
remote are the same remote branch (usually only happens if you have your
`master` branch checked out). The `%s` and `%d` are placeholders for data that
will be filled in from the code.

`REMOTE_SAME` must exist, but you can change it to be whatever you want.
Just make sure that the placeholders are in the same expected order.

#### LOCAL_AHEAD

This defines a string that is used if your current local branch is ahead of its
remote tracking branch. The `%s` and `%d` are placeholders for data that will be
filled in from the code.

`LOCAL_AHEAD` must exist, but you can change it to be whatever you want.
Just make sure that the placeholders are in the same expected order.

#### LOCAL_BEHIND

This defines a string that is used if your current local branch is behind its
remote tracking branch. The `%s` and `%d` are placeholders for data that will be
filled in from the code.

`LOCAL_BEHIND` must exist, but you can change it to be whatever you want.
Just make sure that the placeholders are in the same expected order.

#### LOCAL_DIVERGED

This defines a string that is used if your current local branch is both behind
and ahead of its remote tracking branch. The `%s` and `%d` are placeholders for
data that will be filled in from the code.

`LOCAL_DIVERGED` must exist, but you can change it to be whatever you want.
Just make sure that the placeholders are in the same expected order.

#### CHANGES_STAGED

This defines a string that is used to show staged changes. The `%s` and `%d` are
placeholders for data that will be filled in from the code.

`CHANGES_STAGED` must exist, but you can change it to be whatever you want.
Just make sure that the placeholders are in the same expected order.

#### CHANGES_UNSTAGED

This defines a string that is used to show unstaged changes. The `%s` and `%d`
are placeholders for data that will be filled in from the code.

`CHANGES_UNSTAGED` must exist, but you can change it to be whatever you want.
Just make sure that the placeholders are in the same expected order.

#### CHANGES_CONFLICTED

This defines a string that is used to show conflicts. The `%s` and `%d`
are placeholders for data that will be filled in from the code.

`CHANGES_CONFLICTED` must exist, but you can change it to be whatever you want.
Just make sure that the placeholders are in the same expected order.

#### CHANGES_UNTRACKED

This defines a string that is used to show untracked files. The `%s` and `%d`
are placeholders for data that will be filled in from the code.

`CHANGES_UNTRACKED` must exist, but you can change it to be whatever you want.
Just make sure that the placeholders are in the same expected order.

#### STASH_FORMAT

This defines a string that is used to show the current stash state. The `%s` and
`%d` are placeholders for data that will be filled in from the code.

`STASH_FORMAT` must exist, but you can change it to be whatever you want.
Just make sure that the placeholders are in the same expected order.

#### PROMPT_FORMAT

This defines a string that is used for the full prompt that you see on screen.
The `%s` and `%d` are placeholders for data that will be filled in from the
code.

`PROMPT_FORMAT` must exist, but you can change it to be whatever you want.
Just make sure that the placeholders are in the same expected order.

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

## Contributing

* Follow the code styles of the current files
* Comment your code
* Separate functions that perform side effects as much as possible to make
  testing easier
* Every test should be passing after your change
* At least one test should be written to test your change.
* If any tests are failing after your change, the code must be fixed before it
  will be merged (this may just be fixing the test if the expected behaviour of
  the code has changed).

### Testing

You can run the automated tests with:

```bash
$ make test
```

Or you can run automated tests and open up a coverage report in your web
browser:

```bash
$ make test-report
```
