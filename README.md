HTTPSify <small>v3</small>
=============================
A `Let'sEncrypt` based reverse proxy, that will automatically generate &amp; renew valid `ssl` certs for your domains, it also enables the `http/2` protocol by default, and uses `roundrobin` as an algorithm to loadbalance the incoming requests between multiple `upstreams`

# Install

### # Using Docker
> Just run the following and then have fun !!
```bash
$ docker run -v $HOME:/root/ -p 443:443 alash3al/httpsify
```

### # Building from source
> You must have the `Go` environment installed
```bash
$ go get -u github.com/alash3al/httpsify
```

### Configurations
> Goto your `$HOME` Directory and edit the `hosts.json` to something like this
```json
{
	"example1.com": ["http://localhost"],
	"example2.com": ["http://localhost:8080", "http://localhost:8081"]
}
```
> As you see, the configuration file accepts a `JSON` object/hashmap of `domain` -> `upstreams`,
and yes, it can loadbalance the requests between multiple upstreams using `roundrobin` algorithm.

> Also You don't need to restart the server to reload the configurations, because `httpsify` automatically watches the
configurations file and reload it on any change.
