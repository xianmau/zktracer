
<script src="/static/js/jquery.tablesorter.min.js"></script>
<script type="text/javascript">
	$(function(){
		$('.navigation ul li a').eq(3).addClass("active");

		$('.tablesorter').tablesorter({headers:{0:{sorter:false},2:{sorter:false}}});

		// 处理状态数据
		$('.status-data').each(function(index) {
			var data = $(this).text()
			if(data == "true") {
				$(this).html('<span class="text-success">running</span>')
			}else {
				$(this).html('<span class="text-danger">stop</span>')
			}
		});
	});

</script>
