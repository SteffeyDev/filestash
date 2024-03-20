# Contributing Guide

Thanks for taking the time to join our community and start contributing. This guide will help you get started with the Filestash project.

## How to contribute?

### Before you submit a pull request

For anything else than a typo or a bug fix, please raise an issue to discuss your proposal before submitting any code.

### License for contributions

As the copyright owner, you agree to license your contributions under an irrevocable MIT license.


### Building from source

*Prerequisites*: Git, Make, Node, Go, Glib 2.0

```
# Download the source
git clone https://github.com/mickael-kerjean/filestash
cd filestash

# Install dependencies
npm install --legacy-peer-deps # frontend dependencies
make build_init # install the required static libraries
mkdir -p ./dist/data/state/
cp -R config ./dist/data/state/

# Create the build
make build_frontend
make build_backend

# Run the program
./dist/filestash
```

### Building from source with Docker
```sh
docker build . -f docker/Dockerfile --platform linux/amd64 --push --tag registry.example.com/filestash:x.x.x
```

### Setup Docker development environment

#### Backend
```
docker run -d -it -v .:/app -p 8334:8334 golang:1.19-bookworm
apt update
apt install -y libvips-dev make libjpeg-dev libtiff-dev libpng-dev libwebp-dev libraw-dev libheif-dev libgif-dev nodejs npm
cd /app
make build_init
mkdir -p /app/data/state
ln -s /app/config /app/data/state/config
npm i --legacy-peer-deps
go run cmd/main.go &
npm run start &
```

### Tests
Our tests aren't open source. This comes as an attempt to restrict opportunistic forks (see [1](https://news.ycombinator.com/item?id=17006902#17009852) and [2](https://www.reddit.com/r/selfhosted/comments/a54axs/annoucing_jellyfin_a_free_software_fork_of_emby/ebk92iu/?utm_source=share&utm_medium=web2x)) from creating a stable release without serious commitment and splitting the community in pieces while I'm on holidays. Also the project welcome serious and willing maintainers.
