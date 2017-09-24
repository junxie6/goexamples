var userWindow = (function () {
	function _addEventListeners() {
		$('[id="userWindow"]').on('close', function(event) {
			$('[id="userWindow"]').remove();
		});
	}

	function _createElements() {
	}

	function _createWindow() {
		$('[id="userWindow"]').jqxWindow({
			title: 'Manage User',
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
				$('[id="userWindow"]').jqxWindow('focus');
			},
		});
	}

	function _appendHTML() {
		var html = `
		<div id="userWindow">
			<div id="userWindowHeader">
			</div>
			<div style="overflow: hidden;" id="userWindowContent">
				<div class="row">
					<!-- left column -->
					<div class="col-lg-5">
						<form class="form-horizontal">
							<!-- Username -->
							<div class="form-group">
								<label for="Username" class="col-lg-3 control-label">Username</label>
								<div class="col-lg-9">
									<input type="text" class="form-control" id="Username" placeholder="Username">
								</div>
							</div>
							<br /><input type="submit" class="btn btn-default" id="userFormSubmit" value="Submit">
						</form>
					</div>
					<!-- right column -->
					<div class="col-lg-7">
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

		return html;
	}

	return {
		config: {
			notYet: null,
		},
		init: function () {
			$('body').append(_appendHTML());

			//Creating all jqxWindgets except the window
			_createElements();
			//Attaching event listeners
			_addEventListeners();
			//Adding jqxWindow
			_createWindow();
		},
	};
} ());

