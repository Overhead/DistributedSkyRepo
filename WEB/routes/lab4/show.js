exports.show = function(req, res){
	res.render('lab4show', { title: 'Lab 4 Show', node_addr: req.params.node })
});