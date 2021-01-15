package main

// ███████╗███╗   ██╗██╗   ██╗███████╗██╗      ██████╗ ██████╗ ███████╗███████╗
// ██╔════╝████╗  ██║██║   ██║██╔════╝██║     ██╔═══██╗██╔══██╗██╔════╝██╔════╝
// █████╗  ██╔██╗ ██║██║   ██║█████╗  ██║     ██║   ██║██████╔╝█████╗  ███████╗
// ██╔══╝  ██║╚██╗██║╚██╗ ██╔╝██╔══╝  ██║     ██║   ██║██╔═══╝ ██╔══╝  ╚════██║
// ███████╗██║ ╚████║ ╚████╔╝ ███████╗███████╗╚██████╔╝██║     ███████╗███████║
// ╚══════╝╚═╝  ╚═══╝  ╚═══╝  ╚══════╝╚══════╝ ╚═════╝ ╚═╝     ╚══════╝╚══════╝

// Envelopes are designed to modulate the amplitude of a signal.
// They should be normalized to have values between 0 and 1

import "math"

var (
	sqrt2π float64 = math.Sqrt(τ) // simple optimization
)

// Gaussian is an envelope with height 1 at μ and RMS width of σ
// f(x) = exp(-(x-μ)^2/2σ^2) μ and σ should be specified in seconds
func Gaussian(μ, σ, x Seconds) float64 {

	xu := float64(x - μ)
	return math.Exp(-xu * xu / float64(2*σ*σ))

}

// GaussianRepeat generates a sequence of gaussian envelopes of period λ seconds
func GaussianRepeat(μ, σ, λ, t Seconds) float64 {

	s := math.Mod(float64(t), float64(λ))
	tu := s - float64(μ)
	return math.Exp(-tu * tu / float64(2*σ*σ))

}
