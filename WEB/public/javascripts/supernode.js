//------------------------------------------------------------------------------
/**
*/
function SetupEdgeNodes(edgenodetable)
{
	var table = document.getElementById(edgenodetable);
	
	// get server ip
	var searchVars = GetSearchVars();
	var ip = searchVars['ip'];
	
	// create bitverse socket to this location
	var node = CreateWebNode("ws://" + ip + "/node");	
	
	// get children when we open the socket
	node.connectedCallback = function()
	{
		this.GetSiblings();
	}
	
	// setup table when this procedure is done
	node.childrenReceivedCallback = function(msg)
	{
		for (i = 0; i < msg.length; i++)
		{
			var childId = msg[i];
			var tr = document.createElement("tr");
			
			var cell1 = document.createElement("td");
			var cell2 = document.createElement("td");
			var cell3 = document.createElement("td");
			cell3.className = "text-center";
			
			cell2.id = childId + "_tags"
			
			var cell1Contents = document.createTextNode(childId);
			var cell3Contents = document.createElement("button");
			cell3Contents.className = "btn btn-default";
			cell3Contents.appendChild(document.createTextNode("Manage"));
			cell3Contents.ip = ip;
			cell3Contents.child = childId;
			cell3Contents.onclick = function() { window.location.href = "/manage" + "?addr=" + this.ip + "&" + "node=" + this.child; }
			
			cell1.appendChild(cell1Contents);
			cell3.appendChild(cell3Contents);
			
			tr.appendChild(cell1);
			tr.appendChild(cell2);
			tr.appendChild(cell3);
			table.appendChild(tr);
			
			// get tags
			node.GetTags(childId);
		}
	}
	
	// get node tags
	node.tagsReceivedCallback = function(node, tags)
	{
		// get node tags element
		element = document.getElementById(node + "_tags");
		
		for (var key in tags)
		{
			var keyElement = document.createTextNode(key);			
			var tagElement = document.createTextNode(tags[key]);
			var boldElement = document.createElement("b");
			element.appendChild(boldElement);
			boldElement.appendChild(keyElement);
			element.appendChild(document.createTextNode(" : "));
			element.appendChild(tagElement);
			element.appendChild(document.createElement("br"));
		}
	}
}
