<h1>[[ .Entity.Name ]]</h1>
[[ .Entity.Description | str2html ]]
<h4> Places in 1km around</h4>
<h4><small class="color-gray">[[ .Entity.Address.Street ]]</small></h4>
<ul>
    [[ range $key, $val := .NearResult.Results ]]
    <li>[[ $val.Dis ]]</li>
    [[ end ]]
</ul>