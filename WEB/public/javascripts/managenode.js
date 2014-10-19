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
	var msg = JSON.parse(event.data);
	console.log(msg)
	
	switch(msg.ReplyCode) {
		case 1:
			console.log(msg)
			ShowContainers()
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
		case 7: //List containers
			CreateContainerList(msg)
			break;
		case 10:
			console.log(msg.Images[0])
			break;
	}
}

function CreateContainerList(json, node){
	var table = document.getElementById("container_table");
	setContainerHeaderText(json.Containers.length);
	$("#containers_body").empty();
	for (var i = 0; i < json.Containers.length; i++)
    {
		populateContainerList(i+1, json.Containers[i], node, json, table);
    }
}

function populateContainerList(nr, container, node, json, table){
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
	        var args = new ContainerArgs();
	        args.ID = container.ID;
	    }
	    cell6.appendChild(startButton);
    } else {
	    var stopButton = document.createElement("Button");
	    var buttonText = document.createTextNode("Stop");
	    stopButton.className = "btn btn-warning btn-margins"
	    stopButton.appendChild(buttonText);
	    stopButton.onclick = function()
	    {
	        var args = new ContainerArgs();
	        args.ID = container.ID;
	    }  
	    
	    var killButton = document.createElement("Button");
	    var buttonText = document.createTextNode("Kill");
	    killButton.className = "btn btn-danger btn-margins"
	    killButton.appendChild(buttonText);
	    killButton.onclick = function()
	    {
	        var args = new ContainerArgs();
	        args.ID = container.ID;          
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
	        var args = new RemoveContainerArgs();
	        args.ID = container.ID;
	        args.RemoveVolumes = true;
	        args.Force = true;
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