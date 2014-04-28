<script src="http://cdn.bootcss.com/jstree/3.0.0/jstree.min.js"></script>
<script src="/static/js/jstree.data.json"></script>
<script>
	$(function () {
		$('#zktree_div')
		.on('init.jstree', function(){
		})
		.on('select_node.jstree', function(e, data){
			//console.log('selected');
			//console.log(data.instance.get_path(data.selected[0], '/', false));
			var current_node = data.instance.get_path(data.selected[0], '/', false);
			$('.node_path').html(current_node);
			var code = '';
			code += '<h2>';
			code += current_node;
			code += '</h2>';
			$('.show').html(code);

		})
		.jstree({
			'core':{
				'data':{
					'dataType':'json',
					'url':function(node){
						return node.id === '#' ?
							'/static/js/jstree.data.json' :
							'/admin/getdata';
					},
					'data':function(node){
						if(node.id == '#'){
							return; 
						}
						var nd =$.jstree.reference('#zktree_div');
						return { 'zoneid':$('#cur_zone').val(), 'znode':nd.get_path(node,'/',false) };
					}
				}
			},
			// 其它一些参数设置
			'plugins':["wholerow","sort"]
		});
	});
</script>
