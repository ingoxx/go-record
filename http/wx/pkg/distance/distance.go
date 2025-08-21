package distance

import "math"

// 地球半径（单位：米）
const EarthRadius = 6371000

// 将角度转为弧度
func rad(d float64) float64 {
	return d * math.Pi / 180.0
}

// 计算两个经纬度之间的距离（米）
func Distance(lat1, lng1, lat2, lng2 float64) float64 {
	// 转换成弧度
	rlat1 := rad(lat1)
	rlat2 := rad(lat2)
	rlng1 := rad(lng1)
	rlng2 := rad(lng2)

	dlat := rlat2 - rlat1
	dlng := rlng2 - rlng1

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(rlat1)*math.Cos(rlat2)*
			math.Sin(dlng/2)*math.Sin(dlng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthRadius * c
}
