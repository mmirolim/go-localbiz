<h1>[[ .Entity.Name ]]</h1>
[[ .Entity.Description | str2html ]]
<h4 class="color-gray">[[i18n .Lang "address" ]]</h4>
<address>
    <span class="refloc">[[ .Entity.RefLoc ]]</span>
    <span class="transport">[[ .Entity.Transport ]]</span>
    <span class="street">[[ .Entity.Address.Street ]]</span>
    <span class="district">[[ .Entity.Address.District ]]</span>
    <span class="city">[[ .Entity.Address.City ]]</span>
</address>
<h4 class="color-gray">[[i18n .Lang "cuisine" ]]</h4>
<ul>
    [[ range $key, $val := .Entity.Cuisines ]]
    <li>[[ $val ]]</li>
    [[ end ]]
</ul>
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

<h4 class="color-gray">[[i18n .Lang "around-1km" ]]</h4>
<ul>
    [[ $ctrlSlug := .CtrlSlug]]
    [[ range $key, $val := .Near.Results ]]
    <li><a href="/[[ $ctrlSlug ]]/[[ $val.Obj.Slug ]]">[[ $val.Obj.Name ]]</a> - [[ $val.Dis ]] m</li>
    [[ end ]]
</ul>
