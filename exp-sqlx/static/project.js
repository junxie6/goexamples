var projectWindow = (function () {
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


