
/**
 * Module dependencies.
 */

var express = require('express')
  , routes = require('./routes')
  , lab4 = require('./routes/lab4')
  , manage = require('./routes/manage')

var app = module.exports = express.createServer();
var stringify = require('json-stable-stringify');
//var WebSocket = require('ws')
var WebSocketClient = require('websocket').client;

//Modify param function to take in regex
app.param(function(name, fn){
	  if (fn instanceof RegExp) {
	    return function(req, res, next, val){
	      var captures;
	      if (captures = fn.exec(String(val))) {
	        req.params[name] = captures;
	        next();
	      } else {
	        next('route');
	      }
	    }
	  }
	});

// Configuration

app.configure(function(){
  app.set('views', __dirname + '/views');
  app.set('view engine', 'ejs');
  app.use(express.bodyParser());
  app.use(express.methodOverride());
  app.use(express.compiler({ src : __dirname + '/public', enable: ['less']}));
  app.use(app.router);
  app.use(express.static(__dirname + '/public'));
});

app.configure('development', function(){
  app.use(express.errorHandler({ dumpExceptions: true, showStack: true }));
});

app.configure('production', function(){
  app.use(express.errorHandler());
});

// Compatible

// Now less files with @import 'whatever.less' will work(https://github.com/senchalabs/connect/pull/174)
var TWITTER_BOOTSTRAP_PATH = './vendor/twitter/bootstrap/less';
express.compiler.compilers.less.compile = function(str, fn){
  try {
    var less = require('less');var parser = new less.Parser({paths: [TWITTER_BOOTSTRAP_PATH]});
    parser.parse(str, function(err, root){fn(err, root.toCSS());});
  } catch (err) {fn(err);}
}

// Routes

app.get('/', routes.index);
app.get('/home', routes.index);
app.get('/lab4', lab4.index)

app.param('node', /(.*)/);
app.param('key', /(.*)/);

app.get('/lab4/:node/storage', function(req, res) {
	res.render('lab4show', { title: "Get" , node_addr: req.params.node[0]})
});

app.post('/lab4/:node/storage', function(req, res){
	var msg = {
			Action: 3,
			Key: req.body.newVal.toString(),
		    date: Date.now()
	}
	try {
		/*var ws = new WebSocket('ws://'+ req.params.node[0] + '/node');
		ws.on('open', function() {
		    ws.send(stringify(msg));
		});
		ws.on('message', function(message) {
		    console.log('received: %s', message);
		    res.send(message)
		});
		ws.on('error', function(error) {
			res.send(error)
		})*/
		var client = new WebSocketClient();

		client.on('connectFailed', function(error) {
		    console.log('Connect Error: ' + error.toString());
		    res.send(error.toString())
		});
		
		client.on('connect', function(connection) {
		    connection.on('error', function(error) {
		        console.log("Connection Error: " + error.toString());
		        res.send(error.toString())
		    });
		    connection.on('close', function() {
		        console.log('echo-protocol Connection Closed');
		    });
		    connection.on('message', function(message) {
		        if (message.type === 'utf8') {
		            console.log("Received: '" + message.utf8Data + "'");
		            res.send(message.utf8Data)
		        }
		    });
		    console.log('WebSocket client connected');
	    	connection.sendUTF(stringify(msg).toString());
		});
		console.log("Connect to: ws://" + req.params.node[0] + "/node")
		client.connect('ws://' + req.params.node[0] + '/node', "", "http://" + req.params.node[0]);
	} catch(e) {
		console.log('Error: %s', e);
	}
});

app.put('/lab4/:node/storage/:key', function(req, res){
	var msg = {
			Action: 2,
			Key: req.params.key[0].toString(),
			NewKey: req.body.newVal.toString(),
		    date: Date.now()
	}
	try {
		/*var ws = new WebSocket('ws://'+ req.params.node[0] + '/node');
		ws.on('open', function() {
		    ws.send(stringify(msg));
		});
		ws.on('message', function(message) {
		    console.log('received: %s', message);
		    res.send(message)
		});
		ws.on('error', function(error) {
			res.send(error)
		})*/
		var client = new WebSocketClient();

		client.on('connectFailed', function(error) {
		    console.log('Connect Error: ' + error.toString());
		    res.send(error.toString())
		});
		
		client.on('connect', function(connection) {
		    connection.on('error', function(error) {
		        console.log("Connection Error: " + error.toString());
		        res.send(error.toString())
		    });
		    connection.on('close', function() {
		        console.log('echo-protocol Connection Closed');
		    });
		    connection.on('message', function(message) {
		        if (message.type === 'utf8') {
		            console.log("Received: '" + message.utf8Data + "'");
		            res.send(message.utf8Data)
		        }
		    });
		    console.log('WebSocket client connected');
	    	connection.sendUTF(stringify(msg).toString());
		});
		console.log("Connect to: ws://" + req.params.node[0] + "/node")
		client.connect('ws://' + req.params.node[0] + '/node', "", "http://" + req.params.node[0]);
	} catch(e) {
		console.log('Error: %s', e);
	}
});

