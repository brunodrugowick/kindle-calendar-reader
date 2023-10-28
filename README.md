# Calendar Reader for Kindle

A client app that connects to your calendars and serves a simple HTML page with your events for the day. The
page is simple enough for Kindle and other low-performance devices with rudimentary browsers be able to access it.

Routes:

- `/`: the root route serves an HTML page with a rudimentary list of events starting from the beginning of the current day
- `/setup`: servers an HTML page with a link to redirect you to Google to authorize the app
- `/json`: for convenience, serves the same list of events in a JSON format (starting from the beginning of the current day)
  - `?startDate=2023-02-14T18:00:00-03:00`: this optional query param is support for even more laziness. I personally integrated into a Gnome Extension to always show my next meeting of the day.

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

For Outlook:

- [ ] Get OAuth to work properly (or continue using DeviceCode grant but better)
- [ ] Understand and properly process the events from the API
- [ ] Actually understand Azure and how that crap works with regard to accessing the events of my personal calendar!

Realistically:

- [ ] Proper HTML templates instead of `const`
- [ ] Better UI for the Today page
- [ ] Better UI for the Setup page
- [X] Make use of the refresh token
- [ ] Select/input calendars to show
- [ ] Support Outlook
- [X] Separate into at least two API files (`/` and `/setup`)
- [X] Move the redirect from `/` to `setup`
- [ ] Remove hardcoded things... environment variables FTW
- [X] Better service layer with a Composite of providers
- [ ] Review API layer and propose improvements
- [X] Fix date-related stuff
- [X] New "next events" endpoint to get an X number of events to come
- [X] Take in more inputs on the JSON requests like `limit`

If I can dream (these are prioritized):

1. Unit tests
2. At least Basic auth to access the "Today" page
3. Proper multi-tenant web application
4. Option to have a websocket updating the events automatically
5. ...

## Want to use this?

Why? But ok, I can give you some directions...

### Application Credentials

In the current state, this is a client that you configure to access your (Google only) Calendar events. So you need to 
set up your own client within Google (Outlook soon, I promise).

Basically, you need to set up a project, credentials and enable the Google Calendar API for it. Follow [this guide from
Google](https://developers.google.com/calendar/api/quickstart/go) to get a `credentials.json` file and put it on the 
root of the project directory.

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
      "http://localhost:8080/setup"
    ]
  }
}
```

>_NOTE_: You must append the `/setup` at the end of the host because that's where the app expects to get the 
> redirection back from the provider with the `code` to exchange for an access token.

### Token(s)

This is **not** a multi-tenant application. You currently can set up one Google account to view events from. When you 
hit the `/setup` endpoint you get a chance to configure a token that will be used to request calendar events on your 
behalf.

> _NOTE_: your credentials are store in the app's memory only, this is safe to use, I'm not tricking you. But I must 
> say that if you don't understand how all this works, you better not use this app at all.  

### Run

You can use the `make` target `run` to run a docker container. Default port is `8080`, but you can customize it
directly on the `Makefile` or, since you should know what you're doing because you're still reading this, with the 
environment variable `SERVER_PORT`.

Here's an example:

```bash
make SERVER_PORT=8888 run
```

### Deploy

Well, since this uses OAuth to Google and Outlook, there's a little inconvenience to make this work on remote hosts. 
I'm not complaining, there's a reason for that. But, if you still are reading and really want to deploy in a remote host
on your network, for example:

>_NOTE_: The following instructions assume you understand what you're doing, so, please, DO NOT CONTINUE if you are not
> sure about anything that I said in the previous sessions. And good day. Bye!

1. After starting the app, go to the `/setup` route.
2. Click the generated link and allow the application on your account (that's going to happen all in Google).
3. Google redirects you to something like `http://localhost:8080/setup?<lots-of-things-here>`.
4. Change `http://localhost:8080` to whereever you deployed your app to, like `http://192.168.0.66` for example.

The app then gets the code that Google generates when you allowed the app in your account and exchange that for a token
that has your credentials. It will use that to get the events from this point on.

I personally use the [Executor Gnome extension](https://extensions.gnome.org/extension/2932/executor/) with the 
following command to show my next event right on the top bar of my Operating System:

```bash
curl -s http://192.168.0.66/json?startDate=$(date --rfc-3339=seconds | tr " " "T") | jq '.[0] | select(.allDay == false) | del(.startTimestamp) | del(.allDay)' | jq -r 'to_entries|map("\(.value|tostring)  ")|.[]'
```

Quick explanation:

- `date --rfc-3339=seconds | tr " " "T"` gets current date in the RFC-3339 format, but without the `T` so I add it with `tr`.
- `jq '.[0]` gets the first entry.
- `select(.allDay == false)` excludes all-day events. 
- `del(.startTimestamp) | del(.allDay)'` deletes these fields from the resulting JSON.
- `jq -r 'to_entries|map("\(.value|tostring)  ")|.[]'` prints only the values without the keys...
  - ... also adds a space after each attribute because Executor removes line breaks.
