function NodeResponse(event) {
	var msg = JSON.parse(event.data);
	console.log(msg)
	
	switch(msg.ReplyCode) {
	case 1:
		break;
	case 2:
		break;
	case 3:
		break;
	case 4: 
		break;
	case 5:
		break;
	case 6:
		break;
	case 7:
		console.log(msg.Containers[0])
		break;
	case 10:
		console.log(msg.Images[0])
		break;
	}
}