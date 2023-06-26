# RbacSimplifier
This project aims to simplify and normalize Role and ClusterRole files.

## Usage

You can address the following command for running the simplifier:

```bash
$ go run main.go --input-file role.yaml
```

The command above will output the final yaml in the standard output.
Note that the `input-file` argument is the path to the file you want to simplify and normalize.

In case you want to store the resulting yaml in a file, the output can be redirected
to a file with:

```bash
$ go run main.go --input-file > output.yaml
```