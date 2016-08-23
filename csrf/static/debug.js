function openAjaxLoader(){
	return function(){
		var xhr = new window.XMLHttpRequest();
		$('#ajaxLoader').jqxLoader('open');
		return xhr;
	}
}

$(document).ready(function() {
	var token = $('#JWToken').val();

	// Global AJAX Configuration
	$.ajaxSetup({
	type: 'post',
	contentType: 'application/json; charset=utf-8',
	dataType: 'json',
	timeout: 60 * 1000,
	headers: {
		Authorization: token,
	},
	//xhr: openAjaxLoader(),
	});

	//
	BindLogin();
});

function BindLogin() {
	$('#LoginBtn').on('click', function(event){
		event.preventDefault();

		var optObj = {
			Data: {
				SQLDealerName: 'Care 1',
				SQLIDShipAddr: 111,
				SQLPrice: 111.11,
			},
			ObjArr: [
				{DealerName: 'Health 1', IDShipAddr: 222, Price: 222.22},
			],
		};

		$.ajax({url: '/api/user/1', type: 'post', data: JSON.stringify(optObj),
			headers: {
				Authorization: $('#JWToken').val(),
			},
		}).done(function(dataObj){
			console.log(dataObj);
		}).fail(function(dataObj){
		}).always(function(dataObj) {
		});
	});
}

