[[ $ctrlSlug := .CtrlSlug ]]
[[ $slug := .Entity.Slug ]]
[[ $city := .Entity.Address.City ]]

<h1>[[ .Entity.Name ]]</h1>

<ul>
    [[ range $key, $val := .Entity.Types ]]
    <li class="tags label"><a href="[[getUrl $city $ctrlSlug "types" $val]]">[[ $val ]]</a></li>
    [[ end ]]
</ul>
[[ .Entity.Desc | str2html ]]

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
    <li><a href="[[getUrl $city $ctrlSlug "cuisines" $val]]">[[ $val ]]</a></li>
    [[ end ]]
</ul>
<h4 class="color-gray">[[i18n .Lang "features" ]]</h4>
<ul>
    [[ range $key, $val := .Entity.Features ]]
    <li><a href="[[getUrl $city $ctrlSlug "features" $val]]">[[ $val ]]</a></li>
    [[ end ]]
</ul>
<h4 class="color-gray">[[i18n .Lang "good-for" ]]</h4>
<ul>
    [[ range $key, $val := .Entity.GoodFor ]]
    <li><a href="[[getUrl $city $ctrlSlug "goodFor" $val]]">[[ $val ]]</a></li>
    [[ end ]]
</ul>

<h4 class="color-gray">[[i18n .Lang "around-1km" ]]</h4>
<ul>
    [[ range $key, $val := .Near.Results ]]
    <li><a href="/[[ $ctrlSlug ]]/[[ $val.Obj.Slug ]]">[[ $val.Obj.Name ]]</a> - [[ $val.Dis ]] m</li>
    [[ end ]]
</ul>
