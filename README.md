# Calendar Reader for Kindle

A client app that connects to your calendars and serves a simple HTML page with your events for the day. The
page is simple enough for Kindle and other low-performance devices with rudimentary browsers be able to access it.

## Why?

This is a personal project and a rough idea still under refinement. The motivation is that I wanted to put together 
the events from multiple calendars (personal and professional - which are often from different providers like Google 
and Microsoft) in a very simple daily view that I could load into my Kindle on the wall in my office.

I believe this is going to be very helpful especially because I work from home and in this setup it's very common to 
have small personal tasks that you need to take care during the work hours. Similarly, it's common to have professional
responsibilities that happen on my personal time.

I have a home server (which only means "an old computer hidden somewhere and connected to the network") where I can
deploy a very simple application that serves a "Today" page with my schedule. This is very important here because 
without modifying the Kindle operating system, I don't know of a way to access the Google Calendar, for example, since
the web browser is very limited in these devices.

With all that in mind I thought of a very simple web application where:

- I can set up my credentials to access my calendars from any device in the network (one with a modern browser);
- I can access a simple view from my Kindle to see the events of the day

## TODO

Realistically:

- [ ] Proper HTML templates instead of `const`
- [ ] Better UI for the Today page
- [ ] Better UI for the Setup page
- [ ] Make use of the refresh token
- [ ] Select/input calendars to show
- [ ] Support Outlook
- [X] Separate into at least two API files (`/` and `/setup`)
- [X] Move the redirect from `/` to `setup`
- [ ] Remove hardcoded things... environment variables FTW
- [ ] Better service layer to deal with multiple events from multiple providers (includes fixing date-related stuff)

If I can dream (these are prioritized):

1. Unit tests
2. At least Basic auth to access the "Today" page
3. Proper multi-tenant web application
4. Option to have a websocket updating the events automatically
5. ...

## Want to  use this?

Why? But ok, I can give some directions...

### Application Credentials

In the current state, this is a client that you configure to access your (Google only) Calendar events. So you need to set up your own client within Google (Outlook soon, I promise).

Basically, you need to set up a project, credentials and enable the Google Calendar API for it. Follow [this guide from Google](https://developers.google.com/calendar/api/quickstart/go) to get a `credentials.json` file and put it on the root of the project directory.

This is what the file looks like if you followed the tutorial:

```json
{
  "installed": {
    "client_id":"<something-something>",
    "project_id":"<something-something>",
    "auth_uri":"https://accounts.google.com/o/oauth2/auth",
    "token_uri":"https://oauth2.googleapis.com/token",
    "auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs",
    "client_secret":"<something-something>",
    "redirect_uris":[
      "http://localhost/setup"
    ]
  }
}
```

>_NOTE_: You must set up the `redirect_uris` according to your setup of host and port where you're accessing this
> in your own network. `http://localhost:8080`, for example, is fine if running on your personal machine on port 
> `8080`. For a network deployment you might want something like `http://192.168.0.42` or
> `http://something-that-my-local-dns-resolves` and run on port `80`.

>_NOTE_: You must alse append the `/setup` at the end of the host because that's where the app expects to get the 
> redirection back from the provider with the `code` to exchange for an access token.

### Token(s)

This is **not** a multi-tenant application. You currently can set up one Google account to view events from. When you 
hit the `/setup` endpoint (or when the app detects that there are no tokens configured), you get a chance to configure 
a token that will be used to request calendar events on your behalf.

> _NOTE_: your credentials never leave your computer, this is safe to use, I'm not tricking you. But I must say
> that if you don't understand how all this works, you better not use this app at all.  

### Deploy and run

You can use the `make` target `run` to run a docker container. Default port is `8080`, but you can customize it
directly on the `Makefile` or, since you should know what you're doing because you're still reading this, with the 
environment variable `SERVER_PORT`.

Again, you have to make sure you configure the `redirect_uri` in the `credentials.json` file according to the host and
port you will be using to access the application. For example, if this will run on port 8888 on a host that answers by
`my-docker-server`, the value for `redirect_uri` should be:

```
http://my-docker-server:8888/
```

Then you a run the below command on your docker server (or modify the Makefile to use another context or something):

```bash
SERVER_PORT=8888 make run
```

You could also only `make build` and later run the app on a docker server (or under a context configured to use your
local docker server) with:

```bash
docker container run \
        --name kindle-calendar-reader \
 		-p 80:8080 \
 		drugowick.dev/kindle-calendar-reader
```
