/*
    ErrorCode		= 0
    CreateContainerCode 	= 1
    StartContainerCode      = 2
    StopContainerCode     	= 3
    KillContainerCode	= 4
    RestartContainerCode    = 5
    RemoveContainerCode	= 6
    ListContainersCode      = 7
    PullImageCode		= 8
    RemoveImageCode		= 9
    ListImagesCode		= 10
    CommitContainerCode 	= 11
    PushImageCode		= 12 */
var NodeResponse = function(event) {
	//console.log(event.data)
	var msg = JSON.parse(event.data);
	console.log(msg)
	
	switch(msg.ReplyCode) {
		case 1:
			console.log(msg)
			ShowContainers()
			StartContainer(msg.ID)
			ConnectNodeToRing(msg.ID)
			break;
		case 2:
			ShowContainers()
			break;
		case 3:
			ShowContainers()
			break;
		case 4: 
			ShowContainers()
			break;
		case 5:
			break;
		case 6:
			ShowContainers()
			break;
		case 7: //List containers
			CreateContainerList(msg)
			break;
		case 10:
			console.log(msg.Images[0])
			break;
		case 11:
			console.log("Response: " + msg.Content)
			break;
	}
}

function CreateSkyRingList(data) {
	try {
		var json = JSON.parse(data);
		var table = document.getElementById("sky-ring-table");
		$('#sky-ring-body').empty()
		
		if(json.Nodes.length == 0) {
			$('#no-ring-text').show()
		} else {
			$('#no-ring-text').hide()
			for(var j=0; j < json.Nodes.length; j++) {
				populateSkyringList(j+1, json.Nodes[j], json, table);
			}
		}
	}catch(e) {
		console.log(e)
		$('#no-ring-text').show()
	}	
}

function populateSkyringList(nr, ringnode, json, table){
	var row = document.createElement("tr");
    var cell1 = document.createElement("td");
    var cell2 = document.createElement("td");
    var cell3 = document.createElement("td");
    
    var nrtab = document.createTextNode(nr);
    
    var idTab = document.createTextNode(ringnode.Id.substring(0, 10)+"...");
    var idCell = document.createElement("p");
    idCell.setAttribute('data-toggle', 'tooltip')
    idCell.setAttribute('data-placement', 'top')
    idCell.setAttribute('title', ringnode.Id)
    idCell.appendChild(idTab)
    
    var nameTab = document.createTextNode(ringnode.Addr.IP);
    
    cell1.appendChild(nrtab);
    cell2.appendChild(idCell);
    cell3.appendChild(nameTab);
    row.appendChild(cell1);
    row.appendChild(cell2);
    row.appendChild(cell3);
    
    table.tBodies.item("sky-ring-body").appendChild(row);
    
}

function CreateContainerList(json){
	var table = document.getElementById("container_table");
	setContainerHeaderText(json.Containers.length);
	$("#containers_body").empty();
	for (var i = 0; i < json.Containers.length; i++)
    {
		populateContainerList(i+1, json.Containers[i], json, table);
    }
}

function populateContainerList(nr, container, json, table){
	var row = document.createElement("tr");
    var cell1 = document.createElement("td");
    var cell2 = document.createElement("td");
    var cell3 = document.createElement("td");
    var cell4 = document.createElement("td");
    var cell5 = document.createElement("td");
    var cell6 = document.createElement("td");
    var cell7 = document.createElement("td");
    
    var nrtab = document.createTextNode(nr);
    
    var idTab = document.createTextNode(container.ID.substring(0, 10)+"...");
    var idCell = document.createElement("p");
    idCell.setAttribute('data-toggle', 'tooltip')
    idCell.setAttribute('data-placement', 'top')
    idCell.setAttribute('title', container.ID)
    idCell.appendChild(idTab)
    
    var imageTab = document.createTextNode(container.Image);
    var createdTab = document.createTextNode(container.Created);
    var statusTab = document.createTextNode(container.Status);
     
    //Only show start button if container is not running
    if(container.Status == "" || container.Status.indexOf("Exit") > -1 ) {
	    var startButton = document.createElement("Button");
	    var buttonText = document.createTextNode("Start");
	    startButton.className = "btn btn-success btn-margins";
	    startButton.appendChild(buttonText);
	    startButton.onclick = function()
	    {
	    	StartContainer(container.ID)
	    }
	    cell6.appendChild(startButton);
    } else {
	    var stopButton = document.createElement("Button");
	    var buttonText = document.createTextNode("Stop");
	    stopButton.className = "btn btn-warning btn-margins"
	    stopButton.appendChild(buttonText);
	    stopButton.onclick = function()
	    {
	        StopContainer(container.ID)
	    }  
	    
	    var killButton = document.createElement("Button");
	    var buttonText = document.createTextNode("Kill");
	    killButton.className = "btn btn-danger btn-margins"
	    killButton.appendChild(buttonText);
	    killButton.onclick = function()
	    {
	        KillContainer(container.ID)        
	    }
	    
	    cell6.appendChild(stopButton);
	    cell6.appendChild(killButton);
    }
    
    var deleteButton = document.createElement("Button");
    var buttonText = document.createTextNode("Delete");
    deleteButton.className = "btn btn-danger btn-margins"
    deleteButton.appendChild(buttonText);
    deleteButton.onclick = function()
    {
    	if (confirm("Do you really want to delete " + container.ID + "?") == true) {
    		DeleteContainer(container.ID)
	        table.deleteRow(nr);
    	}
    }       
    
    cell1.appendChild(nrtab);
    cell2.appendChild(idCell);
    cell3.appendChild(imageTab);
    cell4.appendChild(createdTab);
    cell5.appendChild(statusTab);
    cell6.appendChild(deleteButton);
    row.appendChild(cell1);
    row.appendChild(cell2);
    row.appendChild(cell3);
    row.appendChild(cell4);
    row.appendChild(cell5);
    row.appendChild(cell6);
    
    table.tBodies.item("containers_body").appendChild(row);
}

function setContainerHeaderText(count){
	$("#avail_cont_head").text("Containers available: "+ count);
}

function setRingHeaderText(count) {
	$("#sky-ring-title").text("Nodes in Sky-ring: " + count);
}
