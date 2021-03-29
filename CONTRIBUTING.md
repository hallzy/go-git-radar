# Contributing

If You are going to start working on a reported issue, let someone know.
Leave a comment in the issue, and I can assign the issue to you. This is just
to avoid having multiple people working to solve the same problem.

## Rules

* Always branch from the `dev` branch, unless you are working on a hotfix
  (something that requires immediate attention). Hotfix branches are branched
  from `master`.
    * Hotfix branches should be labelled as such. ex:
      `hotfix/description-of-issue`
    * All other branches should describe what you are doing.
        * If your branch is for working on a reported issue, reference the
          issue number in your commit message. ex: `This resolves #22`
* Follow the code styles of the current files
* Separate functions that perform side effects as much as possible to make
  testing easier
    * A side effect is anything that modifies the state of the program or the
      world, or anything that receives data that isn't a function parameter.
      Ex: print to console, file I/O etc.
    * A general rule of thumb is that if a function doesn't perform a side
      effect, then it's return data will always be exactly the same when the
      function is given the exact same inputs.
* Every test should be passing after your change
* At least one test should be written to test your change.
* If any tests are failing after your change, the code must be fixed before it
  will be merged (this may just be fixing the test if the expected behaviour of
  the code has changed).

## Testing

You can run the automated tests with:

```bash
$ make test
```

Or you can run automated tests and open up a coverage report in your web
browser:

```bash
$ make test-report
```
