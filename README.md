skeleton
========

A simple foundation for a web application that provides:
------
* Users and Authentication
* Blog
* Security Middleware (https-redirect, secure headers, csrf prevention, user access control)

Requires a postgres database.

Mailer is stubbed and prints emails to the console. Please extend Mailer with your own email/SMTP code.

Configuration is handled via environment variables in release mode, and with the following default values (DB_DATABASE changes when running tests):

```
  SESSION_COOKIE_NAME = skeleton
  COOKIE_SECRET       = change this secret used to authenticate cookies
  CSRF_SECRET         = change this csrf secret

  DB_ADDR             = localhost:5432
  DB_USER             = postgres
  DB_PASSWORD         = postgres
  DB_DATABASE         = postgres | test

  APP_NAME            = skeleton
  APP_URL             = localhost:8080
```
