# Comment Timemachine

Integration between Coral Talk app and Pachyderm to create 
a comment moderation tool that traverses time.

# Pre-reqs

You'll need to install / run a pachyderm cluster locally.

[Heres a good place to start](http://docs.pachyderm.io/en/stable/getting_started/local_installation.html)

Once that's complete, this command should work:

```
$ pachctl version
COMPONENT           VERSION
pachctl             1.3.3
pachd               1.3.3
```

# Setup

This inputs the data into Pachyderm PFS, and initializes the submodules (the talk / proxy)

```
make one-time-setup
```

# Run

To run (the talk app, the grpc proxy):

```
make run
```
