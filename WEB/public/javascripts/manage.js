/**
 * New node file
 */

function CreateImageList(json, node){
	var edgeNode = document.getElementById("edgeNode_id").innerHTML
	var table = document.getElementById('images_table')
	$("#images_body").empty();
	var count = 0;
	for (var i = 0; i < json.Images.length; i++)
    {
		var repo = json.Images[i].RepoTags[0].split(":");
		if(repo[0] != "<none>"){
			populateImageList(count+1, json.Images[i], edgeNode, repo, node, table);
			count++;
		}
    }
	setImagesHeaderText(count);
}

function populateImageList(nr, image, edgeNode, imageName, node, table){
	var row = document.createElement("tr");
    var cell1 = document.createElement("td");
    var cell2 = document.createElement("td");
    var cell3 = document.createElement("td");
    var cell4 = document.createElement("td");
    var cell5 = document.createElement("td");
    var cell6 = document.createElement("td");
    
    var nrtab = document.createTextNode(nr);
    var repoTab = document.createTextNode(image.RepoTags);
    
    var idTab = document.createTextNode(image.ID.substring(0, 10)+"...");
    var idCell = document.createElement("p");
    idCell.setAttribute('data-toggle', 'tooltip')
    idCell.setAttribute('data-placement', 'top')
    idCell.setAttribute('title', image.ID)
    idCell.appendChild(idTab)
    
    var createdTab = document.createTextNode(image.Created)
    var sizeTab = document.createTextNode(image.Size)
    

    var createButton = document.createElement("a");
    createButton.href = "#create_modal_"+nr;
    createButton.setAttribute('data-toggle', 'modal');
    var buttonText = document.createTextNode("Create container");
    createButton.className = "btn btn-success btn-margins"
    createButton.appendChild(buttonText);
    createButton.onclick = function()
    {
    	CreateContainerPopup(row, nr, edgeNode, image.ID, imageName);
    }
    
    var deleteButton = document.createElement("Button");
    var buttonText2 = document.createTextNode("Remove");
    deleteButton.className = "btn btn-danger btn-margins"
    deleteButton.appendChild(buttonText2);
    deleteButton.onclick = function()
    {
        if (confirm("Do you really want to delete " + imageName + "?") == true) {
        	var args = new RemoveImageArgs();
	        args.name = image.ID;
	        node.CallRPCFunction("EdgeNodeHandler.RemoveImage", args, edgeNode);
	        table.deleteRow(nr);
        }
    }
    
    var pushButton = document.createElement("Button");
    var buttonText3 = document.createTextNode("Push");
    pushButton.className = "btn btn-info"
	pushButton.appendChild(buttonText3);
    pushButton.onclick = function()
    {
        if (confirm("Do you really want to push " + imageName + " to private repo?") == true) {
        	var args = new ImageArgs();
	        args.Repository = imageName[0]+":"+imageName[1];
	        node.CallRPCFunction("EdgeNodeHandler.PushImage", args, edgeNode);
	        ShowProcessDialog("Pushing " + args.Repository + " to private repo")
        }
    }
    
    
    cell1.appendChild(nrtab);
    cell2.appendChild(idCell);
    cell3.appendChild(repoTab);
    cell4.appendChild(createdTab);
    cell5.appendChild(sizeTab);
    cell6.appendChild(createButton);
    cell6.appendChild(deleteButton);
    cell6.appendChild(pushButton);
    row.appendChild(cell1);
    row.appendChild(cell2);
    row.appendChild(cell3);
    row.appendChild(cell4);
    row.appendChild(cell5);
    row.appendChild(cell6);
    table.tBodies.item("images_body").appendChild(row);
}

function CreateContainerList(json, node){
	var table = document.getElementById("container_table");
	var edgeNode = document.getElementById("edgeNode_id").innerHTML
	setContainerHeaderText(json.Containers.length);
	$("#containers_body").empty();
	for (var i = 0; i < json.Containers.length; i++)
    {
		populateContainerList(i+1, json.Containers[i], node, edgeNode, json, table);
    }
}

