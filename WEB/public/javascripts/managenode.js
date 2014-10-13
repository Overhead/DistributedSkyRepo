function NodeResponse(event) {
	var msg = JSON.parse(event.data);
	console.log(msg)
	console.log(msg.text)
	document.getElementById('result').innerHTML = msg.text
}