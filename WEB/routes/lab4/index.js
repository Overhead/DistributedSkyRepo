exports.index = function(req, res){
	res.render('lab4', { title: 'Lab 4', node_addr: req.query.ip })
};