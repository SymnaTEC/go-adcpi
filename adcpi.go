/*
 ================================================
 ABElectronics UK ADC Pi 8-Channel Analogue to Digital Converter
 Version 1.0 Created 21/06/2017
 ================================================

 Reads from the MCP3424 ADC on the ADC Pi and ADC Pi Plus.

 Required package:
    apt-get install libi2c-dev
 */

package adcpi

const I2CBUS = 1

type Interface struct {
    i2cBus *I2C
    i2cAddress byte
    config byte
    currentChannel byte
    bitRate byte
    conversionMode byte
    pga float64
    lsb float64
    signBit byte
}

func ADCPI(address byte, rate byte) Interface {
    adcpi := Interface {
        signBit:0,
        i2cAddress:address,
        config:0x9C,      // PGAx1, 18 bit, continuous conversion, channel 1
        currentChannel:1, // channel variable for adc 1
        bitRate:18,       // current bitrate
        conversionMode:1, // Conversion Mode
        pga:0.5,          // current pga setting
        lsb:0.000007812,  // default lsb value for 18 bit
    }

    return adcpi
}

/*
 Reads the raw value from the selected ADC channel
 @param channel - 1 to 8
 @returns - raw long value from ADC buffer
 */
func (adcpi Interface) ReadRaw(channel byte) int {

    // variables for storing the raw bytes from the ADC
    var h byte = 0
    var l byte = 0
    var m byte = 0
    var s byte = 0
    var address byte = 0
    var t int = 0
    adcpi.signBit = 0

    // get the config and i2c address for the selected channel
    address = adcpi.setChannel(channel)

    // open the i2c bus
    bus,err := NewI2C(address, I2CBUS)
    if err != nil {
        panic(err)
    }
    adcpi.i2cBus = bus

    // if the conversion mode is set to one-shot update the ready bit to 1
    if adcpi.conversionMode == 0 {
       adcpi.i2cBus.WriteByte(adcpi.updateByte(adcpi.config, 7, 1))
    }

    // keep reading the ADC data until the conversion result is ready
    timeout := 1000 // number of reads before a timeout occurs
    x := 0
    for true {
        if adcpi.bitRate == 18 {
            buffer := make([]byte, 4)
            adcpi.i2cBus.Read(buffer)
            h = buffer[0]
            m = buffer[1]
            l = buffer[2]
            s = buffer[3]
        } else {

            buffer := make([]byte, 3)
            adcpi.i2cBus.Read(buffer)
            h = buffer[0]
            m = buffer[1]
            s = buffer[2]
        }

        // check bit 7 of s to see if the conversion result is ready
        if s & (1 << 7) == 0 {
            break
        }
        if x > timeout {
            // timeout occurred
            return 0
        }
        x++
    }

    // close the i2c bus
    adcpi.i2cBus.Close()

    // extract the returned bytes and combine in the correct order
    switch adcpi.bitRate {
    case 18:
        t = int(((h & 3) << 16) | (m << 8) | l)
        if (t >> 17) & 1 {
            adcpi.signBit = 1
            t &= ^(1 << 17)
        }
        break;
    case 16:
        t = int((h << 8) | m)
        if (t >> 15) & 1 {
            adcpi.signBit = 1
            t &= ^(1 << 15)
        }
        break;
    case 14:
        t = int(((h & 63) << 8) | m)
        if (t >> 13) & 1 {
            adcpi.signBit = 1
            t &= ^(1 << 13)
        }
        break;
    case 12:
        t = int(((h & 15) << 8) | m)
        if (t >> 11) & 1 {
            adcpi.signBit = 1
            t &= ^(1 << 11)
        }
        break;
    default:
        panic("ReadRaw() bitrate out of range")
    }
    return t
}

/*
 Returns the voltage from the selected ADC channel
 @param channel - 1 to 8
 @returns - double voltage value from ADC
 */
func (adcpi Interface) ReadVoltage(channel byte) float64 {
    raw := adcpi.ReadRaw(channel)
    if adcpi.signBit == 1 { // if the signbit is 1 the value is negative and most likely noise so it can be ignored.
        return 0
    } else {
        return float64(raw) * (adcpi.lsb / adcpi.pga) * 2.471 // calculate the voltage and return it
    }
}

/*
 Set the sample resolution
 @param rate - 12 = 12 bit(240SPS max), 14 = 14 bit(60SPS max), 16 = 16 bit(15SPS max), 18 = 18 bit(3.75SPS max)
 */
