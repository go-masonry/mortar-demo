# Tutorial - Part 7 Makefile and LDFLAGS

Previously we mentioned that by calling <http://localhost:5382/self/build> you can see your service build information.
To see them you need to inject them during build.

Here you can see in the provided [`Makefile`](Makefile) how it's done. Now you can run your service with

```s
make run
```

And should see something similar to

```s
HTTP/1.1 200 OK
Content-Length: 200
Content-Type: text/plain; charset=utf-8
Date: Thu, 13 Aug 2020 10:54:47 GMT

{
    "build_tag": "42",
    "build_time": "2020-08-13T10:54:41Z",
    "git_commit": "1f034c1",
    "hostname": "Tals-Mac-mini.lan",
    "init_time": "2020-08-13T13:54:43.044148+03:00",
    "up_time": "4.410331293s",
    "version": "v1.2.3"
}
```

> There are also API generate task and Test task...
