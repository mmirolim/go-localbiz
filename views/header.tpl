Navigation should be here
<a href="http://yalp.go">Home</a>
[[ if .isAuth ]]
<a href="http://yalp.go/user/[[ .Uid ]]">[[ .UserName ]]</a>
<a href="http://yalp.go/logout">Logout</a>
[[ else ]]
<a href="http://yalp.go/login">Login</a>
[[ end ]]
<a href="http://yalp.go/logout">Logout</a>
