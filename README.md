# HalFS

This project uses [my friend's URL shortener](https://github.com/hn275/shorturl) as a backend for a file storage system! So far I've only got a blob store, but eventually it might become a FUSE driver as well, something similar to [how IPFS does it](https://docs.ipfs.tech/concepts/file-systems/).

## S3

The blob store backend is called S3 (not associated with AWS S3).

### Build

```console
$ go build github.com/malcolmseyd/halfs/cmd/s3
```

### Run

```console
$ ./s3
Usage: s3 [COMMAND]

Commands:
        put <FILE>          stores FILE remotely and prints a reference name
        get <NAME> <FILE>   retrieves NAME and writes it to FILE

$ ./s3 put trollface.png 
bT.5S.cT.4S.3S.ZS.YS.1S.eT.SS.gT.hT.2S.WS.RS.7S.US.VS.QS.9S.TS.dT.0S.XS.6S.aT.fT.8S
$ ./s3 get bT.5S.cT.4S.3S.ZS.YS.1S.eT.SS.gT.hT.2S.WS.RS.7S.US.VS.QS.9S.TS.dT.0S.XS.6S.aT.fT.8S demo_trollface.png
Successfully wrote 82578 bytes to demo_trollface.png
```