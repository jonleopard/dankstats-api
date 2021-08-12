# dankstats-api




#### Project Overview

This is a project aimed improving my skills in both Go and JavaScript. The frontend is built with Next.js (react) that talks to a Go API.


The backend is a simple wrapper for [Helix](https://github.com/nicklaw5/helix). Provided you have twitch api credentials, Helix will take those and will authorize your app by [retrieving](https://github.com/nicklaw5/helix/blob/main/docs/authentication_docs.md#authentication-documentation) an AppAccessToken. This token refreshes every 60 days and thats something I need to account for. Ideally the token get refreshed and swapped out automatically without service interuption. 

My initial thought was to simply wrap the Helix methods in an API handler and expose it in an endpoint. I’ve implemented that for `top-games` and `top-channels`. Howevver, there were two obvious problems that occurred to me when I implemented this:

1. Each page load is going out and grabbing data from the API. You can see how images are loading fairly slowly. I should proxy the request and data should be stored in something like redis so that images load faster.

3. Rate limiting. The Twitch API docs give about 800 points per minute (point being a request). I suspect this would go pretty quickly if the site gained any traffic. I would rather use those 800 points to do some scanning for data collection.

So what I have to do now is think of a way to proxy requests so that users can hit my API as much as they want, without it affecting my API access to Twitch. This also means that API requests are going to retrieve data from my database, rather than straight from the twitch api. 

Twitch does not provide any historical data, so any charting is impossible unless I start archiving the data myself. This means I need to setup a database and try to mock data based on their API. I need to figure out what data I want to collect. Lastly, how does this collection/scanning happen? What triggers it? How often? Being at my level in Go, this seems quite complicated but I am determined to figure it out.

#### Project structure

I'm keeping the project layout as flat as possible until it makes sense not to. When I start to break things up I will likely just add a `pkg` and `cmd` directory. All other files can live in the project root.

#### Developing

To get started developing, all you need at this moment is Go, and Twitch API credentials.

To get running with docker-compose, you'll first need to build the image locally.

`make create && make up`

You can spin everything down by running `make down`
