# Bitly Go

URL Shortener - Short URLs &amp; Custom Free Link Powered by Goland and ... database

## Routes

### `GET /`

This route will show this README to show the features of the projects to everyone.

### `GET /search`

- STRING `q` (required, The minimum length is 1 and does not allow empty queries also we will trim the value to make it ready to search)
We will remove and skip all non-English and non-Digits characters and after that, if `q` is not empty we will search and fetch the results. Otherwise, throw an error!

- POSITIVE INT `limit` (default is 10, Minimum value is 1 and Maximum is 100. Otherwise, throw an error!)

Example response:

```json
{
   "status": true,
   "items": {
    "google": "https://google.com",
    "facebook": "https://fb.com/?from_my_site",
    "telegram": "https://t.me/"
   }
}
```

OR

```json
{
   "status": false,
   "error": "Error message"
}
```


### `GET /expire-soon`

- POSITIVE INT `limit` (default is 10, Minimum value is 1 and Maximum is 100. Otherwise, throw an error!)

Example response:

```json
{
   "status": true,
   "items": {
    "google": "https://google.com",
    "facebook": "https://fb.com/?from_my_site",
    "telegram": "https://t.me/"
   }
}
```

OR

```json
{
   "status": false,
   "error": "Error message"
}
```

### `POST /add`

- STRING `name` (optional, If not defined, we will generate a short and unique random name)

- STRING `link` (required, and we will check the link should be valid and pass URL standard format)
About link value: we must support **UTF-8** characters or query values.

If you send `API-KEY` in the headers, your short link will be alive for ever, otherwise, all links you are creating will only live for 2 days. It should be nice to easily config this limitation inside the source.

Note: you cannot create a duplicate name and It should throw an error. But it's okay to store the same link in different names.


### `DELETE /:name` or `DELETE /:name/`

You can only delete links created by an API-KEY.
So if the owner of that link is same as your API-KEY, you are allowed to delete.

Example response:

```json
{
   "status": true
}
```

OR

```json
{
   "status": false,
   "error": "Sorry, no permission"
}
```

## Database

It's okay to use **PostgreSQL** or **MariaDB**.

Note: since this project supports short lifetime links, we can use other databases too. but it's possible to use SQL and delete old rows.
