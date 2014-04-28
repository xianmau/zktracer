<script src="http://code.highcharts.com/stock/highstock.js"></script>
<script src="http://code.highcharts.com/stock/modules/exporting.js"></script>
<script src="/static/js/jquery.tablesorter.min.js"></script>
<script type="text/javascript">
	$(function(){
		$('.navigation ul li a').eq(1).addClass("active");

		$('.tablesorter').tablesorter({headers:{0:{sorter:false},2:{sorter:false},3:{sorter:false}}});

		// 处理地址数据
		$('.addr-data').each(function(index) {
			var data = $(this).text()
			var o = jQuery.parseJSON(data);
			data = "";
			for(var i = 0; i < o.length; i ++){
				data += '<span style="width:160px;display:block;float:left;">' + '<span class="label label-info" data-toggle="tooltip" title="'+o[i]+'">IP' + (i + 1) + '</span>' + o[i] + '</span>';
			}
			$(this).html(data);
		});

		// 绑定IP提示效果
		$('.addr-data .label').tooltip('hide');

		// 处理状态数据
		$('.status-data').each(function(index) {
			var data = $(this).text()
			if(data == "true") {
				$(this).html('<span class="text-success">running</span>')
			}else {
				$(this).html('<span class="text-danger">stop</span>')
			}
		});

		// 处理性能指标数据
		$('.float-data').each(function (index) {
			var data = $(this).text()
			if(data == 0) {
				$(this).html("--")
			} else {

				$(this).html(Math.round((data*100)) + "%")
			}
		})
	});

	// 图表生成
	$(function() {
		Highcharts.setOptions({
			global : {
				useUTC : true
			}
		});
		
		// Create the chart
		$('#container').highcharts('StockChart', {
			chart : {
				events : {
					load : function() {
						// set up the updating of the chart each second
						var series0 = this.series[0];
						var series1 = this.series[1];
						var series2 = this.series[2];
						setInterval(function() {
							$.get('/broker/getlatestdata?zoneid='+$('#currentzone').val()+'&brokerid='+ $('#currentbroker').val(), function(data) {
								series0.addPoint(jQuery.parseJSON(data)[0], true, true);
								series1.addPoint(jQuery.parseJSON(data)[1], true, true);
								series2.addPoint(jQuery.parseJSON(data)[2], true, true);
							});
						}, 60000);
					}
				}
			},
			
			rangeSelector: {
				buttons: [{
					count: 30,
					type: 'minute',
					text: '30M'
				}, {
					count: 1,
					type: 'hours',
					text: '1H'
				}, {
					type: 'all',
					text: 'All'
				}],
				inputEnabled: false,
				selected: 0
			},
			
			title : {
				text : 'Performance Evaluation of Broker: ' + $('#currentbroker').val()
			},

			subtitle: {
			    text: 'from zone: ' + $('#currentzone').val()
			},

			yAxis: {
			    title: {
			        text: 'percentage (%)'
			    }
			},

			exporting: {
				enabled: false
			},
			
			tooltip: {
			    shared: true
			},

			legend: {
		    	enabled: true,
		    	align: 'right',
	        	backgroundColor: '#FCFFC5',
	        	borderColor: 'black',
	        	borderWidth: 0,
		    	layout: 'vertical',
		    	verticalAlign: 'top',
		    	y: 0,
		    	shadow: true,
		    	floating: true
		    },

			series : [{
				name : 'Cpu Rate(%)',
				data : jQuery.parseJSON($('#cpuData').val())
			},{
				name : 'Net Rate(%)',
				data : jQuery.parseJSON($('#netData').val())
			},{
				name : 'Disk Rate(%)',
				data : jQuery.parseJSON($('#diskData').val())
			}]
		});

	});
</script>
