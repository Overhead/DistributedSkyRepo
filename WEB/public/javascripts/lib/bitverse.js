// assumes lib/websocket has been loaded before

var secret = "3e606ad97e0a738d8da4c4c74e8cd1f1f2e016c74d85f17ac2fc3b5dab4ed6c4";

//------------------------------------------------------------------------------
/**
*/
function htoa(hexx)
{
    var hex = hexx.toString();//force conversion
    var str = '';
    for (var i = 0; i < hex.length; i += 2)
        str += String.fromCharCode(parseInt(hex.substr(i, 2), 16));
    return str;
}

//------------------------------------------------------------------------------
/**
*/
function atoh(str)
{
    var arr = [];
    for (var i = 0, l = str.length; i < l; i ++)
    {
        var hex = Number(str.charCodeAt(i)).toString(16);
        arr.push(hex);
    }
    return arr.join('');
}

//------------------------------------------------------------------------------
/**
*/
function WebNode()
{
    // generate UUID    
    var uid = generateUUID();
    uid = CryptoJS.SHA1(uid);
    this.id = uid.toString(CryptoJS.enc.Hex);

    // generate new uuid and encode as sha1
    this.superNodeId = "";
    this.socket;

    this.childrenReceivedCallback = function(children)
    {
        alert(children);
    }

    this.messageReceivedCallback = function(message)
    {
        alert(JSON.stringify(message));
    }
	
	this.tagsReceivedCallback = function(node, tags)
	{
		alert(tags);
	}
	
	this.searchTagsCallback = function(nodes)
	{
		alert(nodes);
	}

    this.connectedCallback = function()
    {
        alert("Connected!");
    }
}

//------------------------------------------------------------------------------
/**
*/
WebNode.prototype.OnOpen = function()
{
    // create handshake and send
    var message = new Msg();
    message.Type = MsgTypeEnum.Handshake;
    message.Src = this.id;

    this.Send(message);
    message.Type = MsgTypeEnum.MakeImposter;
    this.Send(message);

    // call callback
    this.connectedCallback();
}

//------------------------------------------------------------------------------
/**
*/
WebNode.prototype.OnClose = function()
{
    //alert("WebNode.OnClose called!");
}

//------------------------------------------------------------------------------
/**
*/
WebNode.prototype.OnMessage = function(msg)
{
    // parse string to message
    var message = JSON.parse(msg.data);

    if (message.Type == MsgTypeEnum.Handshake)
    {
        this.superNodeId = message.Src;
    }
    else if (message.Type == MsgTypeEnum.Children)
    {
        var children = JSON.parse(message.Payload);
        var index = children.indexOf(this.id);
        if (index > -1) { children.splice(index, 1); }
        this.childrenReceivedCallback(children);
    }
    else if (message.Type == MsgTypeEnum.Data)
    {
        var decrypted = DecryptAES(message.Payload);
        obj = JSON.parse(decrypted);
        this.messageReceivedCallback(obj);
    }
	else if (message.Type == MsgTypeEnum.GetTags)
	{
		var tags = JSON.parse(message.Payload);
		this.tagsReceivedCallback(message.Dst, tags);
	}
	else if (message.Type == MsgTypeEnum.SearchTags)
	{
		var nodes = JSON.parse(message.Payload);
		this.searchTagsCallback(nodes);
	}
}


//------------------------------------------------------------------------------
/**
*/
WebNode.prototype.GetTags = function(node)
{
	// create message
	var message = new Msg();
	message.Type = MsgTypeEnum.GetTags;
	message.Dst = node;
	
	// send message
	this.Send(message);
}

//------------------------------------------------------------------------------
/**
*/
WebNode.prototype.GetSiblings = function()
{
    // create message
    var message = new Msg();
    message.Type = MsgTypeEnum.Children;

    // send message
    this.Send(message);
}

//------------------------------------------------------------------------------
/**
*/
function MakeIV(num)
{
    var buf = new Uint8Array(num);
    window.crypto.getRandomValues(buf);
    var str = String.fromCharCode.apply(null, buf);
    return str;
}

//------------------------------------------------------------------------------
/**
*/
function StringToBytes(string)
{
    var bytes = "";
    for (var i = 0; i < string.length; i++)
    {
        bytes + string.charCodeAt(i);
    }
    return bytes;
}

//------------------------------------------------------------------------------
/**
*/
function utf8_to_b64(str)
{
    return btoa(unescape(str));
}

//------------------------------------------------------------------------------
/**
*/
function b64_to_utf8(str)
{
    return escape(window.atob(str));
}

//------------------------------------------------------------------------------
/**
*/
function EncryptAES(params)
{
    var iv = MakeIV(16);
    var keyhex = CryptoJS.enc.Hex.parse(secret);

    var base64 = btoa(params);

    var encrypted = CryptoJS.AES.encrypt(
        base64,            
        keyhex,
        {
            mode: CryptoJS.mode.CFB,
            iv: CryptoJS.enc.Latin1.parse(iv),
            padding: CryptoJS.pad.NoPadding
        }
    );

    var message = htoa(encrypted.ciphertext);
    message = btoa(iv.concat(message));
    return message;
}

//------------------------------------------------------------------------------
/**
*/
function DecryptAES(data)
{    
    var keyhex = CryptoJS.enc.Hex.parse(secret);

    // convert input from base64 to ascii
    var ascii = atob(data);

    // get content
    var content = ascii.substring(16);

    // get iv
    var iv = ascii.substring(0, 16);
    var decrypted = CryptoJS.AES.decrypt(
        {
            ciphertext: CryptoJS.enc.Latin1.parse(content), 
        },
        keyhex, 
        { 
            mode: CryptoJS.mode.CFB, 
            iv: CryptoJS.enc.Latin1.parse(iv),
            padding: CryptoJS.pad.NoPadding
        } 
    );

    // decode message
    var message = decrypted.toString();
    message = htoa(message);
    message = atob(message);
    return message;
}

//------------------------------------------------------------------------------
/**
*/
WebNode.prototype.CallRPCFunction = function(name, args, node)
{
    // create invocation, the first field is the function name, the second is a JSON encoded list of arguments
    var rpcInvoke =
    {
        Rpc_function_name : name,
        Args : JSON.stringify(args)
    };

    // create message
    var message = new Msg();
    var jsonPayload = JSON.stringify(rpcInvoke);
    var encrypted = this.Encrypt(jsonPayload);
    message.Payload = encrypted;
    message.Type = MsgTypeEnum.Data;
    message.MsgServiceName = "RPCMessageService";
    message.Dst = node;

    // send message
    this.Send(message);    
}

//------------------------------------------------------------------------------
/**
*/
WebNode.prototype.Encrypt = function(data)
{
	return EncryptAES(data);
}

//------------------------------------------------------------------------------
/**
*/
WebNode.prototype.Send = function(msg)
{
	msg.Src = this.id;
    msg.Origin = msg.Src;
    msg.GenerateMessageId();
    var json = JSON.stringify(msg);
    this.socket.send(json);
}

//------------------------------------------------------------------------------
/**
*/
function CreateWebNode(address)
{
    var node = new WebNode();
    node.socket = OpenSocket(address, 
        function() { node.OnOpen(); }, 
        function() { node.OnClose(); }, 
        function(msg) { node.OnMessage(msg) }
    );
    return node;
}
