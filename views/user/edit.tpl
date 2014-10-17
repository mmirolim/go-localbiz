<h1>User Edit Page</h1>
<h3>Form</h3>
<div id="validation-msg">
    <ul class="validation-erros">
        [[ range $val := .vErrs ]]
        [[ range $msg := $val ]]
        <li>[[ $msg ]]</li>
        [[ end ]]
        [[ end ]]
    </ul>
</div>
<div id="form-signup">
    <form action="/user/edit[[ if .Lang ]]?lang=[[ .Lang ]] [[ end ]]" method="post">
        First name: <input type="text" name="first_name" value="[[ .user.FirstName ]]"/><br/>
        Last name: <input type="text" name="last_name" value="[[ .user.LastName ]]"/><br/>
        Email: <input type="email" name="email" value="[[ .user.Email ]]"/><br/>
        Birthday: <input type="date" name="bday" value="[[ .user.Bday ]]"/><br/>
        <input type="radio" name="gender" value="male" [[ if eq .user.Gender "male" ]] checked [[ end ]]/>Make<br/>
        <input type="radio" name="gender" value="female" [[ if eq .user.Gender "female" ]] checked [[ end ]]/>Female<br/>
        City: should be select
        <input type="submit" value="SignUp"/>
        <input type="hidden" value="[[ .uid ]]"/>
        [[ .csrfToken ]]
    </form>
</div>
