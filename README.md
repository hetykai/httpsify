HTTPSify <small>v2</small>
=============================
A letsencrypt based reverse proxy to automatically handles ssl termination with no hassle .  
It listens on `:443` by default **DON'T CHANGE IT** because letsencrypt uses that port to confirm ownership of your domains .  

Features
=========
* Auto `SSL Certs` generation and renewal .
* Now you can specify custom backends for custom domains .
* Now serves `websocket` based requestes easily with no problem .

Requirements
=============
* `Golang` >= 1.7

Installation
=============
`go get github.com/alash3al/httpsify`

Usage
=============
> run the following command and you will get some examples at the end .    

`httpsify --help`

Author
========
Mohammed Al Ashaal, a problem solver ;)

