{{define "forumnav"}}
<nav class="navbar navbar-static-top navbar-dark bg-inverse">
  <a class="navbar-brand" href="/">Skeleton</a>
  <ul class="nav navbar-nav">
    <li class="nav-item">
      <a id="btn-boards" class="nav-link" href="/board/list">Boards</a>
    </li>
    <li class="nav-item">
      <a id="btn-recent" class="nav-link" href="#">Recent</a>
    </li>
    <li id="btn-popular" class="nav-item">
      <a class="nav-link" href="#">Popular</a>
    </li>
    <li id="btn-tags" class="nav-item">
      <a class="nav-link" href="#">Tags</a>
    </li>
    <li id="btn-users" class="nav-item">
      <a class="nav-link" href="#">Users</a>
    </li>
    <li id="btn-search" class="nav-item">
      <a class="nav-link" href="#">Search</a>
    </li>
    <ul class="nav navbar-nav pull-xs-right">
      {{ if .Auth }}
      <li class="nav-item">
        <a id="btn-settings" class="nav-link" href="#">Settings</a>
      </li>
      <li class="nav-item">
        <a id="btn-logout" class="nav-link" href="#">Log Out</a>
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
