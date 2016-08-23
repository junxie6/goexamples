$( document ).ready(function() {
  $('body').append('user id: <input id="user_id" value="10"/>');
  $('body').append('<br><br>Token: <input id="token" size="160" />');
  $('body').append('<br><br>ajax-return-data: <textarea id="ajax-return-data" cols="160" rows="5"></textarea>');

  $('body').append('<br><br><button id="getToken">get token</button>');
  $('body').append('&nbsp;&nbsp;<button id="submitBtn">submit</button>');

  $('#getToken').on('click', function(event){
    var _user_id = $('#user_id').val();

    $.ajax({
      url: "/api/user/" + _user_id,
    }).done(function(data, status, xhr) {
      var _token = xhr.getResponseHeader('X-Csrf-Token');

      $('#token').val(_token);
      $('#ajax-return-data').val(JSON.stringify(data));
    }).fail(function(data, status, xhr) {
      console.log(data);
      console.log(status);
    });

  });

  $('#submitBtn').on('click', function(event){
    var _token = $('#token').val();

    $.ajax({
      url: "/api/signup/post",
      method: 'post',
      headers: {'X-CSRF-Token': _token},
      data: {name: 'test', pass: 'test'},
      dataType: 'json',
    }).done(function(data, status, xhr) {
      console.log(data);
      //console.log(xhr.getAllResponseHeaders());
    }).fail(function(data, status, xhr) {
      console.log(data);
      console.log(status);
    });

  });
});