function populateContainerList(nr, container, node, edgeNode, json, table){
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
	        node.CallRPCFunction("EdgeNodeHandler.StartContainer", args, edgeNode);
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
	        node.CallRPCFunction("EdgeNodeHandler.StopContainer", args, edgeNode);
	    }  
	    
	    var killButton = document.createElement("Button");
	    var buttonText = document.createTextNode("Kill");
	    killButton.className = "btn btn-danger btn-margins"
	    killButton.appendChild(buttonText);
	    killButton.onclick = function()
	    {
	        var args = new ContainerArgs();
	        args.ID = container.ID;
	        node.CallRPCFunction("EdgeNodeHandler.KillContainer", args, edgeNode);            
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
	        node.CallRPCFunction("EdgeNodeHandler.RemoveContainer", args, edgeNode);
	        table.deleteRow(nr);
    	}
    }
    
    var commitButton = document.createElement("a");
    var commitButtonText = document.createTextNode("Save State");
    commitButton.className = "btn btn-info btn-margins"
	commitButton.href = "#commit_modal_"+nr;
    commitButton.setAttribute('data-toggle', 'modal');
    commitButton.appendChild(commitButtonText);
    commitButton.onclick = function()
    {
    	CommitContainerPopup(row, nr, edgeNode, container.ID);
    }
       
    
    cell1.appendChild(nrtab);
    cell2.appendChild(idCell);
    cell3.appendChild(imageTab);
    cell4.appendChild(createdTab);
    cell5.appendChild(statusTab);
    cell6.appendChild(deleteButton);
    cell6.appendChild(commitButton);
    row.appendChild(cell1);
    row.appendChild(cell2);
    row.appendChild(cell3);
    row.appendChild(cell4);
    row.appendChild(cell5);
    row.appendChild(cell6);
    
    table.tBodies.item("containers_body").appendChild(row);
}

function CreateContainerPopup(row, nr, edgeNode, image_ID, imageName){
	var modalDiv = document.createElement("div");
	modalDiv.className = "modal fade";
	modalDiv.id = "create_modal_"+nr;
	
	var modalDialogDiv = document.createElement("div");
	modalDialogDiv.className = "modal-dialog";
	
	var modalContentDiv = document.createElement("div");
	modalContentDiv.className = "modal-content";
	
	var modalHeaderDiv = document.createElement("div");
	modalHeaderDiv.className = "modal-header";
	
	var closeX = document.createElement("a");
	var closeXtext = document.createTextNode('x');
	closeX.className = "close";
	closeX.setAttribute('data-dismiss', 'modal')
	closeX.appendChild(closeXtext);
	
	var header = document.createElement("h3");
	header.innerHTML = "Create container for Image:\n" + imageName[0]+":"+imageName[1];
	
	modalHeaderDiv.appendChild(closeX);
	modalHeaderDiv.appendChild(header);
	
	var modalBodyDiv = document.createElement("div"); 
	modalBodyDiv.className = "modal-body";
	
	var infoTextP = document.createElement("p");
    var infoText = document.createTextNode("Type in name for container");
    infoTextP.appendChild(infoText);
	var inputDiv = document.createElement("div");
	inputDiv.className = "input-group input-group-lg";
	var inputField = document.createElement("input");
	inputField.className = "form-control";
	inputDiv.appendChild(inputField)
	inputDiv.setAttribute('placeholder', 'Name of container');

	modalBodyDiv.appendChild(infoTextP);
	modalBodyDiv.appendChild(inputDiv);
		
	var modalFooterDiv = document.createElement("div"); 
	modalFooterDiv.className = "modal-footer";
	
	var CreateButton = document.createElement("a");
	CreateButton.className = "btn btn-primary";
    var createText = document.createTextNode("Create");
    CreateButton.appendChild(createText);
    
    CreateButton.onclick = function()
    {
    	var args = new CreateContainerArgs();
        args.ContainerName = inputField.value;
        args.ImageName = image_ID;
        node.CallRPCFunction("EdgeNodeHandler.CreateContainer", args, edgeNode);
    }
    CreateButton.setAttribute('data-dismiss', 'modal');
    
	var CloseButton = document.createElement("a");
    var closeText = document.createTextNode("Cancel");
    CloseButton.appendChild(closeText);
	CloseButton.className = "btn";
	CloseButton.setAttribute('data-dismiss', 'modal');
	
	modalFooterDiv.appendChild(CreateButton);
	modalFooterDiv.appendChild(CloseButton);
	
	modalContentDiv.appendChild(modalHeaderDiv);
	modalContentDiv.appendChild(modalBodyDiv);
	modalContentDiv.appendChild(modalFooterDiv);
	
	modalDialogDiv.appendChild(modalContentDiv)
	
	modalDiv.appendChild(modalDialogDiv)
	
	row.appendChild(modalDiv);
	
}

