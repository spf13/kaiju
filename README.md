Kaiju
=======

[]
[by Cool Hand Mike](http://cool-hand-mike.deviantart.com/art/Saturday-Night-Kaiju-393944841)


Kaiju is the next evolution of conversations on the Internet.

Own your comments! Kaiju provides everything you need to quickly build an online community.

Kaiju is a fast powered open source server that provides real time comment
capabilities to any existing webpage via a small (provided) javascript script.

The javascript is written utilizing browserify and socket.io and is contained
in it’s own space so you can feel free to run it on any page without worrying
about conflicts.


# Installation

## Client

### Using the client

A pre-built version of the Kaiju browser-side client is inside this repository at ```ui/build/kaiju-client.js```.

To use it, pull it into your HTML page using a ```<script>``` tag. You can then instantiate the Kaiju client class
using the following code:

```javascript
var kaiju = new Kaiju({
    url: "KAIJU_SERVER_URL",
    forum: "FORUM_ID",
    page: "PAGE_ID",
    selector: "CSS_SELECTOR_FOR_COMMENTS_SECTION"
});

kaiju.connect();
```

The ```connect()``` method will initiate a live connection to the Kaiju server and start listening for arriving comments.

See ```ui/index.html``` for a usage example.

### Building the client

To build the client, the following pre-requisites must be available on the ```PATH```:

* Node v0.8+
* NPM v1.3+
* GNU Make or any compatible make tool

```bash
cd ui
make
```

This will drop the build into ```ui/build/kaiju-client.js```.

## Server

for now, storage of data is done using a MongoDB server, so you should install one.
It should be possibble to support different storage mechanism via some code enhancements. (PR's welcome!)