func (adcpi Interface) SetPGA(gain byte) {
    switch gain {
    case 1:
        adcpi.config = adcpi.updateByte(adcpi.config, 0, 0)
        adcpi.config = adcpi.updateByte(adcpi.config, 1, 0)
        adcpi.pga = 0.5
        break
    case 2:
        adcpi.config = adcpi.updateByte(adcpi.config, 0, 1)
        adcpi.config = adcpi.updateByte(adcpi.config, 1, 0)
        adcpi.pga = 1
        break
    case 4:
        adcpi.config = adcpi.updateByte(adcpi.config, 0, 0)
        adcpi.config = adcpi.updateByte(adcpi.config, 1, 1)
        adcpi.pga = 2
        break
    case 8:
        adcpi.config = adcpi.updateByte(adcpi.config, 0, 1)
        adcpi.config = adcpi.updateByte(adcpi.config, 1, 1)
        adcpi.pga = 4
        break
    default:
        panic("SetPga() gain out of range: 1, 2, 4, 8")
    }

    // write the changes
    adcpi.writeByte(adcpi.i2cAddress, adcpi.config)
    adcpi.writeByte(adcpi.i2cAddress + 1, adcpi.config)
}

/*
 Set the sample resolution
 @param rate - 12 = 12 bit(240SPS max), 14 = 14 bit(60SPS max), 16 = 16 bit(15SPS max), 18 = 18 bit(3.75SPS max)
 */
func (adcpi Interface) SetBitRate(rate byte) {
    switch rate {
    case 12:
        adcpi.config = adcpi.updateByte(adcpi.config, 2, 0)
        adcpi.config = adcpi.updateByte(adcpi.config, 3, 0)
        adcpi.bitRate = 12
        adcpi.lsb = 0.0005
        break
    case 14:
        adcpi.config = adcpi.updateByte(adcpi.config, 2, 1)
        adcpi.config = adcpi.updateByte(adcpi.config, 3, 0)
        adcpi.bitRate = 14
        adcpi.lsb = 0.000125
        break
    case 16:
        adcpi.config = adcpi.updateByte(adcpi.config, 2, 0)
        adcpi.config = adcpi.updateByte(adcpi.config, 3, 1)
        adcpi.bitRate = 16
        adcpi.lsb = 0.00003125
        break
    case 18:
        adcpi.config = adcpi.updateByte(adcpi.config, 2, 1)
        adcpi.config = adcpi.updateByte(adcpi.config, 3, 1)
        adcpi.bitRate = 18
        adcpi.lsb = 0.0000078125
        break
    default:
        panic("SetBitRate() rate out of range: 12, 14, 16, 18")
    }

    // write the changes
    adcpi.writeByte(adcpi.i2cAddress, adcpi.config)
    adcpi.writeByte(adcpi.i2cAddress + 1, adcpi.config)
}

/*
 Set the conversion mode for ADC
 @param mode - 0 = One shot conversion mode, 1 = Continuous conversion mode
*/
func (adcpi Interface) SetConversationMode(mode byte) {
    if mode == 0 {
        adcpi.config = adcpi.updateByte(adcpi.config, 4, 0)
        adcpi.conversionMode = 0
    } else if mode == 1 {
        adcpi.config = adcpi.updateByte(adcpi.config, 4, 1)
        adcpi.conversionMode = 1
    } else {
        panic("SetConversationMode() mode out of range: 0 or 1")
    }
}

/*
 private method for writing a byte to the I2C port
 */
func (adcpi Interface) writeByte(address byte, value byte) {
    bus,err := NewI2C(address, I2CBUS)
    if err != nil {
        panic(err)
    }
    bus.WriteByte(value)
    bus.Close()
}

/*
 private method for setting the value of a single bit within a byte
*/
func (adcpi Interface) updateByte(_byte byte, bit byte, value byte) byte {
    if value == 0 {
        return _byte & ^(1 << bit)
    } else if value == 1 {
        return _byte | 1 << bit
    } else {
        panic("updateByte() value out of range: 0 or 1")
    }
}

/*
 private method for setting the channel
 */
func (adcpi Interface) setChannel(channel byte) byte {
    if channel == 1 || channel == 5 {
        adcpi.config = adcpi.updateByte(adcpi.config, 5, 0)
        adcpi.config = adcpi.updateByte(adcpi.config, 6, 0)
    } else if channel == 2 || channel == 6 {
        adcpi.config = adcpi.updateByte(adcpi.config, 5, 1)
        adcpi.config = adcpi.updateByte(adcpi.config, 6, 0)
    } else if channel == 3 || channel == 7 {
        adcpi.config = adcpi.updateByte(adcpi.config, 5, 0)
        adcpi.config = adcpi.updateByte(adcpi.config, 6, 1)
    } else if channel == 4 || channel == 8 {
        adcpi.config = adcpi.updateByte(adcpi.config, 5, 1)
        adcpi.config = adcpi.updateByte(adcpi.config, 6, 1)
    } else {
        panic("setChannel() value out of range: 1 to 8")
    }
    adcpi.currentChannel = channel
    if channel > 4 {
        return adcpi.i2cAddress + 1
    } else {
        return adcpi.i2cAddress
    }
}