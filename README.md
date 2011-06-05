## InstantMessageServer

Just a small app in my work.
It makes my Rails app having real-time communication capability.

##Dependencies:

* go-socket.io https://github.com/madari/go-socket.io
* redis.go https://github.com/hoisie/redis.go

thanks for their great works!

##How it works

1. Use (Socket.IO)[https://github.com/LearnBoost/Socket.IO] on your client side like this:

###client side code:

    <!DOCTYPE html>
    <html>
    <head>
      <title>InstantMessageServer</title>
      <script src="/javascripts/socket.io.js" type="text/javascript"></script> 
      <script> 
        function esc(msg){
          return msg.replace(/</g, '&lt;').replace(/>/g, '&gt;');
        };
              
        var socket = new io.Socket("127.0.0.1", {port: 8080});
        socket.connect();
        socket.on('connect', function(con) {
          pak = {"hash":"1"};//"1" is user id generate by rails, in the example i set user id to 1.
          socket.send(pak);
        });
        socket.on('message', function(obj){
          var msg = document.createElement('p');
          if ('message' in obj) msg.innerHTML = '<b>[' + obj.message[2] + ']'+ esc(obj.message[0]) + ':</b> ' + esc(obj.message[1]);
          document.getElementById('chat').appendChild(message);
        });
    </script> 
    </head>
    <body>
    <h1>InstantMessageServer Sample</h1> 
    <div id="chat"></div> 
    </body>
    </html>

2. Run a Redis Server on port 6379.Get from [here](http://redis.io)

3. Compile and run InstantMessageServer

    $ git clone git://github.com/Quon/InstantMessageServer.git
    $ cd InstantMessageServer/src
    $ make
    $ ./InstantMessageServer
    
4. Open the web page  in step 1 and make show the pages url contains "127.0.0.1"
    
5. Use a redis client such as redis-cli and test following command:
    $ redis-cli
    redis> publish "user:1:general" "{\"message\":[\"User 1\",\"Hello\"]}"
    (integer) 1
    redis>
    
    
Then you will see the message on the web page.
    
## License 

(The MIT License)

Copyright (c) 2010 LearnBoost &lt;quon@quonlu.com&gt;

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
'Software'), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

