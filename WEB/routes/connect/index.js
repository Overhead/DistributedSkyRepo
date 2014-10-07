/*
 * GET images page.
 */

var request = require("request");

exports.index = function(req, res){
	res.render('connect', { title: 'Lab 4' })
};
