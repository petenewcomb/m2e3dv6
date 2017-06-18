package main

import (
	"fmt"
	"math"

	. "github.com/deadsy/sdfx/sdf"
)

/* ======================== */

func epsilon() float64 { return 0.1 }

/* Filament */

// Nominal filament diameter
func filament_diameter() float64 { return 1.75 }
func filament_radius() float64   { return filament_diameter() / 2 }

// Bump up the hole 15% to accomodate variation
func filament_path_diameter() float64 { return filament_diameter() * 1.15 }
func filament_path_radius() float64   { return filament_path_diameter() / 2 }
func filament_path_sleeve_length() float64 {
	return filament_tube_socket_depth() + anchor_bolt_shaft_radius()
}
func filament_path_sleeve_radius() float64 {
	return filament_path_radius() + filament_path_sleeve_rounding() + 2
}
func filament_path_sleeve_rounding() float64 { return 1 }

/* Filament Tube */

func filament_tube_diameter() float64     { return 4 }
func filament_tube_radius() float64       { return filament_tube_diameter() / 2 }
func filament_tube_socket_depth() float64 { return 10 }

/* Stepper Geometry */

func stepper_protrusion_diameter() float64 { return 22 }
func stepper_protrusion_radius() float64   { return stepper_protrusion_diameter() / 2 }
func stepper_protrusion_rounding() float64 { return 0.5 }
func stepper_protrusion_cutout_radius() float64 {
	return stepper_protrusion_radius() + stepper_protrusion_rounding()
}
func stepper_protrusion_cutout_diameter() float64 { return stepper_protrusion_cutout_radius() * 2 }
func stepper_protrusion_depth() float64           { return 2.5 }

func anchor_bolt_shaft_diameter() float64 { return 3 }
func anchor_bolt_shaft_length() float64   { return 20 }
func anchor_bolt_shaft_radius() float64   { return anchor_bolt_shaft_diameter() / 2 }
func anchor_bolt_head_diameter() float64  { return 7 }
func anchor_bolt_head_radius() float64    { return anchor_bolt_head_diameter() / 2 }
func anchor_bolt_spacing() float64        { return 20 }
func anchor_bolt_offset() float64         { return anchor_bolt_spacing() / 2 }
func anchor_bolt_sleeve_radius() float64 {
	return anchor_bolt_head_radius() + anchor_bolt_sleeve_rounding()
}
func anchor_bolt_sleeve_rounding() float64 { return 1 }
func anchor_bolt_sleeve_length() float64   { return anchor_bolt_shaft_length() }

// How far from the stepper face is the center of the filament gear
func stepper_face_to_filament_path() float64 { return 14 }

/* Feed Gear */

func feed_gear_diameter() float64         { return 12 }
func feed_gear_radius() float64           { return feed_gear_diameter() / 2 }
func feed_gear_cutout_radius() float64    { return feed_gear_radius() + 1 }
func feed_gear_cutout_rounding() float64  { return 1 }
func feed_gear_bite_depth_ratio() float64 { return 0.5 }
func feed_gear_bite_depth() float64       { return feed_gear_bite_depth_ratio() * filament_diameter() }
func feed_gear_to_filament_path() float64 {
	return feed_gear_radius() + filament_radius() - feed_gear_bite_depth()
}

/* Heatsink */

// Model the heatsink as a stack of [ radius, height] cylinders
/*
func heatsink_cylinder_specs() []V2 {
	return []V2{
		V2{4, 1},
		V2{8, 4},
		V2{6, 6},
		V2{8, 4}}
}
func heatsink_cylinder_depths() { return running_total(map_subscript(heatsink_cylinder_specs(), 1)) }
func heatsink_depth()           { return heatsink_cylinder_depths()[len(heatsink_cylinder_depths())-1] }
func heatsink_wiggle_room()     { return V2{-10, 0} }
*/

/* Arbitrary Geometry Parameters (adjust to taste) */

func feed_gear_to_filament_tube() float64 { return anchor_bolt_offset() }
func feed_gear_to_heatsink() float64      { return anchor_bolt_offset() + anchor_bolt_head_radius() + 2 }
func feed_gear_cutout_depth() float64     { return anchor_bolt_shaft_length() }

func body_depth() float64 {
	return anchor_bolt_shaft_length()/2 - stepper_protrusion_depth()
}
func body_height() float64 {
	return anchor_bolt_shaft_length()/2 - anchor_bolt_sleeve_rounding()
}
func body_rounding() float64          { return 3 }
func anchor_bolt_head_depth() float64 { return 10 }

