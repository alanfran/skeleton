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
             <p class="blog-post-meta">{{.Date}} by <a href="#">{{.AuthorName}}</a></p>
             {{.Body}}
           </div><!-- /.blog-post -->
          {{end}}

           <nav class="blog-pagination">
             <a class="btn btn-outline-primary" href="#">Older</a>
             <a class="btn btn-outline-secondary disabled" href="#">Newer</a>
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

     <footer class="blog-footer">
       <p>Blog template built for <a href="http://getbootstrap.com">Bootstrap</a> by <a href="https://twitter.com/mdo">@mdo</a>.</p>
       <p>
         <a href="#">Back to top</a>
       </p>
     </footer>

  {{ template "footer" .}}
  <script>
  $(function() {
    $('#news').addClass('active');
  })
  </script>
{{end}}
