
var MsgTypeEnum = 
{
	GetKey : 0,
	PutKey : 1,
	PostKey : 2,
	DeleteKey : 3,
	ShowRing : 4
}

function Msg()
{
	this.Type;
	this.Payload;
	this.PayloadType;
	this.Src;
	this.Dst;
	this.Origin;
	this.Id;
    this.MsgServiceName;
}

var globalMessageCounter = 0;
Msg.prototype.GenerateMessageId = function()
{
    this.Id = this.Src + globalMessageCounter++;
}