/* Useful Computed Values */

func feed_gear_to_bearing() float64 {
	return feed_gear_to_filament_path() + filament_radius() + bearing_radius()
}

//func heatsink_socket_height() float64        { return heatsink_depth() }
/*
func heatsink_socket_width() float64 {
	return heatsink_cylinder_specs()[len(heatsink_cylinder_specs())-1]
}
*/

/* Placement Functions */

func place_at_feed_gear(s SDF3) SDF3 {
	return s
}

func place_at_filament_path(s SDF3) SDF3 {
	return Transform3D(place_at_feed_gear(s), Translate3d(V3{feed_gear_to_filament_path(), 0, 0}))
}

func raise_to_filament_path(s SDF3) SDF3 {
	return Transform3D(s, Translate3d(V3{0, 0, stepper_face_to_filament_path() - anchor_bolt_shaft_length()/2}))
}

func rotate_to_filament_path(s SDF3) SDF3 {
	return Transform3D(s, RotateX(-math.Pi/2))
}

func place_on_filament_path(s SDF3) SDF3 {
	return raise_to_filament_path(place_at_filament_path(rotate_to_filament_path(s)))
}

func place_at_anchor_bolt(x, y int, s SDF3) SDF3 {
	return Transform3D(place_at_feed_gear(s), Translate3d(V3{float64(x), float64(y), 0}.MulScalar(anchor_bolt_offset())))
}

func place_at_bearing(s SDF3) SDF3 {
	return Transform3D(place_at_feed_gear(s), Translate3d(V3{feed_gear_to_bearing(), 0, 0}))
}

func place_at_heatsink_socket(s SDF3) SDF3 {
	return place_at_filament_path(Transform3D(s, Translate3d(V3{0, 0, -feed_gear_to_heatsink()})))
}

func place_at_filament_tube_socket(s SDF3) SDF3 {
	return place_at_filament_path(Transform3D(s, Translate3d(V3{0, 0, feed_gear_to_filament_tube()})))
}

func bearing_cutout() SDF3 {
	swing_radius := V2{feed_gear_to_bearing() - anchor_bolt_offset(), anchor_bolt_offset()}.Length()
	max_angle := math.Atan(anchor_bolt_offset() / (feed_gear_radius() + bearing_radius() - anchor_bolt_offset()))
	resting_angle := math.Atan(anchor_bolt_offset() / (feed_gear_to_bearing() - anchor_bolt_offset()))
	min_angle := math.Pi / 6
	arc := max_angle - min_angle
	s := Transform3D(bearing_body(bearing_margin()), Translate3d(V3{swing_radius, 0, 0}))
	s = Union3D(
		s,
		RevolveTheta3D(Slice2D(s, V3{0, 0, 0}, V3{0, 1, 0}), arc),
		Transform3D(s, RotateZ(arc)))
	s = Transform3D(s, RotateZ(-(resting_angle - min_angle)))
	s = Transform3D(s, Translate3d(V3{-swing_radius, 0, 0}))
	return Transform3D(s, RotateZ(resting_angle))
}

