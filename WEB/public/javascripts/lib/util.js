function GetSearchVars()
{
	//Holds key:value pairs
    var queryStringColl = null;
            
    //Get querystring from url
    var requestUrl = window.location.search.toString();

    if (requestUrl != '')
	{
        //window.location.search returns the part of the URL 
        //that follows the ? symbol, including the ? symbol
        requestUrl = requestUrl.substring(1);

        queryStringColl = new Array();

        //Get key:value pairs from querystring
        var kvPairs = requestUrl.split('&');

        for (var i = 0; i < kvPairs.length; i++)
		{
            var kvPair = kvPairs[i].split('=');
            queryStringColl[kvPair[0]] = kvPair[1];
        }
    }

    return queryStringColl;
}
