//go:build stm32f3

package machine

import (
	"device/stm32"
	"unsafe"
)

const (
	AF4_I2C1_2_3 = 4
)

// CPUFrequency returns the configured core clock for STM32F303 targets.
func CPUFrequency() uint32 {
	return 72000000
}

var deviceIDAddr = []uintptr{0x1FFFF7AC, 0x1FFFF7B0, 0x1FFFF7B4}

const (
	PA0  = portA + 0
	PA1  = portA + 1
	PA2  = portA + 2
	PA3  = portA + 3
	PA4  = portA + 4
	PA5  = portA + 5
	PA6  = portA + 6
	PA7  = portA + 7
	PA8  = portA + 8
	PA9  = portA + 9
	PA10 = portA + 10
	PA11 = portA + 11
	PA12 = portA + 12
	PA13 = portA + 13
	PA14 = portA + 14
	PA15 = portA + 15

	PB0  = portB + 0
	PB1  = portB + 1
	PB2  = portB + 2
	PB3  = portB + 3
	PB4  = portB + 4
	PB5  = portB + 5
	PB6  = portB + 6
	PB7  = portB + 7
	PB8  = portB + 8
	PB9  = portB + 9
	PB10 = portB + 10
	PB11 = portB + 11
	PB12 = portB + 12
	PB13 = portB + 13
	PB14 = portB + 14
	PB15 = portB + 15

	PC0  = portC + 0
	PC1  = portC + 1
	PC2  = portC + 2
	PC3  = portC + 3
	PC4  = portC + 4
	PC5  = portC + 5
	PC6  = portC + 6
	PC7  = portC + 7
	PC8  = portC + 8
	PC9  = portC + 9
	PC10 = portC + 10
	PC11 = portC + 11
	PC12 = portC + 12
	PC13 = portC + 13
	PC14 = portC + 14
	PC15 = portC + 15

	PD0  = portD + 0
	PD1  = portD + 1
	PD2  = portD + 2
	PD3  = portD + 3
	PD4  = portD + 4
	PD5  = portD + 5
	PD6  = portD + 6
	PD7  = portD + 7
	PD8  = portD + 8
	PD9  = portD + 9
	PD10 = portD + 10
	PD11 = portD + 11
	PD12 = portD + 12
	PD13 = portD + 13
	PD14 = portD + 14
	PD15 = portD + 15

	PE0  = portE + 0
	PE1  = portE + 1
	PE2  = portE + 2
	PE3  = portE + 3
	PE4  = portE + 4
	PE5  = portE + 5
	PE6  = portE + 6
	PE7  = portE + 7
	PE8  = portE + 8
	PE9  = portE + 9
	PE10 = portE + 10
	PE11 = portE + 11
	PE12 = portE + 12
	PE13 = portE + 13
	PE14 = portE + 14
	PE15 = portE + 15

	PF0  = portF + 0
	PF1  = portF + 1
	PF2  = portF + 2
	PF3  = portF + 3
	PF4  = portF + 4
	PF5  = portF + 5
	PF6  = portF + 6
	PF7  = portF + 7
	PF8  = portF + 8
	PF9  = portF + 9
	PF10 = portF + 10
	PF11 = portF + 11
	PF12 = portF + 12
	PF13 = portF + 13
	PF14 = portF + 14
	PF15 = portF + 15
)

func (p Pin) getPort() *stm32.GPIO_Type {
	switch p / 16 {
	case 0:
		return stm32.GPIOA
	case 1:
		return stm32.GPIOB
	case 2:
		return stm32.GPIOC
	case 3:
		return stm32.GPIOD
	case 4:
		return stm32.GPIOE
	case 5:
		return stm32.GPIOF
	default:
		panic("machine: unknown port")
	}
}

func (p Pin) enableClock() {
	switch p / 16 {
	case 0:
		stm32.RCC.AHBENR.SetBits(stm32.RCC_AHBENR_GPIOAEN)
	case 1:
		stm32.RCC.AHBENR.SetBits(stm32.RCC_AHBENR_GPIOBEN)
	case 2:
		stm32.RCC.AHBENR.SetBits(stm32.RCC_AHBENR_GPIOCEN)
	case 3:
		stm32.RCC.AHBENR.SetBits(stm32.RCC_AHBENR_GPIODEN)
	case 4:
		stm32.RCC.AHBENR.SetBits(stm32.RCC_AHBENR_GPIOEEN)
	case 5:
		stm32.RCC.AHBENR.SetBits(stm32.RCC_AHBENR_GPIOFEN)
	default:
		panic("machine: unknown port")
	}
}

func enableAltFuncClock(bus unsafe.Pointer) {
	switch bus {
	case unsafe.Pointer(stm32.I2C1):
		stm32.RCC.APB1ENR.SetBits(stm32.RCC_APB1ENR_I2C1EN)
	case unsafe.Pointer(stm32.I2C2):
		stm32.RCC.APB1ENR.SetBits(stm32.RCC_APB1ENR_I2C2EN)
	default:
		// not supported yet
	}
}

// I2C defaults for STM32F3Discovery.
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

func (i2c *I2C) getFreqRange(br uint32) uint32 {
	// Timing values for 36MHz PCLK1 (typical 72MHz SYSCLK / 2 APB1)
	switch br {
	case 10 * KHz:
		return 0xF000F3FE
	case 100 * KHz:
		return 0x10805E89
	case 400 * KHz:
		return 0x00901D5B
	case 500 * KHz:
		return 0x00801949
	default:
		return 0
	}
}
