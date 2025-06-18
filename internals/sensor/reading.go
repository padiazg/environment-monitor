package sensor

import "fmt"

type Reading struct {
	// all big-endian, IEEE754 float value
	MassPM1             float32 // Mass Concentration PM1.0 [μg/m3]
	MassPM25            float32 // Mass Concentration PM2.5 [μg/m3]
	MassPM4             float32 // Mass Concentration PM4.0 [μg/m3]
	MassPM10            float32 // Mass Concentration PM10 [μg/m3]
	NumberPM05          float32 // Number Concentration PM0.5 [#/cm3]
	NumberPM1           float32 // Number Concentration PM1.0 [#/cm3]
	NumberPM25          float32 // Number Concentration PM2.5 [#/cm3]
	NumberPM4           float32 // Number Concentration PM4.0 [#/cm3]
	NumberPM10          float32 // Number Concentration PM10 [#/cm3]
	TypicalParticleSize float32 // Typical Particle Size [μm]
}

func (r *Reading) Show() {
	fmt.Printf("Mass Concentration PM1.0  : %3.2f [μg/m3]\n", r.MassPM1)
	fmt.Printf("Mass Concentration PM2.5  : %3.2f [μg/m3]\n", r.MassPM25)
	fmt.Printf("Mass Concentration PM4.0  : %3.2f [μg/m3]\n", r.MassPM4)
	fmt.Printf("Mass Concentration PM10   : %3.2f [μg/m3]\n", r.MassPM10)
	fmt.Printf("Number Concentration PM0.5: %3.2f [#/cm3]\n", r.NumberPM05)
	fmt.Printf("Number Concentration PM1.0: %3.2f [#/cm3]\n", r.NumberPM1)
	fmt.Printf("Number Concentration PM2.5: %3.2f [#/cm3]\n", r.NumberPM25)
	fmt.Printf("Number Concentration PM4.0: %3.2f [#/cm3]\n", r.NumberPM4)
	fmt.Printf("Number Concentration PM10 : %3.2f [#/cm3]\n", r.NumberPM10)
	fmt.Printf("Typical Particle Size     : %3.2f [μm]\n", r.TypicalParticleSize)
	fmt.Println("---------------------------------------")
}
