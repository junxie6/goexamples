var userWindow = (function () {
	var source = {
		localdata: [],
		datatype: 'array',
		datafields: [
			{ name: 'IDUser', type: 'number' },
			{ name: 'Username', type: 'string' },
		],
	};

	function _addEventListeners() {
		$('[id="userWindow"]').on('close', function(event) {
			$(this).remove();
		});

		$('[id="userFormSubmit"]').on('click', function(event) {
			event.preventDefault();

			var obj = {
				Username: $('[id="Username"]').val(),
			};

			$.ajax({
				url: '/SaveUser',
				method: 'POST',
				data: JSON.stringify(obj),
			}).done(function(data, textStatus, jqXHR) {
				if (data.Status == true) {
					refreshGridData(data.Data);
				}
			}).fail(function(jqXHR, textStatus, errorThrown) {
			}).always(function(data, textStatus, errorThrown) {
			});
		});
	}

	function refreshGridData(data) {
		source.localdata = data;

		// passing "cells" to the 'updatebounddata' method will refresh only the cells values when the new rows count is equal to the previous rows count.
		$('[id="userGrid"]').jqxGrid('updatebounddata', 'cells');
	}

	function _createElements() {
		var dataAdapter = new $.jqx.dataAdapter(source);

		$('[id="userGrid"]').jqxGrid({
			width: '100%',
			source: dataAdapter,
			columnsresize: false,
			sortable: false,
			columns: [
				{ text: 'IDUser', datafield: 'IDUser', width: 200 },
				{ text: 'Username', datafield: 'Username', width: 150 },
			]
		});
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
					<div class="col-lg-4">
						<form class="form-horizontal">
							<!-- Username -->
							<div class="form-group">
								<label for="Username" class="col-lg-4 control-label">Username</label>
								<div class="col-lg-8">
									<input type="text" class="form-control" id="Username" placeholder="Username">
								</div>
							</div>
							<br /><input type="submit" class="btn btn-default" id="userFormSubmit" value="Submit">
						</form>
					</div>
					<!-- right column -->
					<div class="col-lg-8">
						<div class="row">
							<div class="col-lg-12">
								<div id="userGrid">
								</div>
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

