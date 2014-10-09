
/**
 * Module dependencies.
 */

var express = require('express')
  , routes = require('./routes')
  , lab4 = require('./routes/lab4')
  , manage = require('./routes/manage')
  , requests = require('./public/javascripts/lab4.js').Request

var app = module.exports = express.createServer();
var stringify = require('json-stable-stringify');
var WebSocket = require('ws')

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
app.get('/manage', manage.index);
app.get('/lab4', lab4.index)

app.param('node', /(.*)/);
app.param('key', /(.*)/);

app.get('/lab4/:node/storage', function(req, res) {
	res.render('lab4show', { title: "Get" , node_addr: req.params.node[0]})
});

app.post('/lab4/:node/storage', function(req, res){
	//res.send("POST to storage: " + req.body.newVal)
	console.log("Post: " + req.params.node[0] + " : " + req.body.newVal)
	//res.render('lab4show', { title: "Post done", node_addr: req.params.node[0], newVal: req.body.newVal, socket_action: 1 })
	var msg = {
		type: "message",
	    text: "test sending POST",
	    id:   1,
	    date: Date.now()
	}
	try {
		var ws = new WebSocket('ws://'+ req.params.node[0] + '/node');
		ws.on('open', function() {
		    ws.send(stringify(msg));
		});
		ws.on('message', function(message) {
		    console.log('received: %s', message);
		    res.send(message)
		});
		ws.on('error', function(error) {
			res.send(error)
		})
	} catch(e) {
		console.log('Error: %s', e);
	}
});

app.put('/lab4/:node/storage/:key', function(req, res){
	//res.send("UPDATE to storage key: " + req.body.newVal)
	console.log("Post: " + req.params.node[0] + " : " + req.body.newVal)
	//res.render('lab4put', { title: "Put done", node_addr: req.params.node[0], key: req.params.key[0], newVal: req.body.newVal, socket_action: 2 })
	var msg = {
		type: "message",
	    text: "test sending PUT",
	    id:   1,
	    date: Date.now()
	}
	try {
		var ws = new WebSocket('ws://'+ req.params.node[0] + '/node');
		ws.on('open', function() {
		    ws.send(stringify(msg));
		});
		ws.on('message', function(message) {
		    console.log('received: %s', message);
		    res.send(message)
		});
		ws.on('error', function(error) {
			res.send(error)
		})
	} catch(e) {
		console.log('Error: %s', e);
	}
});

app.del('/lab4/:node/storage/:key', function(req, res){
	//res.send("DELETE to storage key: " + req.params.key[0])
	console.log("Post: " + req.params.node[0] + " : " + req.body.newVal)
	//res.render('lab4del', { title: "Delete done", node_addr: req.params.node[0], key: req.params.key[0], socket_action: 3 })
	var msg = {
		type: "message",
	    text: "test sending DELETE",
	    id:   1,
	    date: Date.now()
	}
	try {
		var ws = new WebSocket('ws://'+ req.params.node[0] + '/node');
		ws.on('open', function() {
		    ws.send(stringify(msg));
		});
		ws.on('message', function(message) {
		    console.log('received: %s', message);
		    res.send(message)
		});
		ws.on('error', function(error) {
			res.send(error)
		})
	} catch(e) {
		console.log('Error: %s', e);
	}
});

app.get('/lab4/:node/storage/:key', function(req, res){
	//res.send(requests.DoNodeGet(req.params.node[0], req.params.key[0]));
	//res.send('ip: ' + req.params.node[0])
	//res.render('lab4get', { title: "Get" , node_addr: req.params.node[0], key: req.params.key[0] });//, layout: false });
	var msg = {
			type: "message",
		    text: "test sending GET: " + req.params.key[0],
		    id:   1,
		    date: Date.now()
	}
	try {
		var ws = new WebSocket('ws://'+ req.params.node[0] + '/node');
		ws.on('open', function() {
		    ws.send(stringify(msg));
		});
		ws.on('message', function(message) {
		    console.log('received: %s', message);
		    res.send(message)
		});
		ws.on('error', function(error) {
			res.send(error)
		})
	} catch(e) {
		console.log('Error: %s', e);
	}
});


app.listen(3000, function(){
  console.log("Express server listening on port %d in %s mode", app.address().port, app.settings.env);
});
