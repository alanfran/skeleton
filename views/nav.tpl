{{define "nav"}}
<nav class="navbar navbar-static-top navbar-dark bg-inverse">
  <a class="navbar-brand" href="/">Skeleton</a>
  <ul class="nav navbar-nav">
    <li class="nav-item">
      <a id="news" class="nav-link" href="/blog">News</a>
    </li>
    <li id="portfolio" class="nav-item">
      <a class="nav-link" href="#">Portfolio</a>
    </li>
    <li id="contact" class="nav-item">
      <a class="nav-link" href="#">Contact</a>
    </li>
    <ul class="nav navbar-nav pull-xs-right">
      {{ if .Auth }}
      <li class="nav-item">
        <a id="logout" class="nav-link" href="#">Log Out</a>
      </li>
      {{ else }}
      <li class="nav-item">
        <a class="nav-link" href="#" data-toggle="modal" data-target="#login">Login</a>
      </li>
      <li class="nav-item">
        <a class="nav-link" href="#" data-toggle="modal" data-target="#register">Register</a>
      </li>
    </ul>
    {{end}}
  </ul>
</nav>
{{template "loginregister" .}}
{{end}}

{{define "loginregister"}}
<div class="modal fade" id="login" tabindex="-1" role="dialog" aria-labelledby="login" aria-hidden="true">
  <div class="modal-dialog" role="document">
    <div class="modal-content">
      <form id="loginForm" method="post" action="/api/login">
      <div class="modal-header">
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
        <h4 class="modal-title" id="login">Login</h4>
      </div>
      <div class="modal-body" id="loginBody">
          <p><input id="loginUser" type="text" name="username" class="form-control" placeholder="user name"></p>
          <p><input id="loginPass" type="password" name="password" class="form-control" placeholder="password"></p>
      </div>
      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button>
        <button type="submit" class="btn btn-primary">Login</button>
      </div>
    </form>
    </div>
  </div>
</div>
<div class="modal fade" id="register" tabindex="-1" role="dialog" aria-labelledby="register" aria-hidden="true">
  <div class="modal-dialog" role="document">
    <div class="modal-content">
      <form id="registerForm" method="post" action="/api/register">
      <div class="modal-header">
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
        <h4 class="modal-title" id="register">Register</h4>
      </div>
      <div class="modal-body" id="registerBody">
          <p><input id="registerUser" type="text" name="username" class="form-control" placeholder="user name"></p>
          <p><input id="registerPass" type="password" name="password" class="form-control" placeholder="password"></p>
          <p><input id="registerEmail" type="text" name="email" class="form-control" placeholder="email address"></p>
      </div>
      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button>
        <button type="submit" class="btn btn-primary">Register</button>
      </div>
    </form>
    </div>
  </div>
</div>
{{end}}
