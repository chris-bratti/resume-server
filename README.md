# Resume-Server

## Usage

This is the public repo for my resume website.
There are two ways to run this locally if you would like. Docker is the recommended method.
Note that the app assumes two directories are present:

- /env
    - This directory holds the `.env` environment file. This file is used to store sensitive info that the app pulls from. Find an example `.env` below, or you can pull this repo down and run the `make-configs.sh` file to generate one.
- /media
    - This directory holds all of the images for the site itself. If running locally, you will obviously not have these. But hey, you could always use your own!

Example `.env` file:

```
PORT=8080

EMAIL=to_email_address@gmail.com

FROM_EMAIL=smtp_address@gmail.com

KEY=gmail smtp code here

PHONE=(123) 456 7890
```

### Docker
With either Docker method: to map to a different port locally (say port 1234), change the `8080:8080` to `1234:8080`
#### Docker run

```
docker run -d --name "resume-server" -p 8080:8080 -v /path/to/env:/app/env -v /path/to/media:/app/media rhysbratti/resume-server:latest
```

#### Docker compose

```
version: '3.3'

services:
  resume-server:
    image: rhysbratti/resume-server:latest
    container_name: resume-server
    volumes:
      - /path/to/env:/app/env
      - /path/to/media:/app/media
    ports:
      - 8080:8080
    restart: unless-stopped
```
### Locally

You can run it locally if you would like, though it might be a bit more hassle. Follow these steps to do so:
- [Install go](https://go.dev/doc/install)
- `git clone https://github.com/rhysbratti/resume-server.git`
- Run `make-configs.sh` script
    - This will generate the skeleton `.env` file as well as the /media directory
- `go build go.main`
- `./main`


