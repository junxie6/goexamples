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
	BindLogin();
	BindLogout();
});

function BindSalOrder() {
	$('#SalOrderBtn').on('click', function(event){
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
		event.preventDefault();

		var optObj = {
			SQLDealerName: 'Not Yet',
		};

		$.ajax({url: '/News', data: JSON.stringify(optObj),
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

function BindLogin() {
	$('#LoginBtn').on('click', function(event){
		event.preventDefault();

		var optObj = {
			Data: {
				SQLUsername: $('#Username').val(),
				SQLPassword: $('#Password').val(),
			},
		};

		$.ajax({url: '/Login', data: JSON.stringify(optObj),
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

function BindLogout() {
	$('#LogoutBtn').on('click', function(event){
		event.preventDefault();

		var optObj = {
			Data: {
			},
		};

		$.ajax({url: '/Logout', data: JSON.stringify(optObj),
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
