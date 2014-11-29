Navigation should be here
<a href="http://yalp.go">Home</a>
[[ if ne .AuthUser.UserName "" ]]
<a href="http://yalp.go/user/[[ .AuthUser.UserName ]]">[[ .AuthUser.UserName ]]</a>
<a href="http://yalp.go/logout">Logout</a>
[[ else ]]
<a href="http://yalp.go/login">Login</a>
[[ end ]]
