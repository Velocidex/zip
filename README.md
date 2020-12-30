# A Concurrent Zip writer

This is a fork of the Go standard library's Zip implementation which
provides a better writer for multiple concurrent writers.

The standard library implementation requires files to be written in
series. This is not very efficient and leads to significant slow downs
when we need to write many files since each file is compressed on a
single core.

This implementation is designed to work with many writers - each
writer compresses the file independently into a zip file and then
copies the compressed file into the archive when the writer is
closed. This allows each file to be compressed concurrently.

The main different from the standard library is that created file
writer instances need to be explicitey closed to ensure they are added
to the archive (dont forget to check the error status of Close() as it
confirms if the file was added correctly):

```golang
out_fd, err := zip.Create(name)
if err != nil {
    ...
}
_, err = ioutils.Copy(out_fd, in_reader)
if err != nil {
    ...
}
err = out_fd.Close()
if err != nil {
    ...
}
```


## Pool interface

To make using this even easier, there is a CompressorPool
implementation which accepts readers.

Each reader will be compressed and copied to the zip file by a worker
in the pool. A call to pool.Close() will wait until all workers exit
and allow the zip to be safely closed.

```golang
    pool := NewCompressorPool(context.Background(), zip, 10)

    pool.Compress(&Request{
        Reader: reader,
        Name:   "My filename",
    })

    err = pool.Close()
    if err != nil {
       ...
    }
    err = zip.Close()
    if err != nil {
       ...
    }
```
