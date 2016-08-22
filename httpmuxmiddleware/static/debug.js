function openAjaxLoader(){
	return function(){
		var xhr = new window.XMLHttpRequest();
		$('#ajaxLoader').jqxLoader('open');
		return xhr;
	}
}

$(document).ready(function() {
	var token = $('#User').val();

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
	BindSalOrder();
	BindNews();
});

function BindSalOrder() {
	$('#SalOrderBtn').on('click', function(event){
		var optObj = {
			SQLDealerName: 'Not Yet',
		};

		$.ajax({url: '/SalOrder/5', data: JSON.stringify(optObj),
			headers: {
				Authorization: $('#User').val(),
			},
		}).done(function(dataObj){
			console.log(dataObj);
		}).fail(function(dataObj){
		}).always(function(dataObj) {
		});
	});
}

function BindNews() {
	$('#NewsBtn').on('click', function(event){
		var optObj = {
			SQLDealerName: 'Not Yet',
		};

		$.ajax({url: 'http://news.local:8080/News', data: JSON.stringify(optObj),
			headers: {
				Authorization: $('#User').val(),
			},
		}).done(function(dataObj){
			console.log(dataObj);
		}).fail(function(dataObj){
		}).always(function(dataObj) {
		});
	});
}
