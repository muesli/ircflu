ircflu
======

ircflu is an IRC bot written in Go with a flexible message- & command-handler.
Among its advanced features are a catserver which allows you to forward incoming
messages from a TCP socket straight to an IRC channel and an integrated HTTP
server, ready to process incoming JSON-API calls.

At Tomahawk (http://tomahawk-player.org) we're all pretty much addicted to IRC.
We use it not only to coordinate our development efforts and communicate, but
also to get alerted about source code changes on GitHub, when a ticket gets
updated or a service on our servers runs into an issue. We can even use it to
trigger commands being executed on our servers, e.g. to deploy the latest
source code.

ircflu can do all that for us... except the communicating. It's not that kind
of chatbot.

## Installation

Make sure you have a working Go environment. See the [install instructions](http://golang.org/doc/install.html).

First we need to get the required dependencies. ircflu itself is part of that
list so the main executable can depend on our sub-packages:

    go get github.com/fluffle/goirc/client
    go get github.com/hoisie/web
    go get github.com/muesli/ircflu

Now we can build ircflu:

    git clone git://github.com/muesli/ircflu.git
    cd ircflu
    go build

To run the application you will need to specify a few parameters:

    ./ircflu -ircnick="somenickname" -irchost="some.server:6667" -ircchannel="#yourircchannel" -authpassword "yourpassword"

Run ircflu -help to see a full list of options!

## Make it talk on IRC

When started with its default options, ircflu will listen for incoming
connections on TCP port 12345. You can now send messages to an IRC channel
from your favorite shell script or a terminal:

    echo "This will be sent to the default channel on IRC." | netcat -q0 127.0.0.1 12345
    echo "#somechannel This will be sent to a specific channel on IRC." | netcat -q0 127.0.0.1 12345
    echo "@someuser This will be sent to a specific user on IRC." | netcat -q0 127.0.0.1 12345

## Integrated web hooks support

ircflu also runs an integrated HTTP server on port 12346, processing incoming
GitHub (on /github) & GitLab (on /gitlab) web-hook calls which are triggered
by a commit to your git repository. ircflu will in turn notify you on IRC.

## Remote controlling ircflu

It comes with a simplistic authentication system (!auth), supports aliases for
commands (!alias) and executing external applications (!exec), forwarding their
output to IRC.

To authenticate with ircflu, use the auth command in a private query with it:

    !auth [your_authpassword]

Once you're authed, you can execute a command on ircflu's host:

    !exec /usr/local/bin/some_executable

You can also make it join or part IRC channels:

    !join #test
    !part #test

Here's how you can create and use aliases:

    !alias deploy = exec ssh myserver ~/deploy.sh
    !deploy

More soon!

Continous integration: [![Build Status](https://secure.travis-ci.org/muesli/ircflu.png)](http://travis-ci.org/muesli/ircflu)
