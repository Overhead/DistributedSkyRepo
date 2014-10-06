//------------------------------------------------------------------------------
/**
*/
function CheckOnline(addrid, elementid, buttonid)
{
	var img = document.createElement("img");
	document.body.appendChild(img);
	
	var addrElement = document.getElementById(addrid);
	var statusElement = document.getElementById(elementid);
	var buttonElement = document.getElementById(buttonid);
	statusElement.innerHTML = "Pending";	
	buttonElement.disabled = true;
	
	img.loaded = false;
	img.status = statusElement;
	img.button = buttonElement;
	img.onload = function()
	{
		status = "<font color='green'>Online</font>";
		this.status.innerHTML = status;
		this.button.disabled = false;
		this.loaded = true;
	}
	
	var fail = function()
	{
		if (!this.loaded && this != window)
		{
			status = "<font color='red'>Offline</font>";
			this.status.innerHTML = status;
			this.loaded = true;
			this.src = "";
		}
	}
	
	img.src = "http://" + addrElement.innerHTML + "/ping.bmp";
	img.onerror = img.onabort = fail;
	setTimeout
	(
		fail,
		3000
	);
}

//------------------------------------------------------------------------------
/**
*/
function RedirectToSupernode(serverip)
{
	window.location.href = "/supernode" + "?ip=" + serverip;
}