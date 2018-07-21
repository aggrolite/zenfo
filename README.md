# zenfo.info

zenfo.info crawls Zen center websites for event information.

Events are stored into a Postgres database, and served as JSON via HTTP.

## Sources

So far there are two sources:

1. San Francisco Zen Center (sfzc.org)
2. Angel City Zen Center (aczc.org)

The goal for a 1.0 release is to have crawlers for most of Bay Area Zen centers, and possibly one more around Los Angeles.

Sources are scraped by using the `Worker` interface. All workers share a common HTTP client which provides a custom user agent. And in the future there will need to be rate limiting per site.

## Frontend

No frontend exists yet, but that is also a required goal for 1.0 release. The frontend will make use of the JSON API, which makes queries to the Postgres database.

## Build / Install

For now, simply run `make` to try it out. A Postgres instance is required to be running in the background. I am developing on macOS and vendoring is tied to it, but when I get closer to a 1.0 release my plan is to tune for Linux / setup Docker.
