# Bitly Go

URL Shortener - Short URLs &amp; Custom Free Link Powered by Goland and ... database

## Routes


### `GET /search`

- STRING `q` (required, The minimum length is 1 and does not allow empty queries also we will trim the value to make it ready to search)
We will remove and skip all non-English and non-Digits characters and after that, if `q` is not empty we will search and fetch the results. Otherwise, throw an error!

- POSITIVE INT `limit` (default is 10, Minimum value is 1 and Maximum is 100. Otherwise, throw an error!)

### `GET /expire-soon`

- POSITIVE INT `limit` (default is 10, Minimum value is 1 and Maximum is 100. Otherwise, throw an error!)

### `POST /add`

- STRING `name` (optional, If not defined, we will generate a short and unique random name)

- STRING `link` (required, and we will check the link should be valid and pass URL standard format)
About link value: we must support **UTF-8** characters or query values.

## Database

It's okay to use PostgreSQL or MariaDB.

Note: since this project supports short lifetime links, we can use other databases too. but it's possible to use SQL and delete old rows.
