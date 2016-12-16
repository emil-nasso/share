# Protocol description

There are three instances of the application involved in each transfer of a file:

* uploader - the client that is uploading the file.
* downloader - the client that is downloading the file.
* server - the broker. the instance that is getting the file from the uploader and
passing it along to the downloader

# Example communication

This is an example of a transfer of a file. All messages are of a fixed size
and are padded with the character `:` to fill that number of bytes. The size of the
message is displayed in brackets (in bytes).

uploader -> server : version v1 [64]  
server -> uploader : version v1 [64]  
uploader -> server: share [64]  
server -> uploader: hkjsfaKJ443JDSKf433 [64]  

The uploader and server makes sure that they are running the same version of the
protocol first. Then the uploader tells that server that it has something to share.
The server responds with a session id. The uploader starts listening and waiting
for the request to start uploading the file.

downloader -> server : version v1 [64]  
server -> downloader : version v1 [64]  
downloader -> server : get hkjsfaKJ443JDSKf433 [64]  

The downloader and server negotiates version and the downloader requests to downloader
the file by sending the sessionid that the uploader received from the server. The downloader
starts waiting to receive the file.

server -> uploader: start [64]  
uploader -> server: filename.txt [64]  
server -> downloader: filename.txt [64]  
uploader -> server: 2048 [64]  
server -> downloader: 2048 [64]  
uploader -> server : *file-contents* [1024]  
server -> downloader : *file-contents* [1024]  

The server tells the uploader to start sending the file. The uploader sends the filename,
filesize and then the file (in chunks of 1024 bytes). The downloader then closes the connection.
The uploader await a second `start` request, if another upload is needed to another downloader.
