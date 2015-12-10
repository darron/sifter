sifter
============

[Consul watches](https://www.consul.io/docs/agent/watches.html) act in a bit of a interesting way - when you reload the program - they all fire.

It's a [fairly well known limitation that doesn't have a fix at the moment](https://github.com/hashicorp/consul/issues/571).

If it's a really lightweight process then it's not a big deal - but if it's not so lightweight - then it can be a bit of a problem.

If you take a look at what the watch passes on STDIN during a Consul reload - you will see that if it's not actually firing it just passes:

`[]\n`

So - `sifter` is a small Go binary that helps protect against event watches firing repeatedly:

```
{
  "watches": [
    {
      "type": "event",
      "name": "chef-client",
      "handler": "sifter run -e 'chef-client'"
    }
  ]
}
```

When this gets loaded into Consul - instead of launching copies and copies of `chef-client` processes - which is not awesome I promise - it will just say:

`location='blank' elapsed='226.912µs' exec='chef-client'`

When you actually send the event - the logs will show this - and the event will fire:

```
decoded event='chef-client' ltime='1137'
key='sifter/chef-client/i-c972d41e' value='1137'
location='complete' elapsed='47.448987142s' exec='chef-client'
```

Afterwards - if it sees the same event:

`location='duplicate' elapsed='8.269276ms' exec='chef-client'`

## Real World Usage

During a cluster wide rollout of Consul changes - here are a number of events that were prevented from firing because they were duplicate or blank:

![](http://shared.froese.org/2015/7w7to-10-49.jpg)

I only want Consul watches to fire when they're supposed to.

## Development

`brew install forego` if you want to start up Redis and Consul to test.

Have added some watches in the `config/` folder as well as example Consul output - input to sifter in the `test/` folder.
