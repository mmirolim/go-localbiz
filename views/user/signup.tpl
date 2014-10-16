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
        First name: <input type="text" name="first_name" value="[[ .User.FirstName ]]"/><br/>
        Last name: <input type="text" name="last_name" value="[[ .User.LastName ]]"/><br/>
        Username: <input type="text" name="username" value="[[ .User.UserName ]]"/><br/>
        Email: <input type="email" name="email" value="[[ .User.Email ]]"/><br/>
        Birthday: <input type="date" name="bday" value="[[ .User.Bday ]]"/><br/>
        <input type="radio" name="gender" value="male" [[ if eq .User.Gender "male" ]] checked [[ end ]]/>Male<br/>
        <input type="radio" name="gender" value="female" [[ if eq .User.Gender "female" ]] checked [[ end ]]/>Female<br/>
        City: should be select
        <input type="submit" value="SignUp"/>
        [[ .csrfToken ]]
    </form>
</div>