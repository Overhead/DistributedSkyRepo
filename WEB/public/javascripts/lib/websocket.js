function OpenSocket(address, onopen, onclose, onmessage)
{
	if ("WebSocket" in window)
	{
		var socket = new WebSocket(address);		

		socket.onopen = onopen;
		socket.onclose = onclose;
		socket.onmessage = onmessage;

		return socket;
	}

	return nil;
}
