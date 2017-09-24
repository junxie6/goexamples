function bindMainButtons() {
	$('[id="manageUserBtn"]').on('click', function(event){
		event.preventDefault();

		userWindow.init();
	});
}

$( document ).ready(function() {
	// global setting
	$.ajaxSetup({
		dataType: 'json',
		contentType: 'application/json; charset=utf-8',
	});

	var html = `
		<input class="btn btn-default" value="Manage user" id="manageUserBtn" />
		<input class="btn btn-default" value="Manage project" id="manageProjectBtn" />
		<input class="btn btn-default" value="Manage role" id="manageRoleBtn" />
		<input class="btn btn-default" value="Manage permission" id="managePermissionBtn" />
		<input class="btn btn-default" value="List ticket" id="listTicketBtn" />
	`;

	$('body').append(html);

	bindMainButtons();
});


