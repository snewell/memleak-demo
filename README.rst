memleak-demo
============
This is an attempt to create a small example of a memory leak using XDS with the
go version of grpc.  It's designed to be as self-contained as possible.

Upstream bug: https://github.com/grpc/grpc-go/issues/7657


Building
--------
The example is self-contained, and any dependencies are handled by :code:`go.mod`.

.. code-block:: bash

    $ cd /path/to/memleak-demo
    $ go get
    $ go build

The protobufs are pre-built and checked in, but if you want to build them on
your own, everything is in :code:`internal/pb`.

.. code-block:: bash

    $ protoc --proto_path=. --go_out=. --go-grpc_out=. --go_opt=paths=source_relative service.proto

This will create some extra folders, and I don't know enough about :code:`protoc`
to output everything in the directory.  For this repo, I just copied the
generated files and dropped them where they currently live.


Running and Reproducing
-----------------------
Start the server in one terminal, then start the client in another.

.. code-block:: bash

    $ ./memleak-demo server

    $ ./memleak-demo client --host localhost

Arguments are available to change the default port and server hostname (only for
the client).  The client also supports options to increase the nubmer of workers
and how many requests each will do, but I can reproduce this overnight using the
default values for these.  Logging is minimal in both the server and client.

This bug *only* seems to trigger when XDS is used, so that feature cannot be
turned off.  It also requires some traffic on the server, since an idle server
won't have its memory increase.  It also seems to be connected to both connections
and time, since prior testing didn't demonstrate high memory usage after more
severe load testing was stopped.  There's an included script (:code:`run_test.sh`)
that will repeatedly invoke :code:`memleak-demo client` with any extra arugments
passed along, since this forces new connections every few seconds.
