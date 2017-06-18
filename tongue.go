package main

import (
	"math"

	. "github.com/deadsy/sdfx/sdf"
)

func tongue_depth() float64 {
	return bearing_hinge_length()*2 + bearing_thickness()
}

func tongue_bolt_shaft_diameter() float64 { return 4 }
func tongue_bolt_shaft_radius() float64   { return tongue_bolt_shaft_diameter() / 2 }
func tongue_bolt_depth() float64 {
	return stepper_face_to_filament_path() - tongue_bolt_shaft_radius() - filament_tube_radius() - 1
}
func tongue_bolt_head_diameter() float64 { return 7 }
func tongue_bolt_head_depth() float64    { return 10 }
func tongue_bolt_head_radius() float64   { return tongue_bolt_head_diameter() / 2 }
func tongue_nut_size() float64           { return 7 }
func tongue_nut_depth() float64          { return 20 }

func tongue(margin float64) SDF3 {
	tongue_rounding := anchor_bolt_sleeve_rounding()
	tongue_arc := math.Pi / 2
	min_angle := -math.Acos(
		(bearing_bolt_shaft_radius() + (bearing_radius()-bearing_bolt_shaft_radius())/2) /
			V2{feed_gear_to_bearing() - anchor_bolt_offset(), anchor_bolt_offset()}.Length())
	arm := Union3D(
		place_at_bearing(
			Transform3D(
				RevolveTheta3D(
					Transform2D(
						Box2D(V2{bearing_radius() - bearing_bolt_shaft_radius(), bearing_hinge_length()}, tongue_rounding),
						Translate2d(V2{bearing_bolt_shaft_radius() + (bearing_radius()-bearing_bolt_shaft_radius())/2, 0})),
					tongue_arc),
				RotateZ(min_angle))),
		place_at_anchor_bolt(+1, -1, Cylinder3D(bearing_hinge_length(), anchor_bolt_sleeve_radius(), anchor_bolt_sleeve_rounding())))

	armOffset := bearing_thickness()/2 + bearing_hinge_length()/2
	tongue := Union3D(
		Transform3D(arm, Translate3d(V3{0, 0, armOffset})),
		Transform3D(arm, Translate3d(V3{0, 0, -armOffset})))

	tongue = Difference3D(tongue,
		place_at_anchor_bolt(+1, -1, Cylinder3D(tongue_depth()+2*epsilon(), anchor_bolt_shaft_radius(), 0)))

	return tongue
}
