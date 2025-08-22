package distance

import "math"

// Vincenty 使用 WGS-84 椭球参数
const (
	a = 6378137.0         // 半长轴 (赤道半径) 米
	f = 1 / 298.257223563 // 扁率
	b = (1 - f) * a       // 半短轴 (极半径) 米
)

// Distance 计算两个经纬度之间的椭球距离（米）
func Distance(lat1, lon1, lat2, lon2 float64) (float64, error) {
	const maxIter = 200
	const tol = 1e-12

	φ1 := lat1 * math.Pi / 180
	φ2 := lat2 * math.Pi / 180
	L := (lon2 - lon1) * math.Pi / 180

	U1 := math.Atan((1 - f) * math.Tan(φ1))
	U2 := math.Atan((1 - f) * math.Tan(φ2))

	sinU1, cosU1 := math.Sin(U1), math.Cos(U1)
	sinU2, cosU2 := math.Sin(U2), math.Cos(U2)

	λ := L
	var λPrev float64
	var sinσ, cosσ, σ float64
	var sinα, cos2α, cos2σm float64

	for i := 0; i < maxIter; i++ {
		sinλ := math.Sin(λ)
		cosλ := math.Cos(λ)

		sinσ = math.Sqrt(math.Pow(cosU2*sinλ, 2) +
			math.Pow(cosU1*sinU2-sinU1*cosU2*cosλ, 2))

		if sinσ == 0 {
			return 0, nil // 重合点
		}

		cosσ = sinU1*sinU2 + cosU1*cosU2*cosλ
		σ = math.Atan2(sinσ, cosσ)

		sinα = cosU1 * cosU2 * sinλ / sinσ
		cos2α = 1 - sinα*sinα

		if cos2α == 0 {
			cos2σm = 0 // 赤道线
		} else {
			cos2σm = cosσ - 2*sinU1*sinU2/cos2α
		}

		C := f / 16 * cos2α * (4 + f*(4-3*cos2α))
		λPrev = λ
		λ = L + (1-C)*f*sinα*
			(σ+C*sinσ*(cos2σm+C*cosσ*(-1+2*math.Pow(cos2σm, 2))))

		if math.Abs(λ-λPrev) < tol {
			break
		}
	}

	u2 := cos2α * (a*a - b*b) / (b * b)
	A := 1 + u2/16384*(4096+u2*(-768+u2*(320-175*u2)))
	B := u2 / 1024 * (256 + u2*(-128+u2*(74-47*u2)))

	Δσ := B * sinσ * (cos2σm + B/4*(cosσ*(-1+2*math.Pow(cos2σm, 2))-
		B/6*cos2σm*(-3+4*sinσ*sinσ)*(-3+4*math.Pow(cos2σm, 2))))

	s := b * A * (σ - Δσ)

	return s, nil
}
