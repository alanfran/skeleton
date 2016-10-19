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
             <h2 class="blog-post-title">{{.Title}}</h2>
             {{if $.Admin}}
             <div class="pull-md-right">
              <button class="btn btn-sm btn-info">Edit</button>
              <button class="btn btn-sm btn-danger">Delete</button>
             </div>
             {{end}}
             <p class="blog-post-meta">{{.DateString}} by <a href="#">{{.AuthorName}}</a></p>
             {{.Body}}
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
          <p><textarea id="postBody" rows=5 name="Body" class="form-control" placeholder="Body text goes here."></textarea></p>
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
    {{end}}
  })
  </script>
{{end}}
