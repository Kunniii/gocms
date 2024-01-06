# GoCMS - Headless CMS written in Go

A headless CMS with PostgreSQL and Golang!

## How to use?

First, create a `.env` file for the configurations

```sh
cp sample.env .env
```

Then edit the `.env` file with your own information.

### With docker

You can download the image from GHCR

```sh
docker pull ghcr.io/kunniii/gocms:latest
```

By default the server will listen on port 3000

All configurations is available in [`compose.yaml`](https://github.com/Kunniii/gocms/blob/main/compose.yaml) file, or you can find it in [`sample.env`](https://github.com/Kunniii/gocms/blob/main/sample.env)

### Build from source

Make sure you have Go version 1.21 or above.

The source code is available in the `src/` folder
