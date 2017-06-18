package main

import (
	. "github.com/deadsy/sdfx/sdf"
)

func bearing_diameter() float64  { return 16 }
func bearing_radius() float64    { return bearing_diameter() / 2 }
func bearing_thickness() float64 { return 5 }
func bearing_margin() float64    { return 0.2 }

/*
func bearing_slot_width() float64 {
	return bearing_thickness() + 2*bearing_margin()
}
func bearing_slot_radius() float64 {
	return bearing_radius() + 2*bearing_margin()
}
*/
func bearing_rounding() float64 { return 0.5 }

func bearing_bolt_shaft_diameter() float64 { return 5 }
func bearing_bolt_shaft_radius() float64   { return bearing_bolt_shaft_diameter() / 2 }

/*
func bearing_bolt_shaft_length() float64   { return 15 }
func bearing_nut_size() float64            { return 7 }
*/
func bearing_hinge_length() float64 {
	return anchor_bolt_sleeve_length() - stepper_face_to_filament_path() - bearing_thickness()/2
}

/*
func bearing_nut_depth() float64 {
	return stepper_face_to_filament_path() - bearing_bolt_shaft_length()/2 + 1
}
*/

func bearing_body(margin float64) SDF3 {
	return Cylinder3D(bearing_thickness()+2*margin, bearing_radius()+2*margin, bearing_rounding())
}
func bearing() SDF3 {
	return Difference3D(bearing_body(0), Cylinder3D(bearing_thickness()+epsilon(), bearing_bolt_shaft_radius(), 0))
}
