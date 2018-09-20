package main

import "fmt"

func main() {
	/*创建集合 */
	var countryCapitalMap map[string]string
	countryCapitalMap = make(map[string]string)

	/* map插入key - value对,各个国家对应的首都 */
	countryCapitalMap[ "France" ] = "Paris"
	countryCapitalMap[ "Italy" ] = "罗马"
	countryCapitalMap[ "Japan" ] = "东京"
	countryCapitalMap[ "India " ] = "新德里"

	/*使用键输出地图值 */ 
	for country := range countryCapitalMap {
		fmt.Println(country, "首都是", countryCapitalMap[country])
	}

	/*查看元素在集合中是否存在 */
	captial, ok := countryCapitalMap[ "Japan" ] /*如果确定是真实的,则存在,否则不存在 */
	/*fmt.Println(captial) */
	/*fmt.Println(ok) */
	if (ok) {
		fmt.Println("Japan 的首都是", captial)
	} else {
		fmt.Println("Japan 的首都不存在")
	}
	/*删除元素*/
	delete(countryCapitalMap, "France")
	fmt.Println("法国条目被删除")

	captial2, ok := countryCapitalMap[ "France" ]
	if (ok) {
		fmt.Println("France 的首都是", captial2)
	} else {
		fmt.Println("France 的首都不存在")
	}
}