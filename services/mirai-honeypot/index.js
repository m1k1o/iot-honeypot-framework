var net = require('net');

console.log("Listening for hackers!")  
net.createServer(function(socket){

  // Initialize per-connection username and password 
  var username, password;

  socket.setEncoding('utf8');  
  console.log("connected with " + socket.remoteAddress);
  socket.write("\nI'm The Atlantic's utterly unsecure internet toaster! You should hack me.\n-------------------------------------------------------------------------\n\nlogin:");

  socket.on("error", function(err){
     console.log("error: " + err);
  });

  socket.on("data", function(data){
    data = data.toString().trim();
    
    // Check if this is the introductory telnet data stream, which I'll ignore
    if( data.indexOf("\x05") != -1 ){
      return;
    }
    
    // If the user doesn't enter a username, re-print "login"
    if ( typeof username == "undefined" ){
      if (!data) {
        socket.write("login:");
        return;
      }
      // Otherwise, store username and ask for password
      username = data;
      socket.write("password: ");
    }
    else if( username && !password ){
      // Oho! They've supplied a username and password. Let's now give them a fake shell prompt
      password = data.toString().trim();
      socket.write(username + "@poseidon:# ");
    }
    else if( username && password ) {
      // Now they've entered a command; let's log it
      var command = data.toString().trim();
      var message = "\"" + new Date() + "\",\"" + username + "\",\"" + password + "\",\"" + command + "\"," + socket.remoteAddress;
      console.log(message);
      socket.end();
    } 
  });
}).listen(23);
