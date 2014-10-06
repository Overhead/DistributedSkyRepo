/*
 * GET images page.
 */

var request = require("request");

exports.index = function(req, res){
	res.render('lab4', { title: 'Lab 4' })
};
