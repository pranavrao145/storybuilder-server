# StoryBuilder Server

Server for the StoryBuilder game. Made in Golang.

# Requirements

At the moment, this server is not deployed in the cloud. However, it can be easily
hosted locally, as long as you have the following (mostly standard) tools:

- One Linux or Mac Device (WSL is untested, though you're welcome to try)
  - The `hostname` and `make` tools must be installed (installed on most Linux and Mac
    systems by default)
- Git ([Install here](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git))
- Docker ([Install here](https://docs.docker.com/get-docker/))
- Docker Compose
  - Mac: comes with Docker Desktop by default, no installation needed.
  - Linux: you may have to install it manually. Refer to [this
    guide](https://docs.docker.com/compose/install/compose-plugin/#installing-compose-on-linux-systems)
    for details.

# Usage

First, if it's not already there, clone this repository to your local machine:

```
git clone https://github.com/pranavrao145/storybuilder-server.git
```

Change into the directory you cloned:

```
cd storybuilder-server
```

Then, to start the server, run:

```
make
```

**Take note of the URL that is printed, as this is what you need to run the CLI.**

To stop the server at any point, change back into the same directory and run:

```
make stop
```
