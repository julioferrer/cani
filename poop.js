const SerialPort = require('serialport')

const ByteLength = require('@serialport/parser-byte-length')
const port = new SerialPort("/dev/ttyUSB0", { baudRate: 9600})
const parser = port.pipe(new ByteLength({length: 4}))
parser.on('data', (data) => console.log(data.toString('utf8')))
