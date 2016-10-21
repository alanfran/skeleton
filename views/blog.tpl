{{define "blog"}}
  {{ template "header" .}}
  {{ template "nav" .}}

  <div class="blog-header">
       <div class="container">
         <h1 class="blog-title">The Bootstrap Blog</h1>
         <p class="lead blog-description">An example blog template built with Bootstrap.</p>
       </div>
     </div>

     <div class="container">
       <div class="row">
         <div class="col-sm-8 blog-main">
           {{range .Posts}}
           <div class="blog-post">
             <h2 id="title-{{.ID}}" class="blog-post-title">{{.Title}}</h2>
             {{if $.Admin}}
             <div class="pull-md-right">
              <button id='cancel-{{.ID}}' class='btn btn-sm btn-secondary' style="display: none">Cancel</button>
              <button id='save-{{.ID}}' class='btn btn-sm btn-primary' style="display: none">Save</button>
              <button id="edit-{{.ID}}" class="btn btn-sm btn-info">Edit</button>
              <button id="delete-{{.ID}}" class="btn btn-sm btn-danger">Delete</button>
             </div>
             {{end}}
             <p class="blog-post-meta">{{.DateString}} by <a href="#">{{.AuthorName}}</a></p>
             <p id="body-{{.ID}}">{{.Body}}</p>
           </div><!-- /.blog-post -->
          {{end}}

           <nav class="blog-pagination">
             <a class="btn btn-outline-{{if .older}}{{else}}secondary disabled{{end}}" href="#">Older</a>
             <a class="btn btn-outline-{{if .newer}}primary{{else}}secondary disabled{{end}}" href="#">Newer</a>
           </nav>

         </div><!-- /.blog-main -->

         <div class="col-sm-3 offset-sm-1 blog-sidebar">
           <div class="sidebar-module sidebar-module-inset">
             <h4>About</h4>
             <p>Etiam porta <em>sem malesuada magna</em> mollis euismod. Cras mattis consectetur purus sit amet fermentum. Aenean lacinia bibendum nulla sed consectetur.</p>
           </div>

           <div class="sidebar-module">
             <h4>Elsewhere</h4>
             <ol class="list-unstyled">
               <li><a href="#">GitHub</a></li>
               <li><a href="#">Twitter</a></li>
               <li><a href="#">Facebook</a></li>
             </ol>
           </div>
         </div><!-- /.blog-sidebar -->

       </div><!-- /.row -->

     </div><!-- /.container -->

     {{ if .Admin }}
     <footer class="blog-footer">
       <h3>New Post</h3>
       <br>
       <form id="blogForm" method="post" action="/api/blog">
       <div class="form-group">
          <p><input id="postTitle" type="text" name="Title" class="form-control" placeholder="Title"></p>
          <p><textarea id="postBody" rows=10 name="Body" class="form-control" placeholder="Body text goes here."></textarea></p>
       </div>
       <button type="submit" class="btn btn-primary">Submit</button>
     </form>
     </footer>
     {{ end }}
  {{ template "footer" .}}
  <script>
  $(function() {
    $('#news').addClass("active");

    var blogForm = $('#blogForm');

    {{ if .Admin }}
    blogForm.submit(function(event){
      event.preventDefault();
      $("#blogMsg").remove();
      $.ajax({
          type: "POST",
          url: '/api/blog',
          data: {
          Title: $("#postTitle").val(),
          Body: $("#postBody").val(),
          _csrf: "{{._csrf}}"
          },
          success: function(data) {
            location.reload(true);
          },
          error: function(r) {
            $("#blogForm").append("<div id='blogMsg' class='alert-danger'>" + r.responseText + "</div>");
          }

      });
    });

    {{range .Posts}}
      // delete .ID
      var del = $('#delete-{{.ID}}')
      del.click(function(event) {
        $.ajax({
          type: "DELETE",
          url: '/api/blog/{{.ID}}?_csrf={{$._csrf}}',
          success: function(r) {
            location.reload(true);
          },
          error: function(r) {
            alert(r.responseText)
          }
        });
      })
      // edit .ID
      var edit = $('#edit-{{.ID}}')
      var save = $("#save-{{.ID}}")
      var cancel = $("#cancel-{{.ID}}")

      save.hide()
      cancel.hide()



      edit.click(function(event) {
        // replace fields with editable form
        var title = $("#title-{{.ID}}")
        var body = $("#body-{{.ID}}")
        var editableTitle = $("<textarea id='edit-title-{{.ID}}' rows=1 class='blog-post-title' />")
        var editableBody = $("<textarea id='edit-body-{{.ID}}' rows=10 />")
        editableTitle.val(title.text())
        editableBody.val(body.text())
        $('#title-{{.ID}}').replaceWith(editableTitle)
        $('#body-{{.ID}}').replaceWith(editableBody)

        // replace Edit button with Cancel and Save
        $('#edit-{{.ID}}').hide()
        $('#save-{{.ID}}').show()
        $('#cancel-{{.ID}}').show()

        $('#cancel-{{.ID}}').click(function(event) {
          editableTitle.replaceWith(title)
          editableBody.replaceWith(body)
          $('#edit-{{.ID}}').show()
          $('#save-{{.ID}}').hide()
          $('#cancel-{{.ID}}').hide()
        })

        // Save does ajax PUT w/ csrf
        $('#save-{{.ID}}').click(function(event) {
          $.ajax({
            type: "PUT",
            url: '/api/blog/{{.ID}}',
            data: {
              ID: {{.ID}},
              Title: $('#edit-title-{{.ID}}').val(),
              Body: $('#edit-body-{{.ID}}').val(),
              _csrf: {{$._csrf}}
            },
            success: function(r) {
              location.reload(true);
            },
            error: function(r) {
              alert(r.responseText)
            }
          });
        })
      })

    {{end}}

    {{end}}
  })
  </script>
{{end}}
