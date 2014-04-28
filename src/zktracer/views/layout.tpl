<!DOCTYPE html>

<html>
  <head>
  	<title>Zookeeper Tracer Beta</title>
  	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    
  	<link rel="stylesheet" href="//netdna.bootstrapcdn.com/bootstrap/3.0.3/css/bootstrap.min.css">
    <link rel="stylesheet" type="text/css" href="/static/css/tbsorter.css">
    <link rel="stylesheet" type="text/css" href="/static/css/layout.css">
    <link rel="stylesheet" type="text/css" href="/static/css/usual.css">
	<link href="http://cdn.bootcss.com/jstree/3.0.0/package/dist/themes/default/style.css" rel="stylesheet">
    {{.Styles}}

	</head>
  	
  <body>
  	<div class="header">
      <span class="logo"><a href="/">ZooKeeper Tracer</a></span>
      <span class="login-stat"><a href="/login">Log on</a></span>
    </div>

    <div class="navigation">
      <ul>
        <li><a href="/">Home</a></li>
        <li><a href="/broker/">Broker</a></li>
        <li><a href="/logger/">Logger</a></li>
        <li><a href="/app/">App</a></li>
        <li><a href="/topic/">Topic</a></li>
      </ul>
      <div class="clear"></div>
    </div>

  	<div class="body">
      {{.LayoutContent}}
    </div>

  	<div class="footer">
  		<span class="copy">&copy 2014 YY.COM SDS GROUP</span>
      <a class="link-item" href="http://yy.com">yy.com</a>
      <a class="link-item" href="http://duowan.com">duowan.com</a>
      <a class="link-item" href="http://100.com">100.com</a>
      <a class="link-item" href="#">Link #3</a>
      <a class="link-item" href="#">Link #5</a>
      <a class="link-item" href="#">Link #6</a>
      <a class="link-item" href="#">Link #7</a>
      <a class="link-item" href="#">xianmau.me</a>
  	</div>

    <script type="text/javascript" src="http://code.jquery.com/jquery-2.0.3.min.js"></script>
    <script src="//netdna.bootstrapcdn.com/bootstrap/3.0.3/js/bootstrap.min.js"></script>
    {{.Scripts}}
	</body>
</html>
