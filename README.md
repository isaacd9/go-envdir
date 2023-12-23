# The envdir program
This is a Go port of [djb's `envdir` from `daemontools`](https://cr.yp.to/daemontools/envdir.html).

`envdir` runs another program with environment modified according to files in a specified directory.
## Interface
     envdir d child
d is a single argument. child consists of one or more arguments.
`envdir` sets various environment variables as specified by files in the directory named d. It then runs child.

If d contains a file named s whose first line is t, `envdir` removes an environment variable named s if one exists, and then adds an environment variable named s with value t. The name s must not contain =. Spaces and tabs at the end of t are removed. Nulls in t are changed to newlines in the environment variable.

If the file s is completely empty (0 bytes long), `envdir` removes an environment variable named s if one exists, without adding a new variable.

`envdir` exits 111 if it has trouble reading d, if it runs out of memory for environment variables, or if it cannot run child. Otherwise its exit code is the same as that of child.

There are a few existing implemenations of `envdir` in Go such as [https://github.com/yfuruyama](https://github.com/yfuruyama/envdir), [https://github.com/d10n/envdir](https://github.com/d10n/envdir), and [https://github.com/imorph/go-envdir](https://github.com/imorph/go-envdir). This implementation attempts to maintain strict compatibility with djb's `envidir` (any difference in behavior is a bug in this implementation), avoids any external dependencies, and directly calls `execve(2)`, rather than `fork`ing first, avoiding a subprocess and the need to hook up stdin/stdout/stderr.