function CommitContainerPopup(row, nr, edgeNode, container_id){
	var modalDiv = document.createElement("div");
	modalDiv.className = "modal fade";
	modalDiv.id = "commit_modal_"+nr;
	
	var modalDialogDiv = document.createElement("div");
	modalDialogDiv.className = "modal-dialog";
	
	var modalContentDiv = document.createElement("div");
	modalContentDiv.className = "modal-content";
	
	var modalHeaderDiv = document.createElement("div");
	modalHeaderDiv.className = "modal-header";
	
	var closeX = document.createElement("a");
	var closeXtext = document.createTextNode('x');
	closeX.className = "close";
	closeX.setAttribute('data-dismiss', 'modal')
	closeX.appendChild(closeXtext);
	
	var header = document.createElement("h3");
	header.innerHTML = "Commit container: " + container_id.substring(0,10);
	
	modalHeaderDiv.appendChild(closeX);
	modalHeaderDiv.appendChild(header);
	
	var modalBodyDiv = document.createElement("div"); 
	modalBodyDiv.className = "modal-body";
	
	var infoTextP = document.createElement("p");
    var infoText = document.createTextNode("Type in name for container");
    infoTextP.appendChild(infoText);
	var inputDiv = document.createElement("div");
	inputDiv.className = "input-group input-group";
	var inputField = document.createElement("input");
	inputField.className = "form-control";
	inputDiv.appendChild(inputField)
	inputDiv.setAttribute('placeholder', 'Name of new Image');

	modalBodyDiv.appendChild(infoTextP);
	modalBodyDiv.appendChild(inputDiv);
	
	var infoTextP2 = document.createElement("p");
    var infoText2 = document.createTextNode("Type in a tag for the new image");
    infoTextP2.appendChild(infoText2);
	var inputDiv2 = document.createElement("div");
	inputDiv2.className = "input-group input-group";
	var inputField2 = document.createElement("input");
	inputField2.className = "form-control";
	inputDiv2.appendChild(inputField2)
	inputDiv2.setAttribute('placeholder', 'Image tag');

	modalBodyDiv.appendChild(infoTextP2);
	modalBodyDiv.appendChild(inputDiv2);
	
		
	var modalFooterDiv = document.createElement("div"); 
	modalFooterDiv.className = "modal-footer";
	
	var CommitAndPushButton = document.createElement("a");
	CommitAndPushButton.className = "btn btn-primary";
    var createText = document.createTextNode("Commit container");
    CommitAndPushButton.appendChild(createText);
    
    CommitAndPushButton.onclick = function()
    {
    	if (confirm("Are you sure you want to create a new image from this container?") == true) {
	    	var args = new ContainerCommitArgs();
	        args.ContainerID = container_id;
	        args.Repository = "130.240.134.118:5000/"+inputField.value.toLowerCase();
	        args.Tag = inputField2.value;
	        node.CallRPCFunction("EdgeNodeHandler.CommitContainer", args, edgeNode);
    	}
    }
    CommitAndPushButton.setAttribute('data-dismiss', 'modal');
    
	var CloseButton = document.createElement("a");
    var closeText = document.createTextNode("Cancel");
    CloseButton.appendChild(closeText);
	CloseButton.className = "btn";
	CloseButton.setAttribute('data-dismiss', 'modal');
	
	modalFooterDiv.appendChild(CommitAndPushButton);
	modalFooterDiv.appendChild(CloseButton);
	
	modalContentDiv.appendChild(modalHeaderDiv);
	modalContentDiv.appendChild(modalBodyDiv);
	modalContentDiv.appendChild(modalFooterDiv);
	
	modalDialogDiv.appendChild(modalContentDiv)
	
	modalDiv.appendChild(modalDialogDiv)
	
	row.appendChild(modalDiv);
	
}

function GetCurrentTime(message){
	var currentdate = new Date(); 
	var datetime = "Last Sync: " + currentdate.getDate() + "/"
	                + (currentdate.getMonth()+1)  + "/" 
	                + currentdate.getFullYear() + " @ "  
	                + currentdate.getHours() + ":"  
	                + currentdate.getMinutes() + ":" 
	                + currentdate.getSeconds() + ":"
					+ currentdate.getMilliseconds();
	console.log(message +" : " + datetime)
}

