# go-genssourl

In the current version the __go-genssourl__ service is a web service that
can be used as a backend of a web server (such as Apache HTTPD or NGINX) to
provide a generate redirect URL, which is based on data of a user, whose was
authenticated by the web server. This URL has the following structure:

        <sever_protocol>://<server_host>[:<server_port>]/[<server_path>]?<url_attr_username_key>=username&<url_attr_timestamp_key>=timestamp&<url_attr_hash_key>=hash_val[&<url_attr_id_key>=<id_val>]

__An example:__ If a user has autheticated themself

* as __user1__
* with email address  __user1@my.domain__
* at __2023-11-23 08:15:32 +00:00__

then generated URLs could look like this::

        http://server-one.my.domain/?user=user1&ts=2023-11-23T08%3A15%3A32Z&hash=12e3e5....34bc

or

        https://server-two.my.domain:32146/issue?email=user1%40my.domain&timestamp=2023-11-23T08%3A15%3A32Z&key=12e3e5....34bc&id=id001

## Parameters of the URL
The meaning of the individual parameters:

* __server_protocol__
    * The protocol to use for redirection.
    * Here are Meaningful values _http_ or _https_.
    * optional parameter with default value: _https_
* __server_host__
    * name or IP address of destinatio n host
    * optional parameter with default value:  _localhost_
* __server_port__
    * port of the services to connect to on the destination host
    * optional parameter
* __server_path__
    * web path component in URL on the destination host
    * optional parameter
* __url_attr_username_key__
    * name of the parameter, which should contain the value of the username autheticated by our local web proxy
    * optional parameter with default value: _user_
    * format of the value:
        * The value determined by the web server in the authentication process is used.
        * If necessary, URL encoding occurs for characters that are not allowed to be used in a URL. For example, in an email address the __`@`__ character is replaced by the HTML representation __`%40`__.
* __url_attr_timestamp_key__
    * name of the parameter, which should contain the current timestamp
    * optional parameter with default value: _ts_
    * format of the value:
        * The timestamp has the format as a layout value __`2006-01-02T15:04:05Z`__
        * If necessary, URL encoding occurs for characters that are not allowed to be used in a URL.
        * Hint: The Go Layout format consists of parts of the date __`2006-01-02 15:04:05 -0700`__, of the timezone __`MST`__, of the day names __`Mon`__ or __`Monday`__ etc.
* __url_attr_hash_key__
    * name of the parameter, which should contain a hash value, calculated from the account name and the timestamp.
    * optional parameter with default value: _hash_
    * format of the value:
        * Hexdecimal number of the SHA1 value encrypted via RSA from the account name and timestamp
        * The hash value is always calculated.
* __url_attr_id_key__
    * name of the parameter, which should contain an ID key
    * optional parameter with default value: _id_
    * format of the value:
        * a string
        * If this string is the empty string then this parameter is not set.

## The calculation of the hash value
The hash value is calculated as follows:

* A string is created that consists of the account name and the timestamp, which are joined together without separators.
* The SHA1 hash value is calculated from this string.
* This SHA1 hash value is encrypted with a private RSA key.
* This encrypted byte value is converted into a hexadecimal representation.

