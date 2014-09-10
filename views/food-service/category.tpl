<h1>[[ .Data.Category ]] <span class="color-gray">[[ .Data.Count ]]</span> </h1>
<h2>[[ .Data.City ]]</h2>

[[ range $key, $val := .Data.FdsList ]]
<li><a href="/fs/[[ $val.Slug ]]">[[ $val.Name ]]</a></li>
[[ end ]]