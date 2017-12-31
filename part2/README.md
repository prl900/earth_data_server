# earth_data_server (part 1)

This code is referenced on the Medium blog post:

To run this code you'll need a Go 1.0+ environment set in your system.
In Linux/Mac OSX do:

1.- Download the Blue Marble Next Generation file on this folder:
`$ curl -O "https://eoimages.gsfc.nasa.gov/images/imagerecords/73000/73909/world.topo.bathy.200412.3x21600x10800.png"`

2.- Test all different compression methods on the Blue Marble image:
`$ go run compare_compressors.go`

3.- Generate the PNG, Raw and Snappy tiles for this file:
`$ go run generate_tiles.go`

4.- Request a region providing the coordinates of any place in the world and the RGB channel. The result is computed three times by each method recording the time taken to generate the region:
`$ time go run get_region.go -lat 42 -lon -1 -chan 0`