app.del('/lab4/:node/storage/:key', function(req, res){
	var msg = {
			Action: 4,
			Key: req.params.key[0].toString(),
		    date: Date.now()
	}
	try {
		/*var ws = new WebSocket('ws://'+ req.params.node[0] + '/node');
		ws.on('open', function() {
		    ws.send(stringify(msg));
		});
		ws.on('message', function(message) {
		    console.log('received: %s', message);
		    res.send(message)
		});
		ws.on('error', function(error) {
			res.send(error)
		})*/
		var client = new WebSocketClient();

		client.on('connectFailed', function(error) {
		    console.log('Connect Error: ' + error.toString());
		    res.send(error.toString())
		});
		
		client.on('connect', function(connection) {
		    connection.on('error', function(error) {
		        console.log("Connection Error: " + error.toString());
		        res.send(error.toString())
		    });
		    connection.on('close', function() {
		        console.log('echo-protocol Connection Closed');
		    });
		    connection.on('message', function(message) {
		        if (message.type === 'utf8') {
		            console.log("Received: '" + message.utf8Data + "'");
		            res.send(message.utf8Data)
		        }
		    });
		    console.log('WebSocket client connected');
	    	connection.sendUTF(stringify(msg).toString());
		});
		console.log("Connect to: ws://" + req.params.node[0] + "/node")
		client.connect('ws://' + req.params.node[0] + '/node', "", "http://" + req.params.node[0]);
	} catch(e) {
		console.log('Error: %s', e);
	}
});

app.get('/lab4/:node/storage/:key', function(req, res){
	var msg = {
			Action: 1,
			Key: req.params.key[0].toString(),
		    date: Date.now()
	}
	try {
		/*var ws = new WebSocket('ws://'+ req.params.node[0] + '/node');
		ws.on('open', function() {
		    ws.send(stringify(msg));
		});
		ws.on('message', function(message) {
		    console.log('received: %s', message);
		    res.send(message)
		});
		ws.on('error', function(error) {
			res.send(error)
		})*/
		var client = new WebSocketClient();

		client.on('connectFailed', function(error) {
		    console.log('Connect Error: ' + error.toString());
		    res.send(error.toString())
		});
		
		client.on('connect', function(connection) {
		    connection.on('error', function(error) {
		        console.log("Connection Error: " + error.toString());
		        res.send(error.toString())
		    });
		    connection.on('close', function() {
		        console.log('echo-protocol Connection Closed');
		    });
		    connection.on('message', function(message) {
		        if (message.type === 'utf8') {
		            console.log("Received: '" + message.utf8Data + "'");
		            res.send(message.utf8Data)
		        }
		    });
		    console.log('WebSocket client connected');
	    	connection.sendUTF(stringify(msg).toString());
		});
		console.log("Connect to: ws://" + req.params.node[0] + "/node")
		client.connect('ws://' + req.params.node[0] + '/node', "", "http://" + req.params.node[0]);		
	} catch(e) {
		console.log('Error: %s', e);
	}
});

app.get('/manage/:node/ring', function(req, res){
	var msg = {
			Action: 5,
		    date: Date.now()
	}
	try {
		var client = new WebSocketClient();

		client.on('connectFailed', function(error) {
		    console.log('Connect Error: ' + error.toString());
		    res.send(error.toString())
		});
		
		client.on('connect', function(connection) {
		    connection.on('error', function(error) {
		        console.log("Connection Error: " + error.toString());
		        res.send(error.toString())
		    });
		    connection.on('close', function() {
		        console.log('Connection Closed');
		    });
		    connection.on('message', function(message) {
		        if (message.type === 'utf8') {
		            console.log("Received: '" + message.utf8Data + "'");
		            res.send(message.utf8Data)
		        }
		    });
		    console.log('WebSocket client connected');
	    	connection.sendUTF(stringify(msg).toString());
		});
		console.log("Connect to: ws://" + req.params.node[0] + "/node")
		client.connect('ws://' + req.params.node[0] + '/node', "", "http://" + req.params.node[0]);		
	} catch(e) {
		console.log('Error: %s', e);
	}
});

app.get('/manage', manage.index);

app.get('/manage/:node', function(req, res) {
	res.render('managenode', { title: "Node management" , node_addr: req.params.node[0]})
});

app.listen(3000, function(){
  console.log("Express server listening on port %d in %s mode", app.address().port, app.settings.env);
});
