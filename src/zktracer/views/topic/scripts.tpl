
<script src="/static/js/jquery.tablesorter.min.js"></script>
<script type="text/javascript">
	$(function(){
		$('.navigation ul li a').eq(4).addClass("active");

		$('.tablesorter').tablesorter({headers:{6:{sorter:false}}});

		// 处理状态数据
		$('.status-data').each(function(index) {
			var data = $(this).text();
			if(data == "true") {
				$(this).html('<span class="text-success">running</span>')
			}else {
				$(this).html('<span class="text-danger">stop</span>')
			}
		});

		// 处理段数据
		$('.segments-data').each(function(){
			var data = $(this).text();
			var o = jQuery.parseJSON(data);
			data = "";
			for(var i = 0; i < o.length; i ++){
				var tmp = "<strong>" + o[i].last_confirm_entry + "</strong><br />";
				tmp += "<ul>";
				for(var j = 0; j < o[i].loggers.length; j ++){
					tmp += "<li>" + o[i].loggers[j] + "</li>";
				}
				tmp += "</ul>";
				data += tmp;
			}
			var s = '<span class="label label-info" data-toggle="popover" title data-html="true" data-content="' + data + '" data-original-title="Segments Details">Show Segments Details</span>';
			$(this).html(s);
		});
		$('.segments-data .label').popover('hide');
		$('.segments-data .label').click(function(){
			$('.segments-data .label').not(this).popover('hide');
		});
	});


	$(function(){
		var n= $('#cur_tab').val();
		$('#byTab a').eq(n).tab('show');
	})


</script>
