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

### `GET /top`

- POSITIVE INT `limit` (default is 10, Minimum value is 1 and Maximum is 100. Otherwise, throw an error!)

Example response:

```json
{
   "status": true,
   "items": [
      {
         "name": "google",
         "link": "https://google.com",
         "visits": 300
      },
      {
         "name": "github",
         "link": "https://github.com/test",
         "visits": 255
      },
      {
         "name": "fb",
         "link": "https://fb.com",
         "visits": 200
      }
   ]
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

### `POST /add` or `POST /add/`

- STRING `name` (optional, If not defined, we will generate a short and unique random name)

- STRING `link` (required, and we will check the link should be valid and pass URL standard format)
About link value: we must support **UTF-8** characters or query values.

If you send `API-KEY` in the headers, your short link will be alive for ever, otherwise, all links you are creating will only live for 2 days. It should be nice to easily config this limitation inside the source.

Note: you cannot create a duplicate name and It should throw an error. But it's okay to store the same link in different names.

### `POST /:name` or `POST /:name/`

- STRING `link` (required, and we will check the link should be valid and pass URL standard format)
About link value: we must support **UTF-8** characters or query values.

If you send `API-KEY` in the headers, your short link will be alive for ever, otherwise, all links you are creating will only live for 2 days. It should be nice to easily config this limitation inside the source.

Note: you cannot create a duplicate name and It should throw an error. But it's okay to store the same link in different names.

### `GET /:name`

If the name is available on the databases. we will redirect the clients to the target URL. 301 redirect is fine.

Otherwise, we should alert that is a 404 (HTTP Status) route and display a 404 warning.

### `UPDATE /:name`

- STRING `link` (required, and we will check the link should be valid and pass URL standard format)
About link value: we must support **UTF-8** characters or query values.

You can only update a link created by an API-KEY.
So if the owner of that link is same as your API-KEY, you are allowed to update that.

Note: you cannot change the name, you only can change the link value and link that to another new URL.

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

**Important NOTE:** Note that `search`, `add`, `expire-soon`, and `top` are not allowed for names and we should make sure it's not allowed to create such names. since they are already reserved in the router.

The **minimum** allowed name length is 4 and the **maximum** is 25.

The name can only contain English and numeric characters. Any other character is ignored.

## Database

It's okay to use **PostgreSQL** or **MariaDB**.

Note: since this project supports short lifetime links, we can use other databases too. but it's possible to use SQL and delete old rows.
