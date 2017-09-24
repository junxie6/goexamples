var basicDemo = (function () {
	//Adding event listeners
	function _addEventListeners() {
		$('#showWindowButton').click(function () {
			$('#window').jqxWindow('open');
		});

		$('#userFormSubmit').on('click', function(event){
			event.preventDefault();

			var obj = {
				Username: $('#Username').val(),
				Password: '*****',
				Email: $('#Email').val(),
				Profile: {
					Name: 'Jun de profile',
					Age: 23,
					Title: 'Programmer',
				},
			};

			$('#debugTextArea').val(JSON.stringify(obj, null, 2));
		});

		$('#debugFormSubmit').on('click', function(event){
			event.preventDefault();

			var obj = JSON.parse($('#debugTextArea').val());

			console.log(obj);

			$.ajax({
				url: '/user',
				method: 'POST',
				data: JSON.stringify(obj),
			}).done(function(data, textStatus, jqXHR) {
				$('#debugTextArea').val(JSON.stringify(data, null, 2));
			}).fail(function(jqXHR, textStatus, errorThrown) {
			}).always(function(data, textStatus, errorThrown) {
			});
		});
	};

	//Creating all page elements which are jqxWidgets
	function _createElements() {
		//$('#showWindowButton').jqxButton({ width: '70px' });
		//$('#hideWindowButton').jqxButton({ width: '65px' });
	};

	//Creating the demo window
	function _createWindow() {
		$('#window').jqxWindow({
			title: 'New Candidate',
			//autoOpen: false,
			position: 'center',
			showCollapseButton: false,
			showCloseButton: true,
			maxHeight: 600,
			maxWidth: 1000,
			minHeight: 200,
			minWidth: 200,
			height: 600,
			width: 1000,
			initContent: function () {
				$('#window').jqxWindow('focus');
			}
		});
	};

	function getHTML() {
		return `
		<div id="window">
			<div id="windowHeader">
			</div>
			<div style="overflow: hidden;" id="windowContent">
				<div class="row">
					<div class="col-lg-5">
						<form class="form-horizontal">
							<!-- Username -->
							<div class="form-group">
								<label for="Username" class="col-lg-3 control-label">Username</label>
								<div class="col-lg-9">
									<input type="text" class="form-control" id="Username" placeholder="Username">
								</div>
							</div>
							<!-- Email -->
							<div class="form-group">
								<label for="Email" class="col-lg-3 control-label">Email</label>
								<div class="col-lg-9">
									<input type="text" class="form-control" id="Email" placeholder="Email">
								</div>
							</div>
							<br /><input type="submit" class="btn btn-default" id="userFormSubmit" value="Submit">
						</form>
					</div>
					<div class="col-lg-7">
						<div class="row">
							<div class="col-lg-12">
								<form>
									<textarea class="form-control" rows="12" id="debugTextArea"></textarea>
									<br /><input type="submit" class="btn btn-default" id="debugFormSubmit" value="Submit">
								</form>
							</div>
						</div>
						<div class="row">
							<div class="col-lg-12">
								List
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
			`;
	}

	return {
		config: {
			notYet: null,
		},
		init: function () {
			//Creating all jqxWindgets except the window
			_createElements();
			//Attaching event listeners
			_addEventListeners();
			//Adding jqxWindow
			_createWindow();
		},
	};
} ());

var mainWindow = (function () {
	function _addEventListeners() {
	}

	function _createElements() {
	}

	function _createWindow() {
	}

	return {
		config: {
			notYet: null,
		},
		init: function () {
			//Creating all jqxWindgets except the window
			_createElements();
			//Attaching event listeners
			_addEventListeners();
			//Adding jqxWindow
			_createWindow();
		},
	};
} ());


function bindMainButtons() {
	$('[id="manageUserBtn"]').on('click', function(event){
		event.preventDefault();

		userWindow.init();
		console.log('hi');
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


