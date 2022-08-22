# Bitly Go

URL Shortener - Short URLs &amp; Custom Free Link Powered by Goland and ... database

## Routes


### `GET /search`

### `GET /expire-soon`

POSITIVE INT `limit` (default is 10)

### `POST /add`

STRING `name` (optional, If not defined, we will generate a short and unique random name)

STRING `link`

## Database

It's okay to use PostgreSQL or MariaDB.

Note: since this project supports short lifetime links, we can use other databases too. but it's possible to use SQL and delete old rows.
