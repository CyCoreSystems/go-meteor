# Meteor utilities for Go

This package contains utilities for working with Meteor applications.

Primarily, this contains the meteor-dockerfilegen CLI tool, which will generate
a multi-stage Dockerfile for a container-deployable Meteor application.

To install this tool, install from a Github [release](https://github.com/CyCoreSystems/go-meteor/releases) or run `go get`:

```
  go get github.com/CyCoreSystems/go-meteor/cmd/meteor-dockerfilegen
```

Simply run `meteor-dockerfilegen` from the base directory of your Meteor
application, and it will return the contents of the Dockerfile.


For instance:

```
  meteor-dockerfilegen -port 3000 > Dockerfile
```

You can run `meteor-dockerfilegen --help` for help on the various options.

If you have any ideas for making this tool more flexbile, please feel free to
open an issue or pull request!

