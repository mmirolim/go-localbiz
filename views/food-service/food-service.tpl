<h1>[[ .Entity.Name ]]</h1>
[[ .Entity.Description | str2html ]]
<h4 class="color-gray">[[i18n .Lang "features" ]]</h4>
<ul>
    [[ range $key, $val := .Entity.Features ]]
    <li>[[ $val ]]</li>
    [[ end ]]
</ul>
<h4 class="color-gray">[[i18n .Lang "good-for" ]]</h4>
<ul>
    [[ range $key, $val := .Entity.GoodFor ]]
    <li>[[ $val ]]</li>
    [[ end ]]
</ul>
<h4 class="color-gray">[[i18n .Lang "address" ]]</h4>
<address>
    <span class="refloc">[[ .Entity.RefLoc ]]</span>
    <span class="transport">[[ .Entity.Transport ]]</span>
    <span class="street">[[ .Entity.Address.Street ]]</span>
    <span class="city">[[ .Entity.Address.District ]] [[ .Entity.Address.City ]]</span>
</address>
<h4 class="color-gray">[[i18n .Lang "around-1km" ]]</h4>
<ul>
    [[ range $key, $val := .NearResult.Results ]]
    <li>[[ $val.Obj.Name ]] - [[ $val.Dis ]] m</li>
    [[ end ]]
</ul>