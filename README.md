# share - simple file transfers

I made share because i wanted something that i could use when i just needed to copy 
that configuration file, logfile, image, iso-file or anything else from one server 
to another or from a server to my desktop. I wanted something small and simple, 
just one binary to download with wget with no dependencies. I didn't want to send the 
files to a third party and most of the servers where behind firewalls or behind nat.

That is exactly what share does.

* Transfer a file between two systems by relaying it via a server
* Fast and simple
* Even behind a firewall/nat
* A single binary to download
* The files are not stored on the server, it simply forwards the bytes
* Many OSes supported
* Open source
* Written in go

## How?

A file transfer with share involves three pars, the server, the uploader and the downloader.

As long as both the downloader and uploader can make outbound connections to the server, the file can be shared, even if they are both behind firewalls or NAT.

When run in server mode, share waits for connections from the uploader and then forwards those uploads to the downloader.

In upload mode, share tells the server that it wants to upload a file and then waits for a downloader.

To download the file, the you can either run share in download mode or download the file directly from the server via http in any browser or by using curl or wget. When this happens, the server forwards the file from the uploader to the downloader.

## Example, please!

### The server

Run the server and wait for connections

```
$(server.example.com)> share server
Server started, awaiting connections…
```

### Uploading

Request to upload a file and get a session id and a download url

```
$(uploader.example.com)> share upload server.example.com file.zip
Session-id: 123456
Url: http://server.example.com/get/123456
```

### Downloading

Download the file using the session id, either with the share binary ... 

```
$(download.example.com)> share download server.example.com 123456
Download complete
```

... or by using wget, curl or any web browser.

```
$(download.example.com)> wget http://server.example.com/get/123456
```

## Download

Go to the [releases page](https://github.com/emil-nasso/share/releases/latest) to download the latest version of share.
