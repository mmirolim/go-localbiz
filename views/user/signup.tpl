<h1>SignUp page</h1>
<h3>Form</h3>
<div id="validation-msg">
    <ul class="validation-erros">
        [[ range $val := .ValidationErrors ]]
            [[ range $msg := $val ]]
            <li>[[ $msg ]]</li>
            [[ end ]]
        [[ end ]]
    </ul>
</div>
<div id="form-signup">
    <form action="/signup[[ if .Lang ]]?lang=[[ .Lang ]] [[ end ]]" method="post">
        First name: <input type="text" name="first_name" value="[[ .user.FirstName ]]"/><br/>
        Last name: <input type="text" name="last_name" value="[[ .user.LastName ]]"/><br/>
        Username: <input type="text" name="username" value="[[ .user.UserName ]]"/><br/>
        Email: <input type="email" name="email" value="[[ .user.Email ]]"/><br/>
        Birthday: <input type="date" name="bday" value="[[ .user.Bday ]]"/><br/>
        <input type="radio" name="gender" value="male" [[ if eq .user.Gender "male" ]] checked [[ end ]]/>Male<br/>
        <input type="radio" name="gender" value="female" [[ if eq .user.Gender "female" ]] checked [[ end ]]/>Female<br/>
        City: should be select
        <input type="submit" value="SignUp"/>
        [[ .csrfToken ]]
    </form>
</div>