# MoreMath
 More math functions for Go.

## Usage
Note that MoreMath handles errors in a non-standard way. Functions do not return errors, but use a pointer to an error as a parameter. 

While this may be a bit strange, it allows for chaining functions togther. 

```go
C := A.MultiplyScalar(&err, 2).Multiply(&err, B).Transpose(&err)
if err != nil {
    log.Fatal(err)
    return
}
```

Note than if the error passed in is not nil, the function (and other MoreMath functions in the chain) will abort and return a default value.

## Other Notes

While MoreMath uses generics, types other than float64 are a work in progress and should be used at your own risk.

Many of the functions in MoreMath are not performance optimized. While more performant version may be implmented at some point, don't expect this to be the fastest math module for Go.
