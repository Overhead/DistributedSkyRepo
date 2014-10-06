
var MsgTypeEnum = 
{
	Handshake : 0,
	Data : 1,
	Heartbeat : 2,
	Children : 3,
	ChildJoined : 4,
	ChildLeft : 5,
	UpdateTags : 6,
	SearchTags : 7,
	GetTags : 8,
	MakeImposter : 9,
	Bye : 10	
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
