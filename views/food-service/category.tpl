[[ $city := .Data.City ]]
[[ $ctrlSlug := "fs" ]]
<h2> All Types</h2>
[[ range $key, $val := .Data.CatList ]]
<li><a href="[[getUrl $city $ctrlSlug "types" $val.Id]]">[[ $val.Id ]] - [[ $val.Count ]]</a></li>
[[ end ]]
<h2>[[ .Data.City ]]</h2>

