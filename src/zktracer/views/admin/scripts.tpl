<script src="http://cdn.bootcss.com/jstree/3.0.0/jstree.min.js"></script>
<script src="/static/js/jstree.data.json"></script>
<script>
	$(function () {
		$('#zktree_div')
		.on('init.jstree', function(){
		})
		.on('select_node.jstree', function(e, data){
			var current_node = "/" + data.instance.get_path(data.selected[0], '/', false);
			$('.node_path').html(current_node);	
			$.ajax({
				url: '/admin/getnodedata',
				data: { 'zoneid':$('#cur_zone').val(), 'znode':current_node },
				success: function(data){
						var code = "";
						var  o = data;
						if(o == null){
							o = [];
						}
						code += '<h3>ZNODE: ';
						code += current_node;

						code += '</h3>';
						code += '<h3>INFO: </h3>';
						code += '<h3 style="white-space:pre;">';
						code +=	JSON.stringify(o, null, 4);
						code += '</h3>';
						code += '';
						code += '';

						$('.show').html(code);

						console.log(o);
					},
				dataType: 'json'
			});
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
	
	function CreateNode(){
		var cur_zone = $('#cur_zone').val();
		var cur_znode = $('.node_path').text();
		var new_node =  $('#new_node').val();
		var node_data = $('#node_data').val();

		$.ajax({
			type: 'POST',
			url: '/admin/createnode',
			data: {
				zoneid:cur_zone,
				znode:new_node,
				nodepath:cur_znode,
				data:node_data
			},
			success: function(data){
				//alert(data);
				
			},
			dataType: 'json'
		});
	}

	function DeleteNode(){
		var cur_zone= $('#cur_zone').val();
		var cur_znode = $('.node_path').text();
		$.ajax({
			type: 'POST',
			url: '/admin/deletenode',
			data: {
				zoneid:cur_zone,
				node:cur_znode
			},
			success: function(data){

			},
			dataType: 'json'
		});
	}
</script>
