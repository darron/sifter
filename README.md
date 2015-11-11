sifter
============

Consul watches act in a bit of a peculiar way - when you reload it - they all fire.

It's a [fairly well known bug that doesn't have a fix yet](https://github.com/hashicorp/consul/issues/571).

If it's a really lightweight item - no problem - but if it's not so lightweight - then it can be a bit of a problem.

If you take a look at what the watch passes on STDIN - you will see that if it's not actually firing it just passes:

`[]\n`

So - `sifter` is a small Go binary that helps protect against watches firing accidentally:

```
{
  "watches": [
    {
      "type": "event",
      "name": "chef-client",
      "handler": "sifter run -e chef-client"
    }
  ]
}
```

When this gets loaded into Consul - instead of launching copies and copies of `chef-client` processes - it will just say:

`stdin='blank' NOT running 'chef-client'`

When you actually send the event - the logs will show this - and the event will fire:

`stdin='[{"ID":"long-semi-random-uuid-goes-here","Name":"chef-client","Payload":null,"NodeFilter":"","ServiceFilter":"service-name","TagFilter":"consul-server","Version":1,"LTime":1137}]' exec='chef-client'`

This is just an hour long hack - it doesn't do much - just wanted to see if this would work.
