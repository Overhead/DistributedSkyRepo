/*
 * GET manage page.
 */

var request = require("request");

exports.index = function(req, res){
	res.render('manage', { title: 'Lab 5' })
};
