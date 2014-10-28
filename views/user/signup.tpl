<h1>SignUp page</h1>
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
    [[ .userSignUpForm ]]
</div>