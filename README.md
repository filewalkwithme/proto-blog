# proto-blog

## What is it?
Proto-blog is an experimental blog-engine which let you write posts using the markdown format.

![proto-blog-screenshot-00.png](https://github.com/maiconio/proto-blog/blob/master/screenshots/proto-blog-screenshot-00.png)

## Build

```
go build -o proto-blog
```

## Config
Edit the file `blog.cfg`
```
[blog]
title = The Bootstrap Blog
description = The official example template of creating a blog with Bootstrap.
secret = something-very-secret
session-name = blog-session
theme = bootstrap-blog
port = 8080
database-filename = blog.db

[author]
name = John Doe
username = johndoe
email = your@email.com
```


## Running
```
export blog_password_<username>=<password>
./proto-blog
```

i.e.:
```
export blog_password_johndoe=123456
./proto-blog
```

## Posting

You'll need to login accessing `/admin` and entering your username and password. After that, go back to home and click in `new post`

## More images!
![proto-blog-screenshot-01.png](https://github.com/maiconio/proto-blog/blob/master/screenshots/proto-blog-screenshot-01.png)

![proto-blog-screenshot-02.png](https://github.com/maiconio/proto-blog/blob/master/screenshots/proto-blog-screenshot-02.png)
