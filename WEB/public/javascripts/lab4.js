var RequestClass = function() {
    // run code here, or...
}; 

RequestClass.DoNodeGet = function(addr, key) {
	
	wSocket = new WebSocket("ws://" + addr + "/node");
	
   	wSocket.onmessage = function (event) {
	  console.log(event.data);
	  return event.data
	}
	
	wSocket.onopen = function (event) {
		var msg = {
			    type: "message",
			    text: "test sending GET key: " + key,
			    id:   1,
			    date: Date.now()
		    };
		wSocket.send(JSON.stringify(msg));
	}	
};

//...add a method, which we do in this example:
RequestClass.prototype.getList = function(test) {
    return "My List " + test;
};

// now expose with module.exports:
exports.Request = RequestClass;
