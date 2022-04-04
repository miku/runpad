# runpad

Run code from an etherpad. Designed for collaborative editing, e.g. in a
classroom setting. Works with any publicly accessible pad.

How it works:

* register your pad on [runpad.xyz](https://runpad.xyz)
* go to your own runpad output page and see output of your code, updated in real-time

What code is executed?

* we call a fenced code a snippet
* in the future you will be able to control execution with additional options

```python name=quicksort

def quicksort():
    pass

```

Each instance gets a scratch space for files as well. Each pad also gets access
to a virtual filesystem containing a variety of data.

![](static/RunPad.png)

## Usage

```
$ RUNPAD_BASE_URL=http://example.com/api RUNPAD_APIKEY=123 runpad -h
Usage of runpad:
  -a string
        etherpad api key (default "123")
  -c    show pad contents and info
  -l    list pads
  -p string
        pad name to watch (default "runpad")
  -r int
        run snippet with given id (default -1)
  -s    list snippets
  -u string
        etherpad base URL (default "http://example.com/api")
```

## Security

* each pad gets a sandbox execution environment, that limits operations
* each code snippet is allowed to run for a fixed number of seconds after which it is killed
* each pad can make outgoing requests, but they are monitored and bandwidth limited

## MVP

```
$ watch ./runpad -r 0 -p hello
```

![](static/screenie.png)

## TODO

* [ ] container sandbox
* [ ] only run code on change
* [ ] implement `watch` like monitor
* [ ] webpage to follow output