function setImagesHeaderText(count){
	$("#avail_img_head").text("Images available: " + count);
}
function setContainerHeaderText(count){
	$("#avail_cont_head").text("Containers available: "+ count);
}
function setImageSrc(IP) {
	document.getElementById("graphImg").src = "http://130.240.134.115/zabbix/tst/test.php?ip="+IP+"&item=system.cpu.load[percpu,avg5]&span=14400"
}
function ShowProcessDialog(text){
	$("#waitDialogH1").text(text);
	$("#pleaseWaitDialog").modal();
}

function HideProcessDialog(){
	$("#pleaseWaitDialog").modal('hide');
}
/*
    ErrorCode :  0,
    CreateContainer : 1,
    StartContainer : 2,
    StopContainer : 3,
    KillContainer : 4,
    RestartContainer : 5,
    RemoveContainer : 6,
    ListContainers : 7,
    PullImage : 8,
    RemoveImage : 9,
    ListImages : 10,
    CommitContainer: 11,
    PushImage : 12
 */

function GetImage(array, parent){
	var result;
	for(var i = 0; i<array.length; i++){
		if(array[i] == "130.240.134.118:5000/performanceimg4:testtag" ){
			result = parent;
			break
		}
	}
	return parent;
}

var NodeReceiveCallback = function(reply)
{
	var node = this;
    var edgeNode = document.getElementById('edgeNode_id').innerHTML
    
    var json = JSON.parse(reply);
    if(json.IP) {
    	setImageSrc(json.IP)
    }else if(json.ReplyCode == RPCReplyCode.ListImages){ // ListImages
    	//CreateImageList(json, node)
    	var result;
    	for(var i=0; i<json.Images.length; i++){
	    	result = GetImage(json.Images[i].RepoTags, json.Images[i])
	    	if(!$.isEmptyObject(result)){
	    		break;
	    	}
    	}
    	GetCurrentTime("Creating");
    	doPerformanceTestCreate(result.ID)
    	
    } else if(json.ReplyCode == RPCReplyCode.ListContainers) { // ListContainers
    	CreateContainerList(json, node)
    } else if(json.ReplyCode == RPCReplyCode.CreateContainer/*> 0 && json.ReplyCode < 6*/){ // Create, Start, Stop, Kill, Restart containers
    	//listContainers(edgeNode);
    	//alert(json.Content);
    	GetCurrentTime("Created!");
    	GetCurrentTime("Starting");
    	doPerformanceTestStart(json.ID)
    } else if(json.ReplyCode == RPCReplyCode.StartContainer/*> 0 && json.ReplyCode < 6*/){ // Create, Start, Stop, Kill, Restart containers
    	//listContainers(edgeNode);
    	//alert(json.Content);
    	GetCurrentTime("Started!");
    } else if(json.ReplyCode == RPCReplyCode.RemoveContainer){ // Remove container
    	setContainerHeaderText(document.getElementById("container_table").rows.length-1) //Use -1 since the header-row is counted here
    	alert(json.Content);
    } else if(json.ReplyCode == RPCReplyCode.RemoveImage){ // Remove Image
    	setImagesHeaderText(document.getElementById("images_table").rows.length-1) //Use -1 since the header-row is counted here
    	alert(json.Content);
    } else if(json.ReplyCode == RPCReplyCode.PullImage){ // Pull image
    	//listImages(edgeNode);
    	//HideProcessDialog();
    	//alert(json.Content);
    	GetCurrentTime("Pulled!");
    	listImages();
    	
    }else if(json.ReplyCode == RPCReplyCode.PushImage){ // Pull image
    	//listImages(edgeNode);
    	//HideProcessDialog();
    	//alert(json.Content);
    	GetCurrentTime("Pushed!");
    	GetCurrentTime("Pulling");
    	doPerformanceTestPull();
    	
    } else if(json.ReplyCode == RPCReplyCode.CommitContainer){ //Commit container
    	//listImages(edgeNode);
    	//alert(json.Content);
    	GetCurrentTime("Commited!");
    	GetCurrentTime("Pushing!");
    	doPerformanceTestPush();
    } else {
    	HideProcessDialog();
    	alert(json.Content);
    }
    
}