func body() SDF3 {
	body := Union3D(
		Transform3D(Cylinder3D(body_depth()+body_height(), V2{1, 1}.MulScalar(anchor_bolt_offset()).Length(), body_rounding()), Translate3d(V3{0, 0, (body_height() - body_depth()) / 2})),
		Transform3D(Cylinder3D(anchor_bolt_sleeve_length(), anchor_bolt_sleeve_radius(), anchor_bolt_sleeve_rounding()), Translate3d(V3{-1, +1, 0}.MulScalar(anchor_bolt_offset()))),
		Transform3D(Cylinder3D(anchor_bolt_sleeve_length(), anchor_bolt_sleeve_radius(), anchor_bolt_sleeve_rounding()), Translate3d(V3{-1, -1, 0}.MulScalar(anchor_bolt_offset()))),
		Transform3D(Cylinder3D(anchor_bolt_sleeve_length(), anchor_bolt_sleeve_radius(), anchor_bolt_sleeve_rounding()), Translate3d(V3{+1, +1, 0}.MulScalar(anchor_bolt_offset()))),
		Transform3D(Cylinder3D(anchor_bolt_sleeve_length()-tongue_depth(), anchor_bolt_sleeve_radius(), anchor_bolt_sleeve_rounding()), Translate3d(V3{+1, -1, 0}.MulScalar(anchor_bolt_offset()).Add(V3{0, 0, -tongue_depth() / 2}))),
		Transform3D(raise_to_filament_path(Cylinder3D(bearing_thickness(), anchor_bolt_sleeve_radius(), anchor_bolt_sleeve_rounding())), Translate3d(V3{+1, -1, 0}.MulScalar(anchor_bolt_offset()))),
		place_on_filament_path(Transform3D(Cylinder3D(filament_path_sleeve_length(), filament_path_sleeve_radius(), filament_path_sleeve_rounding()), Translate3d(V3{0, 0, anchor_bolt_offset() + filament_path_sleeve_length()/2}))))
	body.(*UnionSDF3).SetMin(ExpMin(0.8))

	body = Intersect3D(body, Box3D(V3{4 * anchor_bolt_offset(), 4 * anchor_bolt_offset(), anchor_bolt_shaft_length()}, 0))
	body.(*IntersectionSDF3).SetMax(PolyMax(anchor_bolt_sleeve_rounding()))

	bb := body.BoundingBox()
	filament_path_length := bb.Max.Y - bb.Min.Y

	filament_path := Union3D(
		Transform3D(Cone3D(anchor_bolt_offset()/2, 0, bearing_thickness()/2, 0), Translate3d(V3{0, 0, -anchor_bolt_offset() / 2})),
		Transform3D(Cylinder3D(filament_path_length+2*epsilon(), filament_radius(), 0), Translate3d(V3{0, 0, bb.Min.Y + filament_path_length/2})))
	filament_path.(*UnionSDF3).SetMin(ExpMin(1.0))

	body = Difference3D(body, Union3D(
		place_at_feed_gear(Transform3D(Cylinder3D(stepper_protrusion_depth()*2, stepper_protrusion_cutout_radius(), stepper_protrusion_rounding()), Translate3d(V3{0, 0, -anchor_bolt_shaft_length() / 2})))))
	body.(*DifferenceSDF3).SetMax(PolyMax(anchor_bolt_sleeve_rounding()))

	body = Difference3D(body, Union3D(
		place_at_anchor_bolt(-1, +1, Cylinder3D(anchor_bolt_sleeve_length()+epsilon(), anchor_bolt_shaft_radius(), 0)),
		place_at_anchor_bolt(-1, -1, Cylinder3D(anchor_bolt_sleeve_length()+epsilon(), anchor_bolt_shaft_radius(), 0)),
		place_at_anchor_bolt(+1, -1, Cylinder3D(anchor_bolt_sleeve_length()+epsilon(), anchor_bolt_shaft_radius(), 0)),
		place_at_anchor_bolt(+1, +1, Cylinder3D(anchor_bolt_sleeve_length()+epsilon(), anchor_bolt_shaft_radius(), 0)),
		raise_to_filament_path(place_at_bearing(bearing_cutout())),
		place_on_filament_path(Transform3D(Cylinder3D(filament_tube_socket_depth()+epsilon(), filament_tube_radius(), 0), Translate3d(V3{0, 0, bb.Max.Y - filament_tube_socket_depth()/2}))),
		place_on_filament_path(filament_path)))
	body.(*DifferenceSDF3).SetMax(PolyMax(anchor_bolt_sleeve_rounding()))

	feed_gear_cutout := Cylinder3D(anchor_bolt_sleeve_length()+2*epsilon(), feed_gear_cutout_radius(), 0)
	body = Difference3D(body, feed_gear_cutout)
	body.(*DifferenceSDF3).SetMax(PolyMax(feed_gear_cutout_rounding()))

	return body
}

func render(filename string, res float64, s SDF3) {
	bb := s.BoundingBox()
	bbs := bb.Max.Sub(bb.Min)
	cells := int(math.Ceil(bbs.MaxComponent() / res))
	fmt.Printf("size: %v\n", bb.Max.Sub(bb.Min))
	fmt.Printf("res: %v\n", res)
	fmt.Printf("mesh cells: %v\n", cells)
	RenderSTL(s, cells, filename)
}

func main() {
	res := 0.25
	render("m2e3dv6_body.stl", res, body())
	render("m2e3dv6_pressuretongue.stl", res, raise_to_filament_path(tongue(0)))
	render("m2e3dv6_hardware.stl", res, raise_to_filament_path(place_at_bearing(bearing())))
}
