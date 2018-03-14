{{define "pgrestore"}}

<!DOCTYPE html>

{{template "header" .}}

<div class="navbar navbar-default navbar-fixed-top" role="navigation">
  <div class="container">
    <div class="navbar-header">
      <a class="navbar-brand" href="/">
              <img alt="PostgreSQL" src="/images/crunchylogo.png" height="25">
      </a>
    </div>
    <div class="navbar-collapse collapse">
      <ul class="nav navbar-nav">
        <li><a href="/web/psql">PSQL</a></li>
        <li class="active"><a href="/web/pgrestore">PG Restore</a></li>
        <li><a href="/web/pgloader">PG Loader</a></li>
      </ul>
    </div>
  </div>
</div>


  <div class="container">
    <div class="page-header">
        <h1>PSQL Data Loader <small>Crunchy PostgreSQL </small></h1>
    </div>
    <div class="container-fluid">
      <form enctype="multipart/form-data" method="post">
        Run PSQL command prior to running the loader? <input type="checkbox" name="command"> <br/>
        PSQL Command: <input type="text" name="psql-command" />
        Data File: <input type="file" name="data" /> <br/>
        <input type="hidden" name="token" value="{{.}}"/>
        <input type="submit" value="upload" />
      </form>
    </div>
  </div>

<script src="/static/js/jquery.js"></script>

</body>
</html>
{{end}}
