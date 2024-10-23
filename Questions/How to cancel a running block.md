
Say we need to cancel a block 

[Hylo - concurrency/cancelation](https://docs.hylo-lang.org/language-tour/concurrency#cancellation)

e.g. 
```js
  sleep: {
    seconds Int
    time.sleep(seconds, time.Second)
  }
  s : {
    sr: false
    stop_requested : { 
      sr
    }
    request_stop: {
      sr = true
    }
  }
  f: {
    s #(stop_requested #(Bool))
    while {true}, {
      s.stop_requested ? { 
        return
      }
      print("Working")
      sleep(seconds: 1)
    }
  }
  f(s)
  _: sleep(seconds: 3)
  s.request_stop()
```

A cleaner approach would be to have a flag to check if needs to keep running: 

```js
f: {
  keep_running : true
  inner: {
    while {keep_running}, {
      print("Working")
      _: time.sleep(1,time.SECOND)
    }
  }
  cancel: {
    keep_running = false
  }
  inner()
}
f()
... 
f.cancel()
```

Yz has a collaborative approach, where the block would check if needs to stop.