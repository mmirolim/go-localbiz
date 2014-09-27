<h1>SignUp page</h1>
<h3>Form</h3>
<div id="form-signup">
    <form action="/signup" method="post">
        First name: <input type="text" name="first_name" value="[[ .User.FirstName ]]"/><br/>
        Last name: <input type="text" name="last_name" value="[[ .User.LastName ]]"/><br/>
        Username: <input type="text" name="username" value="[[ .User.UserName ]]"/><br/>
        Birthday: <input type="date" name="bday"/><br/>
        <input type="radio" name="gender" value="male" [[ if eq .User.Gender "male" ]] checked [[ end ]]/>Make<br/>
        <input type="radio" name="gender" value="female" [[ if eq .User.Gender "female" ]] checked [[ end ]]/>Female<br/>
        City: should be select
        <input type="submit" value="SignUp"/>
    </form>
</div>