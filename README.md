# `lunch`
>tells you what is in lunch menu today directly from command line

## Install
```bash
  # install with go get and make sure you have $GOPATH/bin in your path 
  λ go get github.com/umayr/lunch/cmd/lunch
  
  # make with source
  λ git clone github.com/umayr/lunch
  λ make build

  # create cross platform binaries
  λ git clone github.com/umayr/lunch
  λ make cross
```

You can also download pre-compiled binaries from [here](https://github.com/umayr/lunch/releases/tag/v0.1).

You might also want to set up a cron to notify you about the lunch menu a bit before lunch with something like this:
```bash
  λ notify-send Lunch (echo -n "You have" (lunch) "in lunch today.")
```
which could result in:

![img](http://i.imgur.com/Fq4Y8To.png)

## License
MIT
