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
	BindNews3();
	BindCSRFToken();
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

		$.ajax({url: '/SalOrder?IDSalOrder=5', data: JSON.stringify(optObj),
			headers: {
				Authorization: $('#User').val(),
				'X-CSRF-Token': $('#CSRFToken').val(),
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

		$.ajax({url: '/News1', data: JSON.stringify(optObj),
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

function BindNews3() {
	$('#News3Btn').on('click', function(event){
		event.preventDefault();

		var optObj = {
			SQLDealerName: 'Not Yet',
		};

		$.ajax({url: '/News3', data: JSON.stringify(optObj),
			headers: {
				Authorization: $('#User').val(),
				'X-CSRF-Token': $('#CSRFToken').val(),
			},
		}).done(function(dataObj){
			console.log(dataObj);
		}).fail(function(dataObj){
		}).always(function(dataObj) {
		});
	});
}

function BindCSRFToken() {
	$('#CSRFBtn').on('click', function(event){
		event.preventDefault();

		var optObj = {
			SQLDealerName: 'Not Yet',
		};

		$.ajax({url: '/CSRFToken', type: 'GET', data: optObj,
		}).done(function(dataObj, textStatus, jqXHR){
			console.log(dataObj);
			console.log(jqXHR.getResponseHeader('X-CSRF-Token'));

			$('#CSRFToken').val(jqXHR.getResponseHeader('X-CSRF-Token'));
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
