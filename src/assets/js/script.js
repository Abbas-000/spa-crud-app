$(document).ready(function() {
  $.get("/all", function(data) {
    if (data.employees === null) {
      $("#all").html("<h3>No Employees found</h3>");
    } else {
      $("#all").empty();
      $("#all").html('<ul class="list-group"></ul>');
      $(data.employees).each(function(i, e) {
        $(".list-group").append(`<li class="list-group-item">${e.name} | ${e._id}</li>`);
      })
    }
    $(".list-group-item").each(function() {
      $('<span class="pull-right"><a href="#" data-toggle="modal" data-target="#modalUpdate"><i class="fa fa-pencil"></i></a><a href="#" data-toggle="modal" data-target="#modalDelete"><i class="fa fa-trash"></i></a></span>').appendTo(this);
    })
  });

  $("#addemployee").click(function(e) {
    e.preventDefault();
    $("#closeAdd").click();
    $.ajax({
      url: '/add',
      type: 'post',
      dataType: 'json',
      data: $('#employeeaddform').serialize(),
      success: function(data) {
        console.log('success');
        //add list
      }
    });
  });

  $("#editemployee").click(function(e) {
    $("#closeEdit").click();
    $.ajax({
      url: '/update/'+$("#emid").val(),
      type: 'put',
      dataType: 'json',
      data: $('#employeeupdateform').serialize(),
      success: function(data) {
        console.log(data);
        //update list
      }
    });
  });

  $("#deleteemployee").click(function(e) {
    var id = $("#empid").val();
    $("#closeDelete").click();
    $.ajax({
      url: '/delete/'+id,
      type: 'delete',
      dataType: 'json',
      data: $('#employeedeleteform').serialize(),
      success: function(data) {
        console.log('success');
        //delete list
      }
    });
  });
});
