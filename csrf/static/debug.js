function openAjaxLoader(){
	return function(){
		var xhr = new window.XMLHttpRequest();
		$('#ajaxLoader').jqxLoader('open');
		return xhr;
	}
}

$(document).ready(function() {
	// Global AJAX Configuration
	$.ajaxSetup({
	type: 'post',
	contentType: 'application/json; charset=utf-8',
	dataType: 'json',
	timeout: 60 * 1000,
	//headers: {
	//	Authorization: token,
	//},
	//xhr: openAjaxLoader(),
	});

	//
	BindCSRF();
	BindLogin();
	BindSalOrder();
});

function BindCSRF() {
	$('#CSRFBtn').on('click', function(event){
		event.preventDefault();

		var optObj = {IDDealer: 999};

		$.ajax({url: '/api/GetCSRFToken', type: 'get', data: optObj,
		}).done(function(dataObj, textStatus, jqXHR){
			$('#CSRFToken').val(jqXHR.getResponseHeader('X-CSRF-Token'));
		}).fail(function(dataObj){
		}).always(function(dataObj) {
		});
	});
}

function BindLogin() {
	$('#LoginBtn').on('click', function(event){
		event.preventDefault();

		$('#JSONResponse').val('');

		var Username = $('#Username').val();
		var Password = $('#Password').val();

		var optObj = {
			Data: {
				SQLDealerName: 'Care 1',
				SQLIDShipAddr: 111,
				SQLPrice: 111.11,
			},
			ObjArr: [
				{Username: Username, Password: Password},
			],
		};

		$.ajax({url: '/api/Login', type: 'post', data: JSON.stringify(optObj),
			headers: {
				'X-CSRF-Token': $('#CSRFToken').val(),
			},
		}).done(function(dataObj){
			console.log(dataObj);

			$('#JSONResponse').val(JSON.stringify(dataObj));
		}).fail(function(dataObj){
		}).always(function(dataObj) {
		});
	});
}

function BindSalOrder() {
	$('#SalOrderBtn').on('click', function(event){
		event.preventDefault();

		$('#JSONResponse').val('');

		var optObj = {
			Data: {
				SQLDealerName: 'Care 1',
				SQLIDShipAddr: 111,
				SQLPrice: 111.11,
			},
			ObjArr: [
				{IDOrder: 123},
			],
		};

		$.ajax({url: '/api/SalOrder', type: 'post', data: JSON.stringify(optObj),
			headers: {
				'X-CSRF-Token': $('#CSRFToken').val(),
			},
		}).done(function(dataObj){
			console.log(dataObj);

			$('#JSONResponse').val(JSON.stringify(dataObj));
		}).fail(function(dataObj){
		}).always(function(dataObj) {
		});
	});
}
