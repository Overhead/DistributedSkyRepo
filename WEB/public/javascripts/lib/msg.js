
var MsgTypeEnum = 
{
	Error : 0,
	GetKey : 1,
	PutKey : 2,
	PostKey : 3,
	DeleteKey : 4,
	ShowRing : 5
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
