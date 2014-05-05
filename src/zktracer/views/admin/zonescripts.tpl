<script>
	function CreateZone(){
		alert("Create zone.")
		var newzoneid = $('#newzoneid').val();
		var remotezones = $('#remotezones').val();

		$.ajax({
			type: 'POST',
			url: '/admin/zone/create',
			data: {
				newzoneid:newzoneid, 
				remotezones:remotezones
			},
			success: function(data){
				
			},
			dataType: 'json'
		});	
	}

</script>
