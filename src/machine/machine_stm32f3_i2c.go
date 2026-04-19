//go:build stm32f3

package machine

import "device/stm32"

const (
	AF4_I2C1_2_3 = 4
)

// I2C pins for STM32F3 defaults to I2C1 on PB6/PB7.
const (
	I2C0_SCL_PIN = PB6
	I2C0_SDA_PIN = PB7
)

var (
	I2C0 = &I2C{
		Bus:             stm32.I2C1,
		AltFuncSelector: AF4_I2C1_2_3,
	}
)

type I2C struct {
	Bus             *stm32.I2C_Type
	AltFuncSelector uint8
}

func (i2c *I2C) configurePins(config I2CConfig) {
	config.SCL.ConfigureAltFunc(PinConfig{Mode: PinModeI2CSCL}, i2c.AltFuncSelector)
	config.SDA.ConfigureAltFunc(PinConfig{Mode: PinModeI2CSDA}, i2c.AltFuncSelector)
}

func (i2c *I2C) getFreqRange(config I2CConfig) uint32 {
	clock := CPUFrequency() / 2 // APB1 on STM32F303 default clock tree
	clock /= 1000000
	if clock < 2 {
		clock = 2
	} else if clock > 50 {
		clock = 50
	}
	if config.Frequency > 100000 && clock < 4 {
		clock = 4
	}
	return clock << stm32.I2C_CR2_FREQ_Pos
}

func (i2c *I2C) getRiseTime(config I2CConfig) uint32 {
	freqRange := i2c.getFreqRange(config)
	if config.Frequency > 100000 {
		freqRange *= 300
		freqRange /= 1000
	}
	return (freqRange + 1) << stm32.I2C_TRISE_TRISE_Pos
}

func (i2c *I2C) getSpeed(config I2CConfig) uint32 {
	ccr := func(pclk uint32, freq uint32, coeff uint32) uint32 {
		return (((pclk - 1) / (freq * coeff)) + 1) & stm32.I2C_CCR_CCR_Msk
	}
	sm := func(pclk uint32, freq uint32) uint32 {
		if s := ccr(pclk, freq, 2); s < 4 {
			return 4
		} else {
			return s
		}
	}
	fm := func(pclk uint32, freq uint32, duty uint8) uint32 {
		if duty == DutyCycle2 {
			return ccr(pclk, freq, 3)
		}
		return ccr(pclk, freq, 25) | stm32.I2C_CCR_DUTY
	}

	clock := CPUFrequency() / 2 // APB1
	if config.Frequency <= 100000 {
		return sm(clock, config.Frequency)
	}
	s := fm(clock, config.Frequency, config.DutyCycle)
	if (s & stm32.I2C_CCR_CCR_Msk) == 0 {
		return 1
	}
	return s | stm32.I2C_CCR_F_S
}
