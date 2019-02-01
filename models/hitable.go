package models

type Hitable interface {
	Hit(r *Ray, tmin, tmax float64) (bool, *HitRecord)
}

type HitRecord struct {
	T        float64
	P        *Vector
	N        *Vector
	Material Material
}

type HitableList struct {
	List []Hitable
}

func (hl *HitableList) AddHitable(h Hitable) {
	hl.List = append(hl.List, h)
}

func (hl *HitableList) Hit(r *Ray, tmin, tmax float64) (bool, *HitRecord) {
	var record *HitRecord = nil
	hitAnything := false
	closestSoFar := tmax
	for _, hitable := range hl.List {
		willHit, point := hitable.Hit(r, tmin, closestSoFar)
		if willHit {
			hitAnything = true
			closestSoFar = point.T
			record = point
		}
	}
	return hitAnything, record
}
