[[ $city := .Data.City ]]
[[ $ctrlSlug := "fs" ]]
<h2> All Types</h2>
[[ range $key, $val := .Data.CatList ]]
<li><a href="[[getUrl $city $ctrlSlug "types" $val.Id]]">[[ $val.Id ]] - [[ $val.Count ]]</a></li>
[[ end ]]
<h1>[[ .Data.Category ]] <span class="color-gray">[[ .Data.Count ]]</span> </h1>
<h2>[[ .Data.City ]]</h2>

[[ range $key, $val := .Data.FdsList ]]
<li><a href="/fs/[[ $val.Slug ]]">[[ $val.Name ]]</a></li>
[[ end ]]