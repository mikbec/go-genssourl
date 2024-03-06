# go-genssourl

In the current version the __go-genssourl__ service is a web service that
can be used as a backend of a web server (such as Apache HTTPD or NGINX) to
provide a generated redirect URL. This URL includes a hash as a parameter,
which is based on data of a user, whose was authenticated by the web server.
The generate URL may has the following structure:

        <sever_protocol>://<server_host>[:<server_port>]/[<server_path>]?<url_attr_username_key>=username&<url_attr_timestamp_key>=timestamp&<url_attr_hash_key>=hash_val[&<url_attr_id_key>=<id_val>]

__An example:__ If a user has authenticated themself

* as __user1__
* with email address  __user1@my.domain__
* at __2023-11-23 08:15:32 +00:00__

then generated URLs could look like this::

        http://server-one.my.domain/?user=user1&ts=2023-11-23T08%3A15%3A32Z&hash=12e3e5....34bc

or

        https://server-two.my.domain:32146/issue?email=user1%40my.domain&timestamp=2023-11-23T08%3A15%3A32Z&key=12e3e5....34bc&id=id001

For more information please have a look at

* [README.en.md](./docs/README.en.md) (english version)
* [README.de.md](./docs/README.de.md) (german version)

# Compilation

To compile the __genssourl__ binary you can do that by the following steps:

First you need tools like

* __make__
* __go__ compiler with version equal or newer that 1.20

Then run 

```bash
make tidy
make build
```

After that you will find the binary __genssourl__ in __./deploy/bin/__

# Installation

For installation instruction please have a look at [Install.md](./docs/Install.md).
