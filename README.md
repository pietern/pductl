# pductl

Toggle PDU outlets with UDP packets.

I'm using this to toggle the power state of audio amplifiers when
their streamers are playing audio. When a streamer stops playing, and
an inactivity delay has passed, the amplifier is powered off.

I'm using this with an APC AP7920 Switched Rack PDU and connect to it
over a telnet connection. The scripting mode of this interface is
rather rudimentary but it works fine with an _expect_-like wrapper
(see `./pdu`).

## License

MIT.
